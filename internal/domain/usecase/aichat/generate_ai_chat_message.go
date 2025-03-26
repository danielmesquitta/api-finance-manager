package aichat

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/tx"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/validator"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/db"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/gpt"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
)

type GenerateAIChatMessageUseCase struct {
	v     *validator.Validator
	tx    tx.TX
	db    *db.DB
	gp    gpt.GPT
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
	gp gpt.GPT,
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
		gp:    gp,
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
		aiChat            *entity.AIChat
		chatRecentHistory []entity.AIChatMessageAndAnswer
	)

	g.Go(func() error {
		aiChat, err = uc.acr.GetAIChatByID(subCtx, in.AIChatID)
		return err
	})

	g.Go(func() error {
		chatRecentHistory, err = uc.acr.ListAIChatMessagesAndAnswers(
			ctx,
			repo.ListAIChatMessagesAndAnswersParams{
				AiChatID: in.AIChatID,
				Limit:    6,
			},
		)
		return err
	})

	if err := g.Wait(); err != nil {
		return nil, errs.New(err)
	}

	if aiChat == nil || aiChat.UserID != in.UserID {
		return nil, errs.ErrAIChatNotFound
	}

	g, subCtx = errgroup.WithContext(ctx)
	var title, answer string

	if aiChat.Title == nil {
		g.Go(func() error {
			title, err = uc.generateAIChatTitle(subCtx, in.Message)
			return err
		})
	}

	g.Go(func() error {
		answer, err = uc.generateAIChatAnswer(
			ctx,
			in,
			chatRecentHistory,
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

func (uc *GenerateAIChatMessageUseCase) generateAIChatAnswer(
	ctx context.Context,
	in GenerateAIChatMessageUseCaseInput,
	chatHistory []entity.AIChatMessageAndAnswer,
) (string, error) {
	systemMessage := fmt.Sprintf(
		"Today is %s and you are a financial specialist with access to all of your clients' financial data. Drawing on your expertise in financial planning, please provide insightful answers and actionable advice.",
		time.Now().Format(time.RFC3339),
	)

	messages := []gpt.Message{
		{
			Role:    gpt.RoleSystem,
			Content: systemMessage,
		},
	}

	prevMsgs := []gpt.Message{}
	mapAuthorToRole := map[entity.AIChatMessageAuthor]gpt.Role{
		entity.AIChatMessageAuthorUser: gpt.RoleUser,
		entity.AIChatMessageAuthorAI:   gpt.RoleAssistant,
	}
	for _, msg := range chatHistory {
		prevMsgs = append(prevMsgs, gpt.Message{
			Role:    mapAuthorToRole[msg.Author],
			Content: msg.Message,
		})
	}

	messages = append(messages, prevMsgs...)

	messages = append(messages, gpt.Message{
		Role:    gpt.RoleUser,
		Content: in.Message,
	})

	tools := []gpt.Tool{
		{
			Name:        "get_user_financial_data",
			Description: "Get user financial data, such as transactions, accounts, and budgets",
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
