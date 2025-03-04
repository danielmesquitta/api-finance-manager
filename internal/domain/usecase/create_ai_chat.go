package usecase

import (
	"context"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
	"github.com/google/uuid"
)

type CreateAIChat struct {
	acr repo.AIChatRepo
}

func NewCreateAIChat(
	acr repo.AIChatRepo,
) *CreateAIChat {
	return &CreateAIChat{
		acr: acr,
	}
}

func (uc *CreateAIChat) Execute(
	ctx context.Context,
	userID uuid.UUID,
) (*entity.AIChat, error) {
	latestAIChat, err := uc.acr.GetLatestAIChatByUserID(ctx, userID)
	if err != nil {
		return nil, errs.New(err)
	}

	if latestAIChat != nil && !latestAIChat.HasMessages {
		return &latestAIChat.AIChat, nil
	}

	aiChat, err := uc.acr.CreateAIChat(ctx, userID)
	if err != nil {
		return nil, errs.New(aiChat)
	}

	return aiChat, nil
}
