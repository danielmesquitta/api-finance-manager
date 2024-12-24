package usecase

import (
	"context"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
	"github.com/google/uuid"
)

type GetUser struct {
	ur repo.UserRepo
}

func NewGetUser(
	ur repo.UserRepo,
) *GetUser {
	return &GetUser{
		ur: ur,
	}
}

func (uc *GetUser) Execute(
	ctx context.Context,
	id uuid.UUID,
) (*entity.User, error) {
	user, err := uc.ur.GetUserByID(ctx, id)
	if err != nil {
		return nil, errs.New(err)
	}
	if user == nil {
		return nil, errs.ErrUserNotFound
	}

	return user, nil
}
