package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/validator"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
)

type CreateFeedback struct {
	v  *validator.Validator
	fr repo.FeedbackRepo
}

func NewCreateFeedback(
	v *validator.Validator,
	fr repo.FeedbackRepo,
) *CreateFeedback {
	return &CreateFeedback{
		v:  v,
		fr: fr,
	}
}

type CreateFeedbackInput struct {
	Message string    `json:"message" validate:"required"`
	UserID  uuid.UUID `json:"-"       validate:"required"`
}

func (uc *CreateFeedback) Execute(
	ctx context.Context,
	in CreateFeedbackInput,
) error {
	if err := uc.v.Validate(in); err != nil {
		return err
	}

	var params repo.CreateFeedbackParams
	if err := copier.Copy(&params, in); err != nil {
		return errs.New(err)
	}

	if err := uc.fr.CreateFeedback(ctx, params); err != nil {
		return errs.New(err)
	}

	return nil
}
