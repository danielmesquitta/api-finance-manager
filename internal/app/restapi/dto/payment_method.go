package dto

import (
	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
)

type ListPaymentMethodsResponse struct {
	entity.PaginatedList[entity.PaymentMethod]
}
