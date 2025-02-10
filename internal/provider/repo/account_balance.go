package repo

import (
	"context"
)

type AccountBalanceRepo interface {
	CreateAccountBalances(
		ctx context.Context,
		params []CreateAccountBalancesParams,
	) error
	GetUserBalanceOnDate(
		ctx context.Context,
		arg GetUserBalanceOnDateParams,
	) (int64, error)
}
