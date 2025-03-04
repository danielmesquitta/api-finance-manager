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

type TransactionCategoryRepo struct {
	db *db.DB
	qb *query.QueryBuilder
}

func NewCategoryRepo(
	db *db.DB,
	qb *query.QueryBuilder,
) *TransactionCategoryRepo {
	return &TransactionCategoryRepo{
		db: db,
		qb: qb,
	}
}

func (r *TransactionCategoryRepo) ListTransactionCategories(
	ctx context.Context,
	opts ...repo.TransactionCategoryOption,
) ([]entity.TransactionCategory, error) {
	categories, err := r.qb.ListTransactionCategories(ctx, opts...)
	if err != nil {
		return nil, errs.New(err)
	}

	return categories, nil
}

func (r *TransactionCategoryRepo) CountTransactionCategories(
	ctx context.Context,
	opts ...repo.TransactionCategoryOption,
) (int64, error) {
	return r.qb.CountTransactionCategories(ctx, opts...)
}

func (r *TransactionCategoryRepo) CountTransactionCategoriesByIDs(
	ctx context.Context,
	ids []uuid.UUID,
) (int64, error) {
	return r.db.CountTransactionCategoriesByIDs(ctx, ids)
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

func (r *TransactionCategoryRepo) ListTransactionCategoriesByExternalIDs(
	ctx context.Context,
	externalIDs []string,
) ([]entity.TransactionCategory, error) {
	categories, err := r.db.ListTransactionCategoriesByExternalIDs(
		ctx,
		externalIDs,
	)
	if err != nil {
		return nil, errs.New(err)
	}

	results := []entity.TransactionCategory{}
	if err := copier.Copy(&results, categories); err != nil {
		return nil, errs.New(err)
	}

	return results, nil
}

func (r *TransactionCategoryRepo) GetTransactionCategory(
	ctx context.Context,
	id uuid.UUID,
) (*entity.TransactionCategory, error) {
	category, err := r.db.GetTransactionCategory(ctx, id)
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
