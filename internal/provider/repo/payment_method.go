package repo

import (
	"context"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/google/uuid"
)

type PaymentMethodOptions struct {
	Limit  uint   `json:"-"`
	Offset uint   `json:"-"`
	Search string `json:"search"`
}

type PaymentMethodRepo interface {
	ListPaymentMethods(
		ctx context.Context,
		opts ...PaymentMethodOptions,
	) ([]entity.PaymentMethod, error)
	CountPaymentMethods(
		ctx context.Context,
		opts ...PaymentMethodOptions,
	) (int64, error)
	CreatePaymentMethods(
		ctx context.Context,
		params []CreatePaymentMethodsParams,
	) error
	GetPaymentMethodByID(
		ctx context.Context,
		id uuid.UUID,
	) (*entity.PaymentMethod, error)
}
