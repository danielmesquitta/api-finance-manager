package pgrepo

import (
	"context"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/db"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/db/sqlc"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
	"github.com/jinzhu/copier"
)

type CategoryPgRepo struct {
	q *db.Queries
}

func NewCategoryPgRepo(q *db.Queries) *CategoryPgRepo {
	return &CategoryPgRepo{
		q: q,
	}
}

func (r *CategoryPgRepo) ListCategories(
	ctx context.Context,
) ([]entity.Category, error) {
	categories, err := r.q.ListCategories(ctx)
	if err != nil {
		return nil, errs.New(err)
	}

	results := []entity.Category{}
	if err := copier.Copy(&results, categories); err != nil {
		return nil, errs.New(err)
	}

	return results, nil
}

func (r *CategoryPgRepo) CreateManyCategories(
	ctx context.Context,
	params []repo.CreateManyCategoriesParams,
) error {
	dbParams := make([]sqlc.CreateManyCategoriesParams, len(params))
	if err := copier.Copy(&dbParams, params); err != nil {
		return errs.New(err)
	}

	tx := r.q.UseTx(ctx)
	_, err := tx.CreateManyCategories(ctx, dbParams)
	if err != nil {
		return errs.New(err)
	}

	return nil
}

func (r *CategoryPgRepo) SearchCategories(
	ctx context.Context,
	params repo.SearchCategoriesParams,
) ([]entity.Category, error) {
	dbParams := sqlc.SearchCategoriesParams{}
	if err := copier.Copy(&dbParams, params); err != nil {
		return nil, errs.New(err)
	}

	categories, err := r.q.SearchCategories(ctx, dbParams)
	if err != nil {
		return nil, errs.New(err)
	}

	results := []entity.Category{}
	if err := copier.Copy(&results, categories); err != nil {
		return nil, errs.New(err)
	}

	return results, nil
}

func (r *CategoryPgRepo) CountSearchCategories(
	ctx context.Context,
	search string,
) (int64, error) {
	return r.q.CountSearchCategories(ctx, search)
}

var _ repo.CategoryRepo = &CategoryPgRepo{}
