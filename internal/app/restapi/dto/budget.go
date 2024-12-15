package dto

import (
	"github.com/danielmesquitta/api-finance-manager/internal/domain/usecase"
)

type UpsertBudgetRequest struct {
	usecase.UpsertBudgetUseCaseInput
}
