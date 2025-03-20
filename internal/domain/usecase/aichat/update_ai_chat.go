package aichat

import (
	"context"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/validator"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

type UpdateAIChatUseCase struct {
	v   *validator.Validator
	acr repo.AIChatRepo
}

func NewUpdateAIChatUseCase(
	v *validator.Validator,
	acr repo.AIChatRepo,
) *UpdateAIChatUseCase {
	return &UpdateAIChatUseCase{
		v:   v,
		acr: acr,
	}
}

type UpdateAIChatUseCaseInput struct {
	ID    uuid.UUID   `json:"-"     validate:"required"`
	Title string      `json:"title" validate:"required"`
	Tier  entity.Tier `json:"-"     validate:"required,oneof=TRIAL PREMIUM"`
}

func (uc *UpdateAIChatUseCase) Execute(
	ctx context.Context,
	in UpdateAIChatUseCaseInput,
) error {
	if err := uc.v.Validate(in); err != nil {
		return errs.New(err)
	}

	aiChat, err := uc.acr.GetAIChatByID(ctx, in.ID)
	if err != nil {
		return errs.New(err)
	}
	if aiChat == nil {
		return errs.ErrAIChatNotFound
	}

	var params repo.UpdateAIChatParams
	if err := copier.Copy(&params, in); err != nil {
		return errs.New(err)
	}

	if err := uc.acr.UpdateAIChat(ctx, params); err != nil {
		return errs.New(aiChat)
	}

	return nil
}
