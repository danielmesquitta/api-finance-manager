package dto

import (
	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
)

type ListTransactionsResponse struct {
	entity.PaginatedList[entity.Transaction]
}
