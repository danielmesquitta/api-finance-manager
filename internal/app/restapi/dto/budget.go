package dto

import (
	"github.com/danielmesquitta/api-finance-manager/internal/domain/usecase"
	"github.com/google/uuid"
)

type UpsertBudgetRequest struct {
	usecase.UpsertBudgetInput
	UserID uuid.UUID `json:"-"`
}

type GetBudgetResponse struct {
	usecase.GetBudgetOutput
}
