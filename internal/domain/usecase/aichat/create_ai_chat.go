package aichat

import (
	"context"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/validator"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
	"github.com/google/uuid"
)

type CreateAIChatUseCase struct {
	v   *validator.Validator
	acr repo.AIChatRepo
}

func NewCreateAIChatUseCase(
	v *validator.Validator,
	acr repo.AIChatRepo,
) *CreateAIChatUseCase {
	return &CreateAIChatUseCase{
		v:   v,
		acr: acr,
	}
}

type CreateAIChatUseCaseInput struct {
	UserID uuid.UUID   `json:"-"`
	Tier   entity.Tier `json:"-" validate:"required,oneof=TRIAL PREMIUM"`
}

func (uc *CreateAIChatUseCase) Execute(
	ctx context.Context,
	in CreateAIChatUseCaseInput,
) (*entity.AIChat, error) {
	if err := uc.v.Validate(in); err != nil {
		return nil, errs.New(err)
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
