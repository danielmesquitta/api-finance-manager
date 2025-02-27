package pgrepo

import (
	"context"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/db"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/db/sqlc"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
	"github.com/jinzhu/copier"
)

type AccountBalancePgRepo struct {
	db *db.DB
}

func NewAccountBalancePgRepo(db *db.DB) *AccountBalancePgRepo {
	return &AccountBalancePgRepo{
		db: db,
	}
}

func (r *AccountBalancePgRepo) CreateAccountBalances(
	ctx context.Context,
	params []repo.CreateAccountBalancesParams,
) error {
	dbParams := make([]sqlc.CreateAccountBalancesParams, len(params))
	if err := copier.Copy(&dbParams, params); err != nil {
		return errs.New(err)
	}

	tx := r.db.UseTx(ctx)
	if _, err := tx.CreateAccountBalances(ctx, dbParams); err != nil {
		return errs.New(err)
	}

	return nil
}

func (r *AccountBalancePgRepo) GetUserBalanceOnDate(
	ctx context.Context,
	params repo.GetUserBalanceOnDateParams,
) (int64, error) {
	dbParams := sqlc.GetUserBalanceOnDateParams{}
	if err := copier.Copy(&dbParams, params); err != nil {
		return 0, errs.New(err)
	}

	return r.db.GetUserBalanceOnDate(ctx, dbParams)
}

var _ repo.AccountBalanceRepo = (*AccountBalancePgRepo)(nil)
