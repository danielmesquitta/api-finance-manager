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

type TransactionCategoryPgRepo struct {
	db *db.DB
	qb *query.QueryBuilder
}

func NewCategoryPgRepo(
	db *db.DB,
	qb *query.QueryBuilder,
) *TransactionCategoryPgRepo {
	return &TransactionCategoryPgRepo{
		db: db,
		qb: qb,
	}
}

func (r *TransactionCategoryPgRepo) ListTransactionCategories(
	ctx context.Context,
	opts ...repo.TransactionCategoryOption,
) ([]entity.TransactionCategory, error) {
	categories, err := r.qb.ListTransactionCategories(ctx, opts...)
	if err != nil {
		return nil, errs.New(err)
	}

	return categories, nil
}

func (r *TransactionCategoryPgRepo) CountTransactionCategories(
	ctx context.Context,
	opts ...repo.TransactionCategoryOption,
) (int64, error) {
	return r.qb.CountTransactionCategories(ctx, opts...)
}

func (r *TransactionCategoryPgRepo) CountTransactionCategoriesByIDs(
	ctx context.Context,
	ids []uuid.UUID,
) (int64, error) {
	return r.db.CountTransactionCategoriesByIDs(ctx, ids)
}

func (r *TransactionCategoryPgRepo) CreateTransactionCategories(
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

func (r *TransactionCategoryPgRepo) ListTransactionCategoriesByExternalIDs(
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

func (r *TransactionCategoryPgRepo) GetTransactionCategory(
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

var _ repo.TransactionCategoryRepo = &TransactionCategoryPgRepo{}
