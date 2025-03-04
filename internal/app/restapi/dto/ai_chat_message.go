package dto

import (
	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
)

type ListAIChatMessagesResponse struct {
	entity.PaginatedList[entity.AIChatMessage]
}
