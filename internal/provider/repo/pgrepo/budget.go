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

type BudgetPgRepo struct {
	q *db.Queries
}

func NewBudgetPgRepo(q *db.Queries) *BudgetPgRepo {
	return &BudgetPgRepo{
		q: q,
	}
}

func (r *BudgetPgRepo) CreateBudget(
	ctx context.Context,
	params repo.CreateBudgetParams,
) (*entity.Budget, error) {
	dbParams := sqlc.CreateBudgetParams{}
	if err := copier.Copy(&dbParams, params); err != nil {
		return nil, errs.New(err)
	}

	tx := r.q.UseTx(ctx)
	budget, err := tx.CreateBudget(ctx, dbParams)
	if err != nil {
		return nil, errs.New(err)
	}

	result := entity.Budget{}
	if err := copier.Copy(&result, budget); err != nil {
		return nil, errs.New(err)
	}

	return &result, nil
}

func (r *BudgetPgRepo) CreateBudgetCategories(
	ctx context.Context,
	params []repo.CreateBudgetCategoriesParams,
) error {
	dbParams := make([]sqlc.CreateBudgetCategoriesParams, len(params))
	if err := copier.Copy(&dbParams, params); err != nil {
		return errs.New(err)
	}

	tx := r.q.UseTx(ctx)
	_, err := tx.CreateBudgetCategories(ctx, dbParams)
	if err != nil {
		return errs.New(err)
	}

	return nil
}

func (r *BudgetPgRepo) DeleteBudgetByID(
	ctx context.Context,
	id uuid.UUID,
) error {
	tx := r.q.UseTx(ctx)
	return tx.DeleteBudgetByID(ctx, id)
}

func (r *BudgetPgRepo) DeleteBudgetCategoriesByBudgetID(
	ctx context.Context,
	budgetID uuid.UUID,
) error {
	tx := r.q.UseTx(ctx)
	return tx.DeleteBudgetCategoriesByBudgetID(ctx, budgetID)
}

func (r *BudgetPgRepo) GetBudgetByUserID(
	ctx context.Context,
	userID uuid.UUID,
) (*entity.Budget, error) {
	budget, err := r.q.GetBudgetByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, errs.New(err)
	}

	result := entity.Budget{}
	if err := copier.Copy(&result, budget); err != nil {
		return nil, errs.New(err)
	}

	return &result, nil
}

func (r *BudgetPgRepo) GetBudgetWithCategoriesByUserID(
	ctx context.Context,
	userID uuid.UUID,
) (*entity.Budget, []entity.BudgetCategory, []entity.Category, error) {
	rows, err := r.q.GetBudgetWithCategoriesByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, nil, nil
		}
		return nil, nil, nil, errs.New(err)
	}

	budget := entity.Budget{}
	if err := copier.Copy(&budget, rows[0].Budget); err != nil {
		return nil, nil, nil, errs.New(err)
	}

	budgetCategories := []entity.BudgetCategory{}
	categories := []entity.Category{}

	for _, row := range rows {
		budgetCategory := entity.BudgetCategory{}
		if err := copier.Copy(&budgetCategory, row.BudgetCategory); err != nil {
			return nil, nil, nil, errs.New(err)
		}
		budgetCategories = append(budgetCategories, budgetCategory)

		category := entity.Category{}
		if err := copier.Copy(&category, row.Category); err != nil {
			return nil, nil, nil, errs.New(err)
		}
		categories = append(categories, category)
	}

	return &budget, budgetCategories, categories, nil
}

func (r *BudgetPgRepo) UpdateBudget(
	ctx context.Context,
	params repo.UpdateBudgetParams,
) error {
	dbParams := sqlc.UpdateBudgetParams{}
	if err := copier.Copy(&dbParams, params); err != nil {
		return errs.New(err)
	}

	tx := r.q.UseTx(ctx)
	return tx.UpdateBudget(ctx, dbParams)
}

var _ repo.BudgetRepo = &BudgetPgRepo{}
