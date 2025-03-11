package dto

import (
	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
)

type ListTransactionCategoriesResponse struct {
	entity.PaginatedList[entity.TransactionCategory]
}
