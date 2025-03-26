package aichat

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/tx"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/validator"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/gpt"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
)

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

	jsonPaymentMethods, jsonTransactionCategories, jsonInstitutions, err := uc.jsonMarshalEntities(
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
			jsonPaymentMethods,
			jsonTransactionCategories,
			jsonInstitutions,
		)
		return err
	})

	if err := g.Wait(); err != nil {
		return nil, errs.New(err)
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

	completion, err := uc.gp.Completion(ctx, messages)
	if err != nil {
		return "", errs.New(err)
	}

	return completion.Content, nil
}

func (uc *GenerateAIChatMessageUseCase) jsonMarshalEntities(
	paymentMethods []entity.PaymentMethod,
	transactionCategories []entity.TransactionCategory,
	institutions []entity.Institution,
) (string, string, string, error) {
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

	paymentMethodsJSON, err := json.Marshal(paymentMethodsMap)
	if err != nil {
		return "", "", "", errs.New(err)
	}

	transactionCategoriesJSON, err := json.Marshal(transactionCategoriesMap)
	if err != nil {
		return "", "", "", errs.New(err)
	}

	institutionsJSON, err := json.Marshal(institutionsMap)
	if err != nil {
		return "", "", "", errs.New(err)
	}

	return string(paymentMethodsJSON),
		string(transactionCategoriesJSON),
		string(institutionsJSON),
		nil
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
	jsonPaymentMethods string,
	jsonTransactionCategories string,
	jsonInstitutions string,
) (string, error) {
	systemMessage := fmt.Sprintf(
		`You are a financial planning specialist with access to a system containing your clients' transactions. The system also includes the following information:

- Current Date and Time: %s
- Payment Methods: %s
- Transaction Categories: %s
- Institutions: %s

Using your expertise in financial planning, please analyze the information provided and offer data-driven insights along with actionable advice.`,
		time.Now().Format(time.RFC3339),
		jsonPaymentMethods,
		jsonTransactionCategories,
		jsonInstitutions,
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

	tools := []gpt.Tool{
		{
			Name:        "list_user_transactions",
			Description: "List user transactions and sum them up. The transactions are filtered by the given parameters.",
			Func:        uc.getUserTransactions(in),
			Args: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"start_date": map[string]any{
						"type":        "string",
						"description": "The start date for filtering transactions (RFC3339 format).",
					},
					"end_date": map[string]any{
						"type":        "string",
						"description": "The end date for filtering transactions (RFC3339 format).",
					},
					"category_ids": map[string]any{
						"type": "array",
						"items": map[string]any{
							"type": "string",
						},
						"description": "The category IDs to filter transactions by",
					},
					"institution_ids": map[string]any{
						"type": "array",
						"items": map[string]any{
							"type": "string",
						},
						"description": "The institution IDs to filter transactions by",
					},
					"payment_method_ids": map[string]any{
						"type": "array",
						"items": map[string]any{
							"type": "string",
						},
						"description": "The payment method IDs to filter transactions by",
					},
					"is_expense": map[string]any{
						"type":        "boolean",
						"description": "Return only expenses if true.",
					},
					"is_income": map[string]any{
						"type":        "boolean",
						"description": "Return only income if true.",
					},
				},
				"required": []string{"start_date", "end_date"},
			},
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

func (uc *GenerateAIChatMessageUseCase) getUserTransactions(
	in GenerateAIChatMessageUseCaseInput,
) gpt.ToolFunc {
	return func(ctx context.Context, args map[string]any) (string, error) {
		g, subCtx := errgroup.WithContext(ctx)

		startDate, err := uc.parseRequiredDateArg(args, "start_date")
		if err != nil {
			return "", errs.New(err)
		}

		endDate, err := uc.parseRequiredDateArg(args, "end_date")
		if err != nil {
			return "", errs.New(err)
		}

		opts := []repo.TransactionOption{
			repo.WithTransactionDateAfter(*startDate),
			repo.WithTransactionDateBefore(*endDate),
		}

		categoryIDs := uc.parseUUIDsArg(args, "category_ids")
		if len(categoryIDs) > 0 {
			opts = append(
				opts,
				repo.WithTransactionCategories(categoryIDs...),
			)
		}

		institutionIDs := uc.parseUUIDsArg(args, "institution_ids")
		if len(institutionIDs) > 0 {
			opts = append(
				opts,
				repo.WithTransactionInstitutions(institutionIDs...),
			)
		}

		paymentMethodIDs := uc.parseUUIDsArg(args, "payment_method_ids")
		if len(paymentMethodIDs) > 0 {
			opts = append(
				opts,
				repo.WithTransactionInstitutions(paymentMethodIDs...),
			)
		}

		isExpense := uc.parseBoolArg(args, "is_expense")
		if isExpense {
			opts = append(opts, repo.WithTransactionIsExpense(isExpense))
		}

		isIncome := uc.parseBoolArg(args, "is_income")
		if isIncome {
			opts = append(opts, repo.WithTransactionIsIncome(isIncome))
		}

		var (
			transactions []entity.Transaction
			sum          int64
		)
		g.Go(func() (err error) {
			transactions, err = uc.tr.ListTransactions(
				subCtx,
				in.UserID,
				opts...)
			return err
		})

		g.Go(func() (err error) {
			sum, err = uc.tr.SumTransactions(subCtx, in.UserID, opts...)
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
	strs, ok := args[key].([]string)
	if !ok {
		return nil
	}

	uuids := make([]uuid.UUID, len(strs))
	for i, str := range strs {
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
