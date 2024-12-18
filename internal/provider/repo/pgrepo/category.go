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
) ([]entity.Category, error) {
	categories, err := r.db.ListCategories(ctx)
	if err != nil {
		return nil, errs.New(err)
	}

	results := []entity.Category{}
	if err := copier.Copy(&results, categories); err != nil {
		return nil, errs.New(err)
	}

	return results, nil
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

func (r *CategoryPgRepo) SearchCategories(
	ctx context.Context,
	params repo.SearchCategoriesParams,
) ([]entity.Category, error) {
	return r.qb.SearchCategories(ctx, params)
}

func (r *CategoryPgRepo) CountSearchCategories(
	ctx context.Context,
	search string,
) (int64, error) {
	return r.qb.CountSearchCategories(ctx, search)
}

var _ repo.CategoryRepo = &CategoryPgRepo{}
