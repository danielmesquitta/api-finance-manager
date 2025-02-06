package dto

import (
	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
)

type ListCategoriesResponse struct {
	entity.PaginatedList[entity.TransactionCategory]
}
