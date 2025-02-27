package pgrepo

import (
	"context"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/db"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/db/query"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/db/sqlc"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
	"github.com/jinzhu/copier"
)

type AccountPgRepo struct {
	db *db.DB
	qb *query.QueryBuilder
}

func NewAccountPgRepo(
	db *db.DB,
	qb *query.QueryBuilder,
) *AccountPgRepo {
	return &AccountPgRepo{
		db: db,
		qb: qb,
	}
}

func (r *AccountPgRepo) ListAccounts(
	ctx context.Context,
	opts ...repo.AccountOption,
) ([]entity.Account, error) {
	return r.qb.ListAccounts(ctx, opts...)
}

func (r *AccountPgRepo) ListFullAccounts(
	ctx context.Context,
	opts ...repo.AccountOption,
) ([]entity.FullAccount, error) {
	return r.qb.ListFullAccounts(ctx, opts...)
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

var _ repo.AccountRepo = (*AccountPgRepo)(nil)
