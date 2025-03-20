package dto

import (
	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
)

type ListInstitutionsResponse struct {
	entity.PaginatedList[entity.Institution]
}
