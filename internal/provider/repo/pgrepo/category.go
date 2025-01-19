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

type CategoryPgRepo struct {
	db *db.DB
	qb *query.QueryBuilder
}

func NewCategoryPgRepo(
	db *db.DB,
	qb *query.QueryBuilder,
) *CategoryPgRepo {
	return &CategoryPgRepo{
		db: db,
		qb: qb,
	}
}

func (r *CategoryPgRepo) ListCategories(
	ctx context.Context,
	opts ...repo.ListCategoriesOption,
) ([]entity.Category, error) {
	categories, err := r.qb.ListCategories(ctx, opts...)
	if err != nil {
		return nil, errs.New(err)
	}

	return categories, nil
}

func (r *CategoryPgRepo) CountCategories(
	ctx context.Context,
	opts ...repo.ListCategoriesOption,
) (int64, error) {
	return r.qb.CountCategories(ctx, opts...)
}

func (r *CategoryPgRepo) CountCategoriesByIDs(
	ctx context.Context,
	ids []uuid.UUID,
) (int64, error) {
	return r.db.CountCategoriesByIDs(ctx, ids)
}

func (r *CategoryPgRepo) CreateCategories(
	ctx context.Context,
	params []repo.CreateCategoriesParams,
) error {
	dbParams := make([]sqlc.CreateCategoriesParams, len(params))
	if err := copier.Copy(&dbParams, params); err != nil {
		return errs.New(err)
	}

	tx := r.db.UseTx(ctx)
	_, err := tx.CreateCategories(ctx, dbParams)
	if err != nil {
		return errs.New(err)
	}

	return nil
}

func (r *CategoryPgRepo) ListCategoriesByExternalIDs(
	ctx context.Context,
	externalIDs []string,
) ([]entity.Category, error) {
	categories, err := r.db.ListCategoriesByExternalIDs(ctx, externalIDs)
	if err != nil {
		return nil, errs.New(err)
	}

	results := []entity.Category{}
	if err := copier.Copy(&results, categories); err != nil {
		return nil, errs.New(err)
	}

	return results, nil
}

var _ repo.CategoryRepo = &CategoryPgRepo{}
