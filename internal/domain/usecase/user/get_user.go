package user

import (
	"context"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
	"github.com/google/uuid"
)

type GetUserUseCase struct {
	ur repo.UserRepo
}

func NewGetUserUseCase(
	ur repo.UserRepo,
) *GetUserUseCase {
	return &GetUserUseCase{
		ur: ur,
	}
}

func (uc *GetUserUseCase) Execute(
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
