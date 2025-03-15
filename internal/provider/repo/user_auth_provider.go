package repo

import (
	"context"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
)

type UserAuthProviderRepo interface {
	CreateUserAuthProvider(
		ctx context.Context,
		params CreateUserAuthProviderParams,
	) error
	UpdateUserAuthProvider(
		ctx context.Context,
		params UpdateUserAuthProviderParams,
	) error
	GetUserAuthProvider(
		ctx context.Context,
		params GetUserAuthProviderParams,
	) (*entity.UserAuthProvider, error)
}
