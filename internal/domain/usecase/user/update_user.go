package user

import (
	"context"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/validator"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

type UpdateUserUseCase struct {
	v  *validator.Validator
	ur repo.UserRepo
}

func NewUpdateUserUseCase(
	v *validator.Validator,
	ur repo.UserRepo,
) *UpdateUserUseCase {
	return &UpdateUserUseCase{
		v:  v,
		ur: ur,
	}
}

type UpdateUserUseCaseInput struct {
	ID    uuid.UUID `json:"-"`
	Name  string    `json:"name"`
	Email string    `json:"email" validate:"omitempty,email"`
}

func (uc *UpdateUserUseCase) Execute(
	ctx context.Context,
	in UpdateUserUseCaseInput,
) error {
	if err := uc.v.Validate(in); err != nil {
		return errs.New(err)
	}

	user, err := uc.ur.GetUserByID(ctx, in.ID)
	if err != nil {
		return errs.New(err)
	}
	if user == nil {
		return errs.ErrUserNotFound
	}

	params := repo.UpdateUserParams{}
	if err := copier.Copy(&params, user); err != nil {
		return errs.New(err)
	}

	if err := copier.CopyWithOption(
		&params,
		in,
		copier.Option{IgnoreEmpty: true},
	); err != nil {
		return errs.New(err)
	}

	if _, err := uc.ur.UpdateUser(ctx, params); err != nil {
		return errs.New(err)
	}

	return nil
}
