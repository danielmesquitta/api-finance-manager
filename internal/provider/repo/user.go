package repo

import (
	"context"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/db/sqlc"
	"github.com/google/uuid"
)

type CreateUserParams struct {
	sqlc.CreateUserParams
}

type UpdateUserParams struct {
	sqlc.UpdateUserParams
}

type UserRepo interface {
	CreateUser(
		ctx context.Context,
		params CreateUserParams,
	) (*entity.User, error)
	UpdateUser(
		ctx context.Context,
		params UpdateUserParams,
	) (*entity.User, error)
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
}
