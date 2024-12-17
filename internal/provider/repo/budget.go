package repo

import (
	"context"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/google/uuid"
)

type BudgetRepo interface {
	CreateBudget(
		ctx context.Context,
		params CreateBudgetParams,
	) (*entity.Budget, error)
	CreateBudgetCategories(
		ctx context.Context,
		params []CreateBudgetCategoriesParams,
	) error
	DeleteBudgetCategories(
		ctx context.Context,
		userID uuid.UUID,
	) error
	DeleteBudgets(ctx context.Context, userID uuid.UUID) error
	GetBudget(
		ctx context.Context,
		params GetBudgetParams,
	) (*entity.Budget, error)
	GetBudgetCategories(
		ctx context.Context,
		budgetID uuid.UUID,
	) ([]entity.BudgetCategory, []entity.Category, error)
	UpdateBudget(ctx context.Context, params UpdateBudgetParams) error
}
