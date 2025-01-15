package repo

import (
	"context"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
)

type ListPaymentMethodsOptions struct {
	Limit  uint   `json:"-"`
	Offset uint   `json:"-"`
	Search string `json:"search"`
}

type ListPaymentMethodsOption func(*ListPaymentMethodsOptions)

func WithPaymentMethodsPagination(
	limit uint,
	offset uint,
) ListPaymentMethodsOption {
	return func(o *ListPaymentMethodsOptions) {
		o.Limit = limit
		o.Offset = offset
	}
}

func WithPaymentMethodsSearch(search string) ListPaymentMethodsOption {
	return func(o *ListPaymentMethodsOptions) {
		o.Search = search
	}
}

type PaymentMethodRepo interface {
	ListPaymentMethods(
		ctx context.Context,
		opts ...ListPaymentMethodsOption,
	) ([]entity.PaymentMethod, error)
	CountPaymentMethods(
		ctx context.Context,
		opts ...ListPaymentMethodsOption,
	) (int64, error)
	CreatePaymentMethods(
		ctx context.Context,
		params []CreatePaymentMethodsParams,
	) error
}
