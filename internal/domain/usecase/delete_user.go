package usecase

import (
	"context"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/hash"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
	"github.com/google/uuid"
)

type DeleteUser struct {
	h  *hash.Hasher
	ur repo.UserRepo
}

func NewDeleteUser(
	h *hash.Hasher,
	ur repo.UserRepo,
) *DeleteUser {
	return &DeleteUser{
		h:  h,
		ur: ur,
	}
}

func (uc *DeleteUser) Execute(ctx context.Context, id uuid.UUID) error {
	user, err := uc.ur.GetUserByID(ctx, id)
	if err != nil {
		return errs.New(err)
	}
	if user == nil {
		return errs.ErrUserNotFound
	}

	hashedName, err := uc.h.Hash(user.Name)
	if err != nil {
		return errs.New(err)
	}

	hashedEmail, err := uc.h.Hash(user.Email)
	if err != nil {
		return errs.New(err)
	}

	params := repo.DeleteUserParams{
		ID:    id,
		Name:  hashedName,
		Email: hashedEmail,
	}

	if err := uc.ur.DeleteUser(ctx, params); err != nil {
		return errs.New(err)
	}

	return nil
}
