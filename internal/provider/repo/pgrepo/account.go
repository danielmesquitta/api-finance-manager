package pgrepo

import (
	"context"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/db"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/db/sqlc"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

type AccountPgRepo struct {
	db *db.DB
}

func NewAccountPgRepo(db *db.DB) *AccountPgRepo {
	return &AccountPgRepo{
		db: db,
	}
}

func (r *AccountPgRepo) ListAccountsByUserID(
	ctx context.Context,
	userID uuid.UUID,
) ([]entity.Account, error) {
	accounts, err := r.db.ListAccountsByUserID(ctx, userID)
	if err != nil {
		return nil, errs.New(err)
	}

	results := []entity.Account{}
	if err := copier.Copy(&results, accounts); err != nil {
		return nil, errs.New(err)
	}

	return results, nil
}

func (r *AccountPgRepo) CreateAccounts(
	ctx context.Context,
	params []repo.CreateAccountsParams,
) error {
	dbParams := make([]sqlc.CreateAccountsParams, len(params))
	if err := copier.Copy(&dbParams, params); err != nil {
		return errs.New(err)
	}

	tx := r.db.UseTx(ctx)
	_, err := tx.CreateAccounts(ctx, dbParams)
	if err != nil {
		return errs.New(err)
	}

	return nil
}

var _ repo.AccountRepo = &AccountPgRepo{}
