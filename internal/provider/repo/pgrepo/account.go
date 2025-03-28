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

type AccountRepo struct {
	db *db.DB
	qb *query.QueryBuilder
}

func NewAccountRepo(
	db *db.DB,
	qb *query.QueryBuilder,
) *AccountRepo {
	return &AccountRepo{
		db: db,
		qb: qb,
	}
}

func (r *AccountRepo) ListAccounts(
	ctx context.Context,
	opts ...repo.AccountOptions,
) ([]entity.Account, error) {
	return r.qb.ListAccounts(ctx, opts...)
}

func (r *AccountRepo) ListFullAccounts(
	ctx context.Context,
	opts ...repo.AccountOptions,
) ([]entity.FullAccount, error) {
	return r.qb.ListFullAccounts(ctx, opts...)
}

func (r *AccountRepo) CreateAccounts(
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

var _ repo.AccountRepo = (*AccountRepo)(nil)
