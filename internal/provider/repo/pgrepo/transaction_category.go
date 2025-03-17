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

type TransactionCategoryRepo struct {
	db *db.DB
}

func NewCategoryRepo(
	db *db.DB,
) *TransactionCategoryRepo {
	return &TransactionCategoryRepo{
		db: db,
	}
}

func (r *TransactionCategoryRepo) ListTransactionCategories(
	ctx context.Context,
	opts ...repo.TransactionCategoryOption,
) ([]entity.TransactionCategory, error) {
	categories, err := r.db.ListTransactionCategories(ctx, opts...)
	if err != nil {
		return nil, errs.New(err)
	}

	return categories, nil
}

func (r *TransactionCategoryRepo) CountTransactionCategories(
	ctx context.Context,
	opts ...repo.TransactionCategoryOption,
) (int64, error) {
	return r.db.CountTransactionCategories(ctx, opts...)
}

func (r *TransactionCategoryRepo) CreateTransactionCategories(
	ctx context.Context,
	params []repo.CreateTransactionCategoriesParams,
) error {
	dbParams := make([]sqlc.CreateTransactionCategoriesParams, len(params))
	if err := copier.Copy(&dbParams, params); err != nil {
		return errs.New(err)
	}

	tx := r.db.UseTx(ctx)
	_, err := tx.CreateTransactionCategories(ctx, dbParams)
	if err != nil {
		return errs.New(err)
	}

	return nil
}

func (r *TransactionCategoryRepo) GetTransactionCategoryByID(
	ctx context.Context,
	id uuid.UUID,
) (*entity.TransactionCategory, error) {
	category, err := r.db.GetTransactionCategoryByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, errs.New(err)
	}

	result := entity.TransactionCategory{}
	if err := copier.Copy(&result, category); err != nil {
		return nil, errs.New(err)
	}

	return &result, nil
}

func (r *TransactionCategoryRepo) GetDefaultTransactionCategory(
	ctx context.Context,
) (*entity.TransactionCategory, error) {
	category, err := r.db.GetDefaultTransactionCategory(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, errs.New(err)
	}

	result := entity.TransactionCategory{}
	if err := copier.Copy(&result, category); err != nil {
		return nil, errs.New(err)
	}

	return &result, nil
}

var _ repo.TransactionCategoryRepo = (*TransactionCategoryRepo)(nil)
