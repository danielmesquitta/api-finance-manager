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
	db *db.DB
}

func NewBudgetPgRepo(db *db.DB) *BudgetPgRepo {
	return &BudgetPgRepo{
		db: db,
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

	tx := r.db.UseTx(ctx)
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

func (r *BudgetPgRepo) DeleteBudgetCategories(
	ctx context.Context,
	userID uuid.UUID,
) error {
	tx := r.db.UseTx(ctx)
	return tx.DeleteBudgetCategories(ctx, userID)
}

func (r *BudgetPgRepo) DeleteBudgets(
	ctx context.Context,
	userID uuid.UUID,
) error {
	tx := r.db.UseTx(ctx)
	return tx.DeleteBudgets(ctx, userID)
}

func (r *BudgetPgRepo) GetBudget(
	ctx context.Context,
	params repo.GetBudgetParams,
) (*entity.Budget, error) {
	dbParams := sqlc.GetBudgetParams{}
	if err := copier.Copy(&dbParams, params); err != nil {
		return nil, errs.New(err)
	}

	budget, err := r.db.GetBudget(ctx, dbParams)
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

func (r *BudgetPgRepo) CreateBudgetCategories(
	ctx context.Context,
	params []repo.CreateBudgetCategoriesParams,
) error {
	dbParams := make([]sqlc.CreateBudgetCategoriesParams, len(params))
	if err := copier.Copy(&dbParams, params); err != nil {
		return errs.New(err)
	}

	tx := r.db.UseTx(ctx)
	_, err := tx.CreateBudgetCategories(ctx, dbParams)
	if err != nil {
		return errs.New(err)
	}

	return nil
}

func (r *BudgetPgRepo) GetBudgetCategories(
	ctx context.Context,
	budgetID uuid.UUID,
) ([]entity.BudgetCategory, []entity.Category, error) {
	rows, err := r.db.GetBudgetCategories(ctx, budgetID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, nil
		}
		return nil, nil, errs.New(err)
	}

	budgetCategories := []entity.BudgetCategory{}
	categories := []entity.Category{}
	for _, row := range rows {
		budgetCategory := entity.BudgetCategory{}
		if err := copier.Copy(&budgetCategory, row.BudgetCategory); err != nil {
			return nil, nil, errs.New(err)
		}
		budgetCategories = append(budgetCategories, budgetCategory)

		category := entity.Category{}
		if err := copier.Copy(&category, row.Category); err != nil {
			return nil, nil, errs.New(err)
		}
		categories = append(categories, category)
	}

	return budgetCategories, categories, nil
}

func (r *BudgetPgRepo) UpdateBudget(
	ctx context.Context,
	params repo.UpdateBudgetParams,
) error {
	dbParams := sqlc.UpdateBudgetParams{}
	if err := copier.Copy(&dbParams, params); err != nil {
		return errs.New(err)
	}

	tx := r.db.UseTx(ctx)
	return tx.UpdateBudget(ctx, dbParams)
}

var _ repo.BudgetRepo = &BudgetPgRepo{}
