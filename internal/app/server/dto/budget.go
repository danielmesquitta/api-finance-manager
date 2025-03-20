package dto

import (
	"github.com/danielmesquitta/api-finance-manager/internal/domain/usecase/budget"
)

type UpsertBudgetRequest struct {
	budget.UpsertBudgetUseCaseInput
}

type GetBudgetResponse struct {
	budget.GetBudgetUseCaseOutput
}

type GetBudgetCategoryResponse struct {
	budget.GetBudgetCategoryUseCaseOutput
}
