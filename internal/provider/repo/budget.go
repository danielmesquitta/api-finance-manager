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
	DeleteBudgetCategoriesByBudgetID(
		ctx context.Context,
		budgetID uuid.UUID,
	) error
	GetBudgetByUserID(
		ctx context.Context,
		userID uuid.UUID,
	) (*entity.Budget, error)
	UpdateBudget(ctx context.Context, arg UpdateBudgetParams) error
}
