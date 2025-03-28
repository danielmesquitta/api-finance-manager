package pgrepo

import (
	"context"
	"time"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/db"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/db/sqlc"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

type AccountBalanceRepo struct {
	db *db.DB
}

func NewAccountBalanceRepo(db *db.DB) *AccountBalanceRepo {
	return &AccountBalanceRepo{
		db: db,
	}
}

func (r *AccountBalanceRepo) CreateAccountBalances(
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

func (r *AccountBalanceRepo) GetUserBalanceOnDate(
	ctx context.Context,
	userID uuid.UUID,
	date time.Time,
	opts ...repo.AccountBalanceOptions,
) (int64, error) {
	return r.db.GetUserBalanceOnDate(ctx, userID, date, opts...)
}

var _ repo.AccountBalanceRepo = (*AccountBalanceRepo)(nil)
