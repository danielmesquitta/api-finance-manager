package repo

import (
	"context"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
)

type PaymentMethodOptions struct {
	Limit  uint   `json:"-"`
	Offset uint   `json:"-"`
	Search string `json:"search"`
}

type PaymentMethodOption func(*PaymentMethodOptions)

func WithPaymentMethodPagination(
	limit uint,
	offset uint,
) PaymentMethodOption {
	return func(o *PaymentMethodOptions) {
		o.Limit = limit
		o.Offset = offset
	}
}

func WithPaymentMethodSearch(search string) PaymentMethodOption {
	return func(o *PaymentMethodOptions) {
		o.Search = search
	}
}

type PaymentMethodRepo interface {
	ListPaymentMethods(
		ctx context.Context,
		opts ...PaymentMethodOption,
	) ([]entity.PaymentMethod, error)
	CountPaymentMethods(
		ctx context.Context,
		opts ...PaymentMethodOption,
	) (int64, error)
	CreatePaymentMethods(
		ctx context.Context,
		params []CreatePaymentMethodsParams,
	) error
}
