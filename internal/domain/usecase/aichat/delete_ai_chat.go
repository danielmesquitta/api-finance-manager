package aichat

import (
	"context"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
	"github.com/google/uuid"
)

type DeleteAIChatUseCase struct {
	acr repo.AIChatRepo
}

func NewDeleteAIChatUseCase(
	acr repo.AIChatRepo,
) *DeleteAIChatUseCase {
	return &DeleteAIChatUseCase{
		acr: acr,
	}
}

func (uc *DeleteAIChatUseCase) Execute(
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
