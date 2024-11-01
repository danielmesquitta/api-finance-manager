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
	id string,
) (*entity.User, error) {
	uuidID, err := uuid.Parse(id)
	if err != nil {
		return nil, errs.New(err)
	}

	user, err := uc.ur.GetUserByID(ctx, uuidID)
	if err != nil {
		return nil, errs.New(err)
	}
	if user == nil {
		return nil, errs.ErrUserNotFound
	}

	return user, nil
}
