package aichat

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/usecase/account"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/usecase/budget"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/ptr"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/tx"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/validator"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/gpt"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
)

var errorMessages = []string{
	"Desculpe, mas não entendemos sua pergunta. Poderia reformulá-la?",
	"Lamentamos, mas sua pergunta não ficou clara. Poderia reformular?",
	"Desculpe, não foi possível interpretar sua dúvida. Poderia formular de outra forma?",
	"Desculpe, não compreendemos sua questão. Poderia reescrevê-la?",
	"Sentimos muito, mas não conseguimos entender o que você perguntou. Poderia reformular a questão?",
}

type GenerateAIChatMessageUseCase struct {
	v     *validator.Validator
	tx    tx.TX
	gp    gpt.GPT
	acr   repo.AIChatRepo
	acmr  repo.AIChatMessageRepo
	acmrr repo.AIChatAnswerRepo
	pmr   repo.PaymentMethodRepo
	tcr   repo.TransactionCategoryRepo
	ir    repo.InstitutionRepo
	tr    repo.TransactionRepo
	gbuc  *budget.GetBudgetUseCase
	gabuc *account.GetAccountsBalanceUseCase
}

func NewGenerateAIChatMessageUseCase(
	v *validator.Validator,
	tx tx.TX,
	gp gpt.GPT,
	acr repo.AIChatRepo,
	acmr repo.AIChatMessageRepo,
	acmrr repo.AIChatAnswerRepo,
	pmr repo.PaymentMethodRepo,
	tcr repo.TransactionCategoryRepo,
	ir repo.InstitutionRepo,
	tr repo.TransactionRepo,
	gbuc *budget.GetBudgetUseCase,
	gabuc *account.GetAccountsBalanceUseCase,
) *GenerateAIChatMessageUseCase {
	return &GenerateAIChatMessageUseCase{
		v:     v,
		tx:    tx,
		gp:    gp,
		acr:   acr,
		acmr:  acmr,
		acmrr: acmrr,
		pmr:   pmr,
		tcr:   tcr,
		ir:    ir,
		tr:    tr,
		gbuc:  gbuc,
		gabuc: gabuc,
	}
}

type GenerateAIChatMessageUseCaseInput struct {
	UserID   uuid.UUID   `json:"-"       validate:"required"`
	AIChatID uuid.UUID   `json:"-"       validate:"required"`
	Message  string      `json:"message" validate:"required,max=512"`
	Tier     entity.Tier `json:"-"       validate:"required,oneof=TRIAL PREMIUM"`
}

type GenerateAIChatMessageUseCaseOutput struct {
	*entity.AIChat
	AIChatAnswer *entity.AIChatAnswer `json:"ai_chat_answer"`
}

func (uc *GenerateAIChatMessageUseCase) Execute(
	ctx context.Context,
	in GenerateAIChatMessageUseCaseInput,
) (*GenerateAIChatMessageUseCaseOutput, error) {
	err := uc.v.Validate(in)
	if err != nil {
		return nil, errs.New(err)
	}

	g, subCtx := errgroup.WithContext(ctx)

	var (
		aiChat                *entity.AIChat
		chatRecentHistory     []entity.AIChatMessageAndAnswer
		paymentMethods        []entity.PaymentMethod
		transactionCategories []entity.TransactionCategory
		institutions          []entity.Institution
	)

	g.Go(func() (err error) {
		aiChat, err = uc.acr.GetAIChatByID(subCtx, in.AIChatID)
		return err
	})

	g.Go(func() (err error) {
		chatRecentHistory, err = uc.acr.ListAIChatMessagesAndAnswers(
			ctx,
			repo.ListAIChatMessagesAndAnswersParams{
				AiChatID: in.AIChatID,
				Limit:    6,
			},
		)
		return err
	})

	g.Go(func() (err error) {
		paymentMethods, err = uc.pmr.ListPaymentMethods(subCtx)
		return err
	})

	g.Go(func() (err error) {
		transactionCategories, err = uc.tcr.ListTransactionCategories(subCtx)
		return err
	})

	g.Go(func() (err error) {
		institutions, err = uc.ir.ListInstitutions(subCtx)
		return err
	})

	if err := g.Wait(); err != nil {
		return nil, errs.New(err)
	}

	if aiChat == nil || aiChat.UserID != in.UserID {
		return nil, errs.ErrAIChatNotFound
	}

	entitiesMap, err := uc.simplifyEntities(
		paymentMethods,
		transactionCategories,
		institutions,
	)
	if err != nil {
		return nil, errs.New(err)
	}

	chatHistory := uc.parseChatHistory(chatRecentHistory)

	g, subCtx = errgroup.WithContext(ctx)
	var title, answer string

	if aiChat.Title == nil {
		g.Go(func() (err error) {
			title, err = uc.generateAIChatTitle(subCtx, in.Message)
			return err
		})
	}

	g.Go(func() (err error) {
		answer, err = uc.generateAIChatAnswer(
			ctx,
			in,
			chatHistory,
			entitiesMap...,
		)
		return err
	})

	if err := g.Wait(); err != nil {
		slog.Error(
			"failed to generate AI chat answer",
			"user_id",
			in.UserID,
			"message",
			in.Message,
			"error",
			err,
		)
		title = ""
		answer = errorMessages[rand.IntN(len(errorMessages))]
	}

	var aiChatAnswer *entity.AIChatAnswer
	err = uc.tx.Do(ctx, func(ctx context.Context) error {
		if title != "" {
			err = uc.acr.UpdateAIChat(ctx, repo.UpdateAIChatParams{
				ID:    in.AIChatID,
				Title: &title,
			})

			aiChat.Title = &title
			aiChat.UpdatedAt = time.Now()
		}

		aiChatMessage, err := uc.acmr.CreateAIChatMessage(
			ctx,
			repo.CreateAIChatMessageParams{
				AiChatID: in.AIChatID,
				Message:  in.Message,
			},
		)
		if err != nil {
			return errs.New(err)
		}

		aiChatAnswer, err = uc.acmrr.CreateAIChatAnswer(
			ctx,
			repo.CreateAIChatAnswerParams{
				AiChatMessageID: aiChatMessage.ID,
				Message:         answer,
			},
		)
		if err != nil {
			return errs.New(err)
		}

		return nil
	})
	if err != nil {
		return nil, errs.New(err)
	}

	return &GenerateAIChatMessageUseCaseOutput{
		AIChat:       aiChat,
		AIChatAnswer: aiChatAnswer,
	}, nil
}

func (uc *GenerateAIChatMessageUseCase) generateAIChatTitle(
	ctx context.Context,
	message string,
) (string, error) {
	systemMessage := "Generate a short title the user message"

	messages := []gpt.Message{
		{
			Role:    gpt.RoleSystem,
			Content: systemMessage,
		},
		{
			Role:    gpt.RoleUser,
			Content: message,
		},
	}

	completion, err := uc.gp.Completion(
		ctx,
		messages,
		gpt.WithModel(gpt.Model4oMini),
	)
	if err != nil {
		return "", errs.New(err)
	}

	return completion.Content, nil
}

func (uc *GenerateAIChatMessageUseCase) simplifyEntities(
	paymentMethods []entity.PaymentMethod,
	transactionCategories []entity.TransactionCategory,
	institutions []entity.Institution,
) ([][]map[string]string, error) {
	paymentMethodsMap := []map[string]string{}
	for _, pm := range paymentMethods {
		paymentMethodsMap = append(paymentMethodsMap, map[string]string{
			"id":   pm.ID.String(),
			"name": pm.Name,
		})
	}

	transactionCategoriesMap := []map[string]string{}
	for _, tc := range transactionCategories {
		transactionCategoriesMap = append(
			transactionCategoriesMap,
			map[string]string{
				"id":   tc.ID.String(),
				"name": tc.Name,
			},
		)
	}

	institutionsMap := []map[string]string{}
	for _, i := range institutions {
		institutionsMap = append(institutionsMap, map[string]string{
			"id":   i.ID.String(),
			"name": i.Name,
		})
	}

	entities := [][]map[string]string{
		paymentMethodsMap,
		transactionCategoriesMap,
		institutionsMap,
	}

	return entities, nil
}

func (uc *GenerateAIChatMessageUseCase) jsonMarshalEntities(
	entities ...[]map[string]string,
) ([]string, error) {
	var jsonStrings []string
	for _, entity := range entities {
		entityJSON, err := json.Marshal(entity)
		if err != nil {
			return nil, errs.New(err)
		}
		jsonStrings = append(jsonStrings, string(entityJSON))
	}

	return jsonStrings, nil
}

func (uc *GenerateAIChatMessageUseCase) parseEntityNames(
	entities ...[]map[string]string,
) ([][]string, error) {
	names := make([][]string, len(entities))
	for i, entity := range entities {
		names[i] = make([]string, len(entity))
		for j, e := range entity {
			name, ok := e["name"]
			if !ok {
				return nil, fmt.Errorf("name not found in entity %d", i)
			}
			names[i][j] = name
		}
	}

	return names, nil
}

func (uc *GenerateAIChatMessageUseCase) parseChatHistory(
	chatHistory []entity.AIChatMessageAndAnswer,
) []gpt.Message {
	parsed := []gpt.Message{}
	mapAuthorToRole := map[entity.AIChatMessageAuthor]gpt.Role{
		entity.AIChatMessageAuthorUser: gpt.RoleUser,
		entity.AIChatMessageAuthorAI:   gpt.RoleAssistant,
	}
	for _, msg := range chatHistory {
		parsed = append(parsed, gpt.Message{
			Role:    mapAuthorToRole[msg.Author],
			Content: msg.Message,
		})
	}
	return parsed
}

func (uc *GenerateAIChatMessageUseCase) generateAIChatAnswer(
	ctx context.Context,
	in GenerateAIChatMessageUseCaseInput,
	chatHistory []gpt.Message,
	entities ...[]map[string]string,
) (string, error) {
	systemMessage := fmt.Sprintf(
		`Today is %s and you are a financial planning specialist with access to a system containing your clients' transactions.
Using your expertise in financial planning, please analyze the information provided and offer data-driven insights along with actionable advice.`,
		time.Now().Format(time.RFC3339),
	)

	messages := append([]gpt.Message{
		{
			Role:    gpt.RoleSystem,
			Content: systemMessage,
		},
	}, chatHistory...)

	messages = append(messages, gpt.Message{
		Role:    gpt.RoleUser,
		Content: in.Message,
	})

	entityNames, err := uc.parseEntityNames(entities...)
	if err != nil {
		return "", errs.New(err)
	}

	paymentMethodNames := strings.Join(entityNames[0], ", ")
	transactionCategoryNames := strings.Join(entityNames[1], ", ")
	institutionNames := strings.Join(entityNames[2], ", ")

	jsonEntities, err := uc.jsonMarshalEntities(entities...)
	if err != nil {
		return "", errs.New(err)
	}
	jsonPaymentMethods := jsonEntities[0]
	jsonTransactionCategories := jsonEntities[1]
	jsonInstitutions := jsonEntities[2]

	tools := []gpt.Tool{
		{
			Name: "list_user_transactions",
			Description: fmt.Sprintf(
				"List user transactions and sum them up. The transactions are about all user financial data movement and it can be filtered by start and end date, categories (%s), institutions (%s), payment methods (%s) and if it is expense or income.",
				transactionCategoryNames,
				institutionNames,
				paymentMethodNames,
			),
			Func: uc.listUserTransactions(in.UserID),
			Args: buildListTransactionArgs(
				jsonPaymentMethods,
				jsonTransactionCategories,
				jsonInstitutions,
			),
		},
		{
			Name:        "get_user_accounts_balance",
			Description: "Get user's accounts balance filtered by institution, start and end date. If no institution is provided all accounts are considered.",
			Func:        uc.getUserAccountsBalance(in.UserID),
			Args:        buildGetUserAccountsBalanceArgs(jsonInstitutions),
		},
		{
			Name:        "get_user_budget",
			Description: "Get user's budget definitions and the amount they spent in a specific month.",
			Func:        uc.getUserBudget(in.UserID),
			Args:        buildGetBudgetArgs(),
		},
	}

	message, err := uc.gp.Completion(
		ctx,
		messages,
		gpt.WithTools(tools),
	)
	if err != nil {
		return "", errs.New(err)
	}

	return message.Content, nil
}

func (uc *GenerateAIChatMessageUseCase) listUserTransactions(
	userID uuid.UUID,
) gpt.ToolFunc {
	return func(ctx context.Context, args map[string]any) (string, error) {
		g, subCtx := errgroup.WithContext(ctx)

		opts, err := uc.buildTransactionOptionsFromArgs(args)
		if err != nil {
			return "", errs.New(err)
		}

		var (
			transactions []entity.Transaction
			sum          int64
		)
		g.Go(func() (err error) {
			transactions, err = uc.tr.ListTransactions(
				subCtx,
				userID,
				*opts)
			return err
		})

		g.Go(func() (err error) {
			sum, err = uc.tr.SumTransactions(subCtx, userID, *opts)
			return err
		})

		if err := g.Wait(); err != nil {
			return "", errs.New(err)
		}

		response := map[string]any{
			"description":  "List of user transactions, where amounts and sums are given in cents. Negative values represent expenses; positive values represent income.",
			"transactions": transactions,
			"sum":          sum,
		}
		responseJSON, err := json.Marshal(response)
		if err != nil {
			return "", errs.New(err)
		}

		return string(responseJSON), nil
	}
}

func (uc *GenerateAIChatMessageUseCase) buildTransactionOptionsFromArgs(
	args map[string]any,
) (*repo.TransactionOptions, error) {
	startDate, err := uc.parseRequiredDateArg(args, ArgKeyStartDate)
	if err != nil {
		return nil, errs.New(err)
	}

	endDate, err := uc.parseRequiredDateArg(args, ArgKeyEndDate)
	if err != nil {
		return nil, errs.New(err)
	}

	opts := repo.TransactionOptions{
		StartDate: *startDate,
		EndDate:   *endDate,
		IsIgnored: ptr.New(false),
	}

	categoryIDs := uc.parseUUIDsArg(args, ArgKeyCategoryIDs)
	if len(categoryIDs) > 0 {
		opts.CategoryIDs = categoryIDs
	}

	institutionIDs := uc.parseUUIDsArg(args, ArgKeyInstitutionIDs)
	if len(institutionIDs) > 0 {
		opts.InstitutionIDs = institutionIDs
	}

	paymentMethodIDs := uc.parseUUIDsArg(args, ArgKeyPaymentMethodIDs)
	if len(paymentMethodIDs) > 0 {
		opts.PaymentMethodIDs = paymentMethodIDs
	}

	isExpense := uc.parseBoolArg(args, ArgKeyIsExpense)
	if isExpense {
		opts.IsExpense = isExpense
	}

	isIncome := uc.parseBoolArg(args, ArgKeyIsIncome)
	if isIncome {
		opts.IsIncome = isIncome
	}

	search, _ := args[ArgKeySearch].(string)
	if search != "" {
		opts.Search = search
	}

	return &opts, nil
}

func (uc *GenerateAIChatMessageUseCase) getUserAccountsBalance(
	userID uuid.UUID,
) gpt.ToolFunc {
	return func(ctx context.Context, args map[string]any) (string, error) {
		opts, err := uc.buildAccountsBalanceOptionsFromArgs(args)
		if err != nil {
			return "", errs.New(err)
		}

		budgetOutput, err := uc.gabuc.Execute(
			ctx,
			account.GetAccountsBalanceUseCaseInput{
				UserID:             userID,
				TransactionOptions: *opts,
			},
		)
		if err != nil {
			return "", errs.New(err)
		}

		response := map[string]any{
			"description": "The user's budget for the given month, where amounts are given in cents and percentages are given as integers (example: 1055 represents 10.55%).",
			"data":        budgetOutput,
		}
		responseJSON, err := json.Marshal(response)
		if err != nil {
			return "", errs.New(err)
		}

		return string(responseJSON), nil
	}
}

func (uc *GenerateAIChatMessageUseCase) buildAccountsBalanceOptionsFromArgs(
	args map[string]any,
) (*repo.TransactionOptions, error) {
	startDate, err := uc.parseRequiredDateArg(args, ArgKeyStartDate)
	if err != nil {
		return nil, errs.New(err)
	}

	endDate, err := uc.parseRequiredDateArg(args, ArgKeyEndDate)
	if err != nil {
		return nil, errs.New(err)
	}

	opts := repo.TransactionOptions{
		StartDate: *startDate,
		EndDate:   *endDate,
		IsIgnored: ptr.New(false),
	}

	institutionIDs := uc.parseUUIDsArg(args, ArgKeyInstitutionIDs)
	if len(institutionIDs) > 0 {
		opts.InstitutionIDs = institutionIDs
	}

	return &opts, nil
}

func (uc *GenerateAIChatMessageUseCase) getUserBudget(
	userID uuid.UUID,
) gpt.ToolFunc {
	return func(ctx context.Context, args map[string]any) (string, error) {
		g, subCtx := errgroup.WithContext(ctx)

		date, err := uc.parseRequiredDateArg(args, ArgKeyDate)
		if err != nil {
			return "", errs.New(err)
		}

		budgetOutput, err := uc.gbuc.Execute(
			subCtx,
			budget.GetBudgetUseCaseInput{
				UserID: userID,
				Date:   *date,
			},
		)
		if err != nil {
			return "", errs.New(err)
		}

		if err := g.Wait(); err != nil {
			return "", errs.New(err)
		}

		response := map[string]any{
			"description": "The user's budget for the given month, where amounts are given in cents and percentages are given as integers (example: 1055 represents 10.55%).",
			"data":        budgetOutput,
		}
		responseJSON, err := json.Marshal(response)
		if err != nil {
			return "", errs.New(err)
		}

		return string(responseJSON), nil
	}
}

func (uc *GenerateAIChatMessageUseCase) parseRequiredDateArg(
	args map[string]any,
	key string,
) (*time.Time, error) {
	dateStr, ok := args[key].(string)
	if !ok {
		return nil, errs.New(fmt.Sprintf("%s is required", key))
	}

	date, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return nil, errs.New(err)
	}

	return &date, nil
}

func (uc *GenerateAIChatMessageUseCase) parseUUIDsArg(
	args map[string]any,
	key string,
) []uuid.UUID {
	values, ok := args[key].([]any)
	if !ok {
		return nil
	}

	uuids := make([]uuid.UUID, len(values))
	for i, v := range values {
		str, ok := v.(string)
		if !ok {
			return nil
		}
		id, err := uuid.Parse(str)
		if err != nil {
			return nil
		}
		uuids[i] = id
	}

	return uuids
}

func (uc *GenerateAIChatMessageUseCase) parseBoolArg(
	args map[string]any,
	key string,
) bool {
	parsed, ok := args[key].(bool)
	if !ok {
		return false
	}
	return parsed
}

type ArgKey = string

const (
	ArgKeyStartDate        ArgKey = "start_date"
	ArgKeyEndDate          ArgKey = "end_date"
	ArgKeyCategoryIDs      ArgKey = "category_ids"
	ArgKeyInstitutionIDs   ArgKey = "institution_ids"
	ArgKeyPaymentMethodIDs ArgKey = "payment_method_ids"
	ArgKeyIsExpense        ArgKey = "is_expense"
	ArgKeyIsIncome         ArgKey = "is_income"
	ArgKeySearch           ArgKey = "search"
	ArgKeyDate             ArgKey = "date"
)

func buildListTransactionArgs(
	jsonPaymentMethods,
	jsonTransactionCategories,
	jsonInstitutions string,
) map[string]any {
	var listTransactionArgs = map[string]any{
		"type": "object",
		"properties": map[string]any{
			ArgKeyStartDate: map[string]any{
				"type":        "string",
				"description": "The start date for filtering transactions (RFC3339 format in the GMT-3 timezone).",
			},
			ArgKeyEndDate: map[string]any{
				"type":        "string",
				"description": "The end date for filtering transactions (RFC3339 format in the GMT-3 timezone).",
			},
			ArgKeyCategoryIDs: map[string]any{
				"type": "array",
				"items": map[string]any{
					"type": "string",
				},
				"description": fmt.Sprintf(
					"The category IDs to filter transactions by (use empty array to list all). Here is all possible categories: %s",
					jsonTransactionCategories,
				),
			},
			ArgKeyInstitutionIDs: map[string]any{
				"type": "array",
				"items": map[string]any{
					"type": "string",
				},
				"description": fmt.Sprintf(
					"The institution IDs to filter transactions by (use empty array to list all). Here is all possible institutions: %s",
					jsonInstitutions,
				),
			},
			ArgKeyPaymentMethodIDs: map[string]any{
				"type": "array",
				"items": map[string]any{
					"type": "string",
				},
				"description": fmt.Sprintf(
					"The payment method IDs to filter transactions by (use empty array to list all). Here is all possible payment methods: %s",
					jsonPaymentMethods,
				),
			},
			ArgKeyIsExpense: map[string]any{
				"type":        "boolean",
				"description": "Return only expenses if true. (use false to return all)",
			},
			ArgKeyIsIncome: map[string]any{
				"type":        "boolean",
				"description": "Return only income if true. (use false to return all)",
			},
			ArgKeySearch: map[string]any{
				"type":        "string",
				"description": "Search for transactions by name (Prefer using category, institution and payment method IDs; use empty string to return all).",
			},
		},
		"required": []string{
			ArgKeyStartDate,
			ArgKeyEndDate,
			ArgKeyCategoryIDs,
			ArgKeyInstitutionIDs,
			ArgKeyPaymentMethodIDs,
			ArgKeyIsExpense,
			ArgKeyIsIncome,
			ArgKeySearch,
		},
		"additionalProperties": false,
	}

	return listTransactionArgs
}

func buildGetUserAccountsBalanceArgs(
	jsonInstitutions string,
) map[string]any {
	var getAccountsBalanceArgs = map[string]any{
		"type": "object",
		"properties": map[string]any{
			ArgKeyStartDate: map[string]any{
				"type":        "string",
				"description": "The start date for filtering balances (RFC3339 format in the GMT-3 timezone).",
			},
			ArgKeyEndDate: map[string]any{
				"type":        "string",
				"description": "The end date for filtering balances (RFC3339 format in the GMT-3 timezone).",
			},
			ArgKeyInstitutionIDs: map[string]any{
				"type": "array",
				"items": map[string]any{
					"type": "string",
				},
				"description": fmt.Sprintf(
					"The institution IDs to filter balances by (use empty array to list all). Here is all possible institutions: %s",
					jsonInstitutions,
				),
			},
		},
		"required": []string{
			ArgKeyStartDate,
			ArgKeyEndDate,
			ArgKeyInstitutionIDs,
		},
		"additionalProperties": false,
	}

	return getAccountsBalanceArgs
}

func buildGetBudgetArgs() map[string]any {
	var getBudgetArgs = map[string]any{
		"type": "object",
		"properties": map[string]any{
			ArgKeyDate: map[string]any{
				"type":        "string",
				"description": "The first day of the month to get the budget for (RFC3339 format in the GMT-3 timezone).",
			},
		},
		"required":             []string{ArgKeyDate},
		"additionalProperties": false,
	}
	return getBudgetArgs
}
