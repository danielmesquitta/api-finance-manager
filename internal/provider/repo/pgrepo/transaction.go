package pgrepo

import (
	"context"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/db"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/db/query"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/db/sqlc"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

type TransactionPgRepo struct {
	db *db.DB
	qb *query.QueryBuilder
}

func NewTransactionPgRepo(
	db *db.DB,
	qb *query.QueryBuilder,
) *TransactionPgRepo {
	return &TransactionPgRepo{
		db: db,
		qb: qb,
	}
}

func (r *TransactionPgRepo) ListTransactions(
	ctx context.Context,
	userID uuid.UUID,
	opts ...repo.ListTransactionsOption,
) ([]entity.Transaction, error) {
	transactions, err := r.qb.ListTransactions(ctx, userID, opts...)
	if err != nil {
		return nil, errs.New(err)
	}
	return transactions, nil
}

func (r *TransactionPgRepo) ListTransactionsWithCategoriesAndInstitutions(
	ctx context.Context,
	userID uuid.UUID,
	opts ...repo.ListTransactionsOption,
) ([]entity.TransactionWithCategoryAndInstitution, error) {
	transactions, err := r.qb.
		ListTransactionsWithCategoriesAndInstitutions(ctx, userID, opts...)
	if err != nil {
		return nil, errs.New(err)
	}
	return transactions, nil
}

func (r *TransactionPgRepo) CountTransactions(
	ctx context.Context,
	userID uuid.UUID,
	opts ...repo.ListTransactionsOption,
) (int64, error) {
	return r.qb.CountTransactions(ctx, userID, opts...)
}

func (r *TransactionPgRepo) CreateTransactions(
	ctx context.Context,
	params []repo.CreateTransactionsParams,
) error {
	dbParams := make([]sqlc.CreateTransactionsParams, len(params))
	if err := copier.Copy(&dbParams, params); err != nil {
		return errs.New(err)
	}

	tx := r.db.UseTx(ctx)
	_, err := tx.CreateTransactions(ctx, dbParams)
	if err != nil {
		return errs.New(err)
	}

	return nil
}

var _ repo.TransactionRepo = &TransactionPgRepo{}
