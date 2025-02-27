package pgrepo

import (
	"context"
	"database/sql"
	"errors"

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
	opts ...repo.TransactionOption,
) ([]entity.Transaction, error) {
	transactions, err := r.qb.ListTransactions(ctx, userID, opts...)
	if err != nil {
		return nil, errs.New(err)
	}
	return transactions, nil
}

func (r *TransactionPgRepo) ListFullTransactions(
	ctx context.Context,
	userID uuid.UUID,
	opts ...repo.TransactionOption,
) ([]entity.FullTransaction, error) {
	transactions, err := r.qb.
		ListFullTransactions(ctx, userID, opts...)
	if err != nil {
		return nil, errs.New(err)
	}
	return transactions, nil
}

func (r *TransactionPgRepo) CountTransactions(
	ctx context.Context,
	userID uuid.UUID,
	opts ...repo.TransactionOption,
) (int64, error) {
	return r.qb.CountTransactions(ctx, userID, opts...)
}

func (r *TransactionPgRepo) SumTransactions(
	ctx context.Context,
	userID uuid.UUID,
	opts ...repo.TransactionOption,
) (int64, error) {
	return r.qb.SumTransactions(ctx, userID, opts...)
}

func (r *TransactionPgRepo) SumTransactionsByCategory(
	ctx context.Context,
	userID uuid.UUID,
	opts ...repo.TransactionOption,
) (map[uuid.UUID]int64, error) {
	return r.qb.SumTransactionsByCategory(ctx, userID, opts...)
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

func (r *TransactionPgRepo) CreateTransaction(
	ctx context.Context,
	params repo.CreateTransactionParams,
) error {
	var dbParams sqlc.CreateTransactionParams
	if err := copier.Copy(&dbParams, params); err != nil {
		return errs.New(err)
	}

	tx := r.db.UseTx(ctx)
	if err := tx.CreateTransaction(ctx, dbParams); err != nil {
		return errs.New(err)
	}

	return nil
}

func (r *TransactionPgRepo) GetTransaction(
	ctx context.Context,
	params repo.GetTransactionParams,
) (*entity.FullTransaction, error) {
	dbParams := sqlc.GetTransactionParams{}
	if err := copier.Copy(&dbParams, params); err != nil {
		return nil, errs.New(err)
	}

	transaction, err := r.db.GetTransaction(ctx, dbParams)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, errs.New(err)
	}

	result := entity.FullTransaction{}
	if err := copier.Copy(&result, transaction); err != nil {
		return nil, errs.New(err)
	}

	return &result, nil
}

func (r *TransactionPgRepo) UpdateTransaction(
	ctx context.Context,
	params repo.UpdateTransactionParams,
) error {
	dbParams := sqlc.UpdateTransactionParams{}
	if err := copier.Copy(&dbParams, params); err != nil {
		return errs.New(err)
	}

	tx := r.db.UseTx(ctx)
	return tx.UpdateTransaction(ctx, dbParams)
}

var _ repo.TransactionRepo = (*TransactionPgRepo)(nil)
