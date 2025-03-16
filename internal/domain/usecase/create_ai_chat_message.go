package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/openai/openai-go"
	"golang.org/x/sync/errgroup"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/tx"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/validator"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/db"
	gptOpenAI "github.com/danielmesquitta/api-finance-manager/internal/provider/gpt/openai"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
)

type CreateAIChatMessage struct {
	v     *validator.Validator
	tx    tx.TX
	db    *db.DB
	oa    *gptOpenAI.OpenAI
	acr   repo.AIChatRepo
	acmr  repo.AIChatMessageRepo
	acmrr repo.AIChatAnswerRepo
	pmr   repo.PaymentMethodRepo
	tcr   repo.TransactionCategoryRepo
	ir    repo.InstitutionRepo
}

func NewCreateAIChatMessage(
	v *validator.Validator,
	tx tx.TX,
	db *db.DB,
	oa *gptOpenAI.OpenAI,
	acr repo.AIChatRepo,
	acmr repo.AIChatMessageRepo,
	acmrr repo.AIChatAnswerRepo,
	pmr repo.PaymentMethodRepo,
	tcr repo.TransactionCategoryRepo,
	ir repo.InstitutionRepo,
) *CreateAIChatMessage {
	return &CreateAIChatMessage{
		v:     v,
		tx:    tx,
		db:    db,
		oa:    oa,
		acr:   acr,
		acmr:  acmr,
		acmrr: acmrr,
		pmr:   pmr,
		tcr:   tcr,
		ir:    ir,
	}
}

type CreateAIChatMessageInput struct {
	UserID   uuid.UUID `json:"-"       validate:"required"`
	AIChatID uuid.UUID `json:"-"       validate:"required"`
	Message  string    `json:"message" validate:"required,max=512"`
}

func (uc *CreateAIChatMessage) Execute(
	ctx context.Context,
	in CreateAIChatMessageInput,
) (*entity.AIChatAnswer, error) {
	if err := uc.v.Validate(in); err != nil {
		return nil, errs.New(err)
	}

	messageResponse, err := uc.generateMessageResponse(ctx, in)
	if err != nil {
		return nil, errs.New(err)
	}

	var aiChatAnswer *entity.AIChatAnswer
	err = uc.tx.Do(ctx, func(ctx context.Context) error {
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
				Message:         messageResponse,
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

	return aiChatAnswer, nil
}

func (uc *CreateAIChatMessage) generateMessageResponse(
	ctx context.Context,
	in CreateAIChatMessageInput,
) (string, error) {
	const FunctionGetChatHistory = "get_chat_history"
	const FunctionGetUserFinancialData = "get_user_financial_data"

	params := openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(in.Message),
		}),
		Tools: openai.F([]openai.ChatCompletionToolParam{
			{
				Type: openai.F(openai.ChatCompletionToolTypeFunction),
				Function: openai.F(openai.FunctionDefinitionParam{
					Name: openai.String(FunctionGetChatHistory),
					Description: openai.String(
						"Get chat history, including previous messages and responses",
					),
				}),
			},
			{
				Type: openai.F(openai.ChatCompletionToolTypeFunction),
				Function: openai.F(openai.FunctionDefinitionParam{
					Name: openai.String(FunctionGetUserFinancialData),
					Description: openai.String(
						"Get user financial data, such as transactions, accounts, and budgets",
					),
				}),
			},
		}),
		Seed:  openai.Int(0),
		Model: openai.F(openai.ChatModelO3Mini),
	}

	completion, err := uc.oa.Client.Chat.Completions.New(ctx, params)
	if err != nil {
		return "", errs.New(err)
	}

	toolCalls := completion.Choices[0].Message.ToolCalls
	if len(toolCalls) == 0 {
		return completion.Choices[0].Message.Content, nil
	}

	params.Messages.Value = append(
		params.Messages.Value,
		completion.Choices[0].Message,
	)

	g, gCtx := errgroup.WithContext(ctx)
	mu := sync.Mutex{}
	for _, toolCall := range toolCalls {
		g.Go(func() error {
			switch toolCall.Function.Name {
			case FunctionGetChatHistory:
				chatHistory, err := uc.getChatHistory(gCtx, in.AIChatID)
				if err != nil {
					return errs.New(err)
				}
				mu.Lock()
				params.Messages.Value = append(
					params.Messages.Value,
					openai.ToolMessage(toolCall.ID, chatHistory),
				)
				mu.Unlock()

			case FunctionGetUserFinancialData:
				userFinancialData, err := uc.getUserFinancialData(
					gCtx,
					in.UserID,
					in.Message,
				)
				if err != nil {
					return errs.New(err)
				}
				mu.Lock()
				params.Messages.Value = append(
					params.Messages.Value,
					openai.ToolMessage(toolCall.ID, userFinancialData),
				)
				mu.Unlock()

			default:
				return errs.New("unknown tool call function")
			}

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return "", errs.New(err)
	}

	completion, err = uc.oa.Client.Chat.Completions.New(ctx, params)
	if err != nil {
		panic(err)
	}

	return completion.Choices[0].Message.Content, nil
}

func (uc *CreateAIChatMessage) getChatHistory(
	ctx context.Context,
	aiChatID uuid.UUID,
) (string, error) {
	chatRecentHistory, err := uc.acr.ListAIChatMessagesAndAnswers(
		ctx,
		repo.ListAIChatMessagesAndAnswersParams{
			AiChatID: aiChatID,
			Limit:    6,
		},
	)
	if err != nil {
		return "", errs.New(err)
	}
	jsonData, err := json.Marshal(chatRecentHistory)
	if err != nil {
		return "", errs.New(err)
	}
	return string(jsonData), nil
}

func (uc *CreateAIChatMessage) getUserFinancialData(
	ctx context.Context,
	userID uuid.UUID,
	message string,
) (string, error) {
	sqlQuery, err := uc.generateSQLQuery(
		ctx,
		userID,
		message,
	)
	if err != nil {
		return "", errs.New(err)
	}
	return uc.executeSQLQuery(ctx, sqlQuery)
}

func (uc *CreateAIChatMessage) generateSQLQuery(
	ctx context.Context,
	userID uuid.UUID,
	message string,
) (string, error) {
	g, gCtx := errgroup.WithContext(ctx)
	var (
		paymentMethods []entity.PaymentMethod
		categories     []entity.TransactionCategory
		institutions   []entity.Institution
	)

	g.Go(func() (err error) {
		paymentMethods, err = uc.pmr.ListPaymentMethods(gCtx)
		return err
	})

	g.Go(func() (err error) {
		categories, err = uc.tcr.ListTransactionCategories(gCtx)
		return err
	})

	g.Go(func() (err error) {
		institutions, err = uc.ir.ListInstitutions(
			gCtx,
			repo.WithInstitutionUsers([]uuid.UUID{userID}),
		)
		return err
	})

	if err := g.Wait(); err != nil {
		return "", errs.New(err)
	}

	paymentMethodsMap := map[string]any{}
	for _, pm := range paymentMethods {
		paymentMethodsMap["id"] = pm.ID.String()
		paymentMethodsMap["name"] = pm.Name
	}
	paymentMethodsJSON, err := json.Marshal(paymentMethodsMap)
	if err != nil {
		return "", errs.New(err)
	}

	categoriesMap := map[string]any{}
	for _, c := range categories {
		categoriesMap["id"] = c.ID.String()
		categoriesMap["name"] = c.Name
	}
	categoriesJSON, err := json.Marshal(categoriesMap)
	if err != nil {
		return "", errs.New(err)
	}

	institutionsMap := map[string]any{}
	for _, i := range institutions {
		institutionsMap["id"] = i.ID.String()
		institutionsMap["name"] = i.Name
	}
	institutionsJSON, err := json.Marshal(institutionsMap)
	if err != nil {
		return "", errs.New(err)
	}

	systemPrompt := fmt.Sprintf(
		`You're an expert in SQL who generates PostgreSQL queries based on financial questions.
Here's the database schema (simplified Prisma format):

model Transaction {
  id String @id @db.Uuid
  name String
  amount BigInt           // Stored in cents, divide by 100 for display
  is_ignored Boolean      // If true, ignore this transaction in calculations
  date DateTime
  deleted_at DateTime?
  payment_method_id String? @db.Uuid
  user_id String @db.Uuid
  category_id String @db.Uuid
  account_id String? @db.Uuid
  institution_id String? @db.Uuid

  @@map("transactions")
}

model AccountBalance {
  id String @id @db.Uuid
  amount BigInt           // Stored in cents, divide by 100 for display
  created_at DateTime
  deleted_at DateTime?
  account_id String @db.Uuid

  @@map("account_balances")
}

model Account {
  id String @id @db.Uuid
  name String
  type String  // 'BANK' or 'CREDIT'
  deleted_at DateTime?
  user_institution_id String @db.Uuid

  @@map("accounts")
}

model UserInstitution {
  id String @id @db.Uuid
  user_id String @db.Uuid
  institution_id String @db.Uuid

  @@map("user_institutions")
}

model Budget {
  id String @id @db.Uuid
  amount BigInt           // Stored in cents, divide by 100 for display
  date DateTime
  deleted_at DateTime?
  user_id String @db.Uuid

  @@map("budgets")
}

model BudgetCategory {
  id String @id @db.Uuid
  amount BigInt           // Stored in cents, divide by 100 for display
  deleted_at DateTime?
  budget_id String @db.Uuid
  category_id String @db.Uuid

  @@map("budget_categories")
}

model TransactionCategory {
  id String @id @db.Uuid
  name String
  deleted_at DateTime?

  @@map("transaction_categories")
}

Here is the payment_methods table data, represented as JSON:
%s

Here is the transaction_categories table data, represented as JSON:
%s

Here is the institutions table data, with only institutions that the user has an account, represented as JSON:
%s

IMPORTANT RULES:
1. ALWAYS filter for the specific user: user_id = '%s'
2. ALWAYS exclude deleted records: deleted_at IS NULL
3. Return only the minimal necessary data to answer the question
4. Format amounts as dollars using amount::float/100 in your SELECT
5. Include appropriate JOINs to get related data when needed
6. You cannot write destructive queries (INSERT, UPDATE, DELETE)
7. Instead of returning columns with default or technical names, use the AS keyword to create clear, self-explanatory column names in your query results
8. Respond ONLY with a single SQL query, no explanations or comments.`,
		string(paymentMethodsJSON),
		string(categoriesJSON),
		string(institutionsJSON),
		userID,
	)

	userPrompt := fmt.Sprintf(
		"Generate a PostgreSQL query for this question: %s",
		message,
	)

	chatCompletion, err := uc.oa.Client.Chat.Completions.New(
		ctx,
		openai.ChatCompletionNewParams{
			Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
				openai.SystemMessage(systemPrompt),
				openai.UserMessage(userPrompt),
			}),
			Model: openai.F(openai.ChatModelO3Mini),
		},
	)
	if err != nil {
		return "", errs.New(err)
	}

	return chatCompletion.Choices[0].Message.Content, nil
}

func (uc *CreateAIChatMessage) executeSQLQuery(
	ctx context.Context,
	sqlQuery string,
) (string, error) {
	dest := map[string]any{}
	if err := uc.db.ScanRaw(ctx, sqlQuery, &dest); err != nil {
		return "", errs.New(err)
	}
	jsonData, err := json.Marshal(dest)
	if err != nil {
		return "", errs.New(err)
	}
	return string(jsonData), nil
}
