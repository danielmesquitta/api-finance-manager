package repo

import (
	"context"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/google/uuid"
)

type AccountRepo interface {
	ListAccountsByUserID(
		ctx context.Context,
		userID uuid.UUID,
	) ([]entity.Account, error)
	ListAccountsByExternalIDs(
		ctx context.Context,
		externalIDs []string,
	) ([]entity.Account, error)
	CreateAccounts(
		ctx context.Context,
		params []CreateAccountsParams,
	) error
}
