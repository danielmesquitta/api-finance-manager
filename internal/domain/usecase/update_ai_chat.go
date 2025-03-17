package usecase

import (
	"context"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/validator"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

type UpdateAIChat struct {
	v   *validator.Validator
	acr repo.AIChatRepo
}

func NewUpdateAIChat(
	v *validator.Validator,
	acr repo.AIChatRepo,
) *UpdateAIChat {
	return &UpdateAIChat{
		v:   v,
		acr: acr,
	}
}

type UpdateAIChatInput struct {
	ID    uuid.UUID   `json:"-"     validate:"required"`
	Title string      `json:"title" validate:"required"`
	Tier  entity.Tier `json:"-"     validate:"required,oneof=TRIAL PREMIUM"`
}

func (uc *UpdateAIChat) Execute(
	ctx context.Context,
	in UpdateAIChatInput,
) error {
	if err := uc.v.Validate(in); err != nil {
		return errs.New(err)
	}

	aiChat, err := uc.acr.GetAIChat(ctx, in.ID)
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
