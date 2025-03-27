package repo

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type AccountBalanceOptions struct {
	InstitutionIDs []uuid.UUID `json:"institution_ids"`
}

type AccountBalanceOption func(*AccountBalanceOptions)

func WithAAccountBalanceInstitutions(
	institutionIDs []uuid.UUID,
) AccountBalanceOption {
	return func(o *AccountBalanceOptions) {
		o.InstitutionIDs = institutionIDs
	}
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
		options ...AccountBalanceOption,
	) (int64, error)
}
