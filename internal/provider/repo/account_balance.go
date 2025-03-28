package repo

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type AccountBalanceOptions struct {
	InstitutionIDs []uuid.UUID `json:"institution_ids"`
}

type AccountBalanceRepo interface {
	CreateAccountBalances(
		ctx context.Context,
		params []CreateAccountBalancesParams,
	) error
	GetUserBalanceOnDate(
		ctx context.Context,
		userID uuid.UUID,
		date time.Time,
		opts ...AccountBalanceOptions,
	) (int64, error)
}
