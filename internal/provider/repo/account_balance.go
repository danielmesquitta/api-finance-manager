package repo

import (
	"context"

	"github.com/google/uuid"
)

type AccountBalanceRepo interface {
	CreateAccountBalance(
		ctx context.Context,
		params CreateAccountBalanceParams,
	) error
	GetUserBalance(
		ctx context.Context,
		userID uuid.UUID,
	) (int64, error)
	GetUserBalanceOnDate(
		ctx context.Context,
		arg GetUserBalanceOnDateParams,
	) (int64, error)
}
