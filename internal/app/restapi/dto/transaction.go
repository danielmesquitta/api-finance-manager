package dto

import (
	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/usecase"
)

type ListTransactionsResponse struct {
	entity.PaginatedList[entity.FullTransaction]
}

type GetTransactionResponse struct {
	entity.FullTransaction
}

type UpdateTransactionRequest struct {
	usecase.UpdateTransactionInput
}
