package usecase

import (
	"context"
	"time"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/validator"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
	"github.com/jinzhu/copier"
)

type CreateUser struct {
	v  validator.Validator
	ur repo.UserRepo
}

func NewCreateUser(
	v validator.Validator,
	ur repo.UserRepo,
) *CreateUser {
	return &CreateUser{
		v:  v,
		ur: ur,
	}
}

type CreateUserInput struct {
	Name  string `validate:"required"`
	Email string `validate:"required,email"`
}

func (uc *CreateUser) Execute(ctx context.Context, in CreateUserInput) (*entity.User, error) {
	if err := uc.v.Validate(in); err != nil {
		return nil, errs.New(err)
	}

	userWithSameEmail, err := uc.ur.GetUserByEmail(ctx, in.Email)
	if err != nil {
		return nil, errs.New(err)
	}
	if userWithSameEmail != nil {
		return nil, errs.ErrUserAlreadyRegistered
	}

	params := repo.CreateUserParams{}
	if err := copier.Copy(&params, in); err != nil {
		return nil, errs.New(err)
	}

	params.Tier = entity.TierTRIAL

	twoWeeksFromNow := time.Now().AddDate(0, 0, 14)
	params.SubscriptionExpiresAt = twoWeeksFromNow

	user, err := uc.ur.CreateUser(ctx, params)
	if err != nil {
		return nil, errs.New(err)
	}

	if user == nil {
		return nil, errs.ErrUserNotFound
	}

	return user, nil
}
