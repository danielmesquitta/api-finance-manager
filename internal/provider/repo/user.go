package repo

import (
	"context"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/google/uuid"
)

type PremiumUserOptions struct {
	Limit  uint `json:"-"`
	Offset uint `json:"-"`
}

type PremiumUserOption func(*PremiumUserOptions)

func WithPremiumUserPagination(
	limit uint,
	offset uint,
) PremiumUserOption {
	return func(o *PremiumUserOptions) {
		o.Limit = limit
		o.Offset = offset
	}
}

type UserRepo interface {
	CreateUser(
		ctx context.Context,
		params CreateUserParams,
	) (*entity.User, error)
	DeleteUser(
		ctx context.Context,
		params DeleteUserParams,
	) error
	DestroyUser(
		ctx context.Context,
		id uuid.UUID,
	) error
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	UpdateUser(
		ctx context.Context,
		params UpdateUserParams,
	) (*entity.User, error)
	UpdateUserSynchronizedAt(
		ctx context.Context,
		arg UpdateUserSynchronizedAtParams,
	) error
}
