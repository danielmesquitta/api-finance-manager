package aichat

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

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

type GenerateAIChatMessageUseCase struct {
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

func NewGenerateAIChatMessageUseCase(
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
) *GenerateAIChatMessageUseCase {
	return &GenerateAIChatMessageUseCase{
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

type GenerateAIChatMessageUseCaseInput struct {
	UserID   uuid.UUID   `json:"-"       validate:"required"`
	AIChatID uuid.UUID   `json:"-"       validate:"required"`
	Message  string      `json:"message" validate:"required,max=512"`
	Tier     entity.Tier `json:"-"       validate:"required,oneof=TRIAL PREMIUM"`
}

func (uc *GenerateAIChatMessageUseCase) Execute(
	ctx context.Context,
	in GenerateAIChatMessageUseCaseInput,
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
		aiChatMessage, err := uc.acmr.GenerateAIChatMessage(
			ctx,
			repo.GenerateAIChatMessageParams{
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

func (uc *GenerateAIChatMessageUseCase) generateMessageResponse(
	ctx context.Context,
	in GenerateAIChatMessageUseCaseInput,
) (string, error) {
	const FunctionGetChatHistory = "get_chat_history"
	const FunctionGetUserFinancialData = "get_user_financial_data"

	params := openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(
				fmt.Sprintf(
					"Today is %s and you are a financial specialist with access to all of your clients' financial data. Drawing on your expertise in financial planning, please provide insightful answers and actionable advice.",
					time.Now().Format(time.RFC3339),
				),
			),
			openai.UserMessage(in.Message),
		},
		Tools: []openai.ChatCompletionToolParam{
			{
				Function: openai.FunctionDefinitionParam{
					Name: FunctionGetChatHistory,
					Description: openai.String(
						"Get chat history, including previous messages and responses",
					),
				},
			},
			{
				Function: openai.FunctionDefinitionParam{
					Name: FunctionGetUserFinancialData,
					Description: openai.String(
						"Get user financial data, such as transactions, accounts, and budgets",
					),
				},
			},
		},
		Model: openai.ChatModelO3Mini,
	}

	completion, err := uc.oa.Client.Chat.Completions.New(ctx, params)
	if err != nil {
		return "", errs.New(err)
	}

	toolCalls := completion.Choices[0].Message.ToolCalls
	if len(toolCalls) == 0 {
		return completion.Choices[0].Message.Content, nil
	}

	params.Messages = append(
		params.Messages,
		completion.Choices[0].Message.ToParam(),
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
				params.Messages = append(
					params.Messages,
					openai.ToolMessage(chatHistory, toolCall.ID),
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
				params.Messages = append(
					params.Messages,
					openai.ToolMessage(userFinancialData, toolCall.ID),
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

func (uc *GenerateAIChatMessageUseCase) getChatHistory(
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

func (uc *GenerateAIChatMessageUseCase) getUserFinancialData(
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

func (uc *GenerateAIChatMessageUseCase) generateSQLQuery(
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

	paymentMethodMaps := []map[string]any{}
	for _, pm := range paymentMethods {
		paymentMethodMaps = append(paymentMethodMaps, map[string]any{
			"id":   pm.ID.String(),
			"name": pm.Name,
		})
	}
	paymentMethodsJSON, err := json.Marshal(paymentMethodMaps)
	if err != nil {
		return "", errs.New(err)
	}

	categoryMaps := []map[string]any{}
	for _, c := range categories {
		categoryMaps = append(categoryMaps, map[string]any{
			"id":   c.ID.String(),
			"name": c.Name,
		})
	}
	categoriesJSON, err := json.Marshal(categoryMaps)
	if err != nil {
		return "", errs.New(err)
	}

	institutionMaps := []map[string]any{}
	for _, i := range institutions {
		institutionMaps = append(institutionMaps, map[string]any{
			"id":   i.ID.String(),
			"name": i.Name,
		})
	}
	institutionsJSON, err := json.Marshal(institutionMaps)
	if err != nil {
		return "", errs.New(err)
	}

	systemPrompt := fmt.Sprintf(
		`You're an expert in SQL who generates PostgreSQL queries based on financial questions.
Here's the database schema (simplified Prisma format):

model Transaction {
  id String @id @db.Uuid
  name String
  amount BigInt           // Stored in cents, divide by 100 for display, negative for spent and positive for earned
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

CONSIDERATIONS:
- Budgets are set for a specific month and repeat. If a new budget isn't created the next month, the most recent one applies. For example, to get the latest budget for a user:

SELECT *
FROM budgets
WHERE user_id = $1
  AND date <= $2
  AND deleted_at IS NULL
ORDER BY date DESC
LIMIT 1;

IMPORTANT RULES:
1. ALWAYS filter for the specific user: user_id = '%s'
2. ALWAYS exclude deleted records: deleted_at IS NULL
3. Return only the minimal necessary data to answer the question
4. Format amounts as dollars using amount::float/100 in your SELECT
5. Include appropriate JOINs to get related data when needed
6. You cannot write destructive queries (INSERT, UPDATE, DELETE)
7. Use the AS keyword to create clear, self-explanatory column names in your query results
8. Respond ONLY with a single SQL query, no explanations or comments
9. Favor searching for one or multiple transaction_categories, payment_methods, and institutions by their IDs instead of doing full-text searches.`,
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
			Messages: []openai.ChatCompletionMessageParamUnion{
				openai.SystemMessage(systemPrompt),
				openai.UserMessage(userPrompt),
			},
			Model: openai.ChatModelO3Mini,
		},
	)
	if err != nil {
		return "", errs.New(err)
	}

	return chatCompletion.Choices[0].Message.Content, nil
}

func (uc *GenerateAIChatMessageUseCase) executeSQLQuery(
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
