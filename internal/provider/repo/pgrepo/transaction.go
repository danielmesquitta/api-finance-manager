package pgrepo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/db"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/db/sqlc"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

type TransactionRepo struct {
	db *db.DB
}

func NewTransactionRepo(
	db *db.DB,
) *TransactionRepo {
	return &TransactionRepo{
		db: db,
	}
}

func (r *TransactionRepo) ListTransactions(
	ctx context.Context,
	userID uuid.UUID,
	opts ...repo.TransactionOption,
) ([]entity.Transaction, error) {
	transactions, err := r.db.ListTransactions(ctx, userID, opts...)
	if err != nil {
		return nil, errs.New(err)
	}
	return transactions, nil
}

func (r *TransactionRepo) ListFullTransactions(
	ctx context.Context,
	userID uuid.UUID,
	opts ...repo.TransactionOption,
) ([]entity.FullTransaction, error) {
	transactions, err := r.db.
		ListFullTransactions(ctx, userID, opts...)
	if err != nil {
		return nil, errs.New(err)
	}
	return transactions, nil
}

func (r *TransactionRepo) CountTransactions(
	ctx context.Context,
	userID uuid.UUID,
	opts ...repo.TransactionOption,
) (int64, error) {
	return r.db.CountTransactions(ctx, userID, opts...)
}

func (r *TransactionRepo) SumTransactions(
	ctx context.Context,
	userID uuid.UUID,
	opts ...repo.TransactionOption,
) (int64, error) {
	return r.db.SumTransactions(ctx, userID, opts...)
}

func (r *TransactionRepo) SumTransactionsByCategory(
	ctx context.Context,
	userID uuid.UUID,
	opts ...repo.TransactionOption,
) (map[uuid.UUID]int64, error) {
	return r.db.SumTransactionsByCategory(ctx, userID, opts...)
}

func (r *TransactionRepo) CreateTransactions(
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

func (r *TransactionRepo) CreateTransaction(
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

func (r *TransactionRepo) GetTransaction(
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

func (r *TransactionRepo) UpdateTransaction(
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

var _ repo.TransactionRepo = (*TransactionRepo)(nil)
