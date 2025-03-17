package usecase

import (
	"context"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
	"github.com/google/uuid"
)

type DeleteAIChat struct {
	acr repo.AIChatRepo
}

func NewDeleteAIChat(
	acr repo.AIChatRepo,
) *DeleteAIChat {
	return &DeleteAIChat{
		acr: acr,
	}
}

func (uc *DeleteAIChat) Execute(
	ctx context.Context,
	id uuid.UUID,
) error {
	aiChat, err := uc.acr.GetAIChatByID(ctx, id)
	if err != nil {
		return errs.New(err)
	}
	if aiChat == nil {
		return errs.ErrAIChatNotFound
	}

	if err := uc.acr.DeleteAIChat(ctx, id); err != nil {
		return errs.New(aiChat)
	}

	return nil
}
