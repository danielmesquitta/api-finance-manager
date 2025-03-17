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

type CreateAIChatInput struct {
	UserID uuid.UUID   `json:"-"`
	Tier   entity.Tier `json:"-"`
}

func (uc *CreateAIChat) Execute(
	ctx context.Context,
	in CreateAIChatInput,
) (*entity.AIChat, error) {
	if in.Tier == entity.TierFree {
		return nil, errs.ErrUnauthorizedTier
	}

	latestAIChat, err := uc.acr.GetLatestAIChatByUserID(ctx, in.UserID)
	if err != nil {
		return nil, errs.New(err)
	}

	if latestAIChat != nil && !latestAIChat.HasMessages {
		return &latestAIChat.AIChat, nil
	}

	aiChat, err := uc.acr.CreateAIChat(ctx, in.UserID)
	if err != nil {
		return nil, errs.New(aiChat)
	}

	return aiChat, nil
}
