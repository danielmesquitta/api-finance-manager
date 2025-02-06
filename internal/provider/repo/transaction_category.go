package repo

import (
	"context"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/google/uuid"
)

type TransactionCategoryOptions struct {
	Limit  uint   `json:"-"`
	Offset uint   `json:"-"`
	Search string `json:"search"`
}

type TransactionCategoryOption func(*TransactionCategoryOptions)

func WithTransactionCategoryPagination(
	limit uint,
	offset uint,
) TransactionCategoryOption {
	return func(o *TransactionCategoryOptions) {
		o.Limit = limit
		o.Offset = offset
	}
}

func WithTransactionCategorySearch(search string) TransactionCategoryOption {
	return func(o *TransactionCategoryOptions) {
		o.Search = search
	}
}

type CategoryRepo interface {
	ListTransactionCategories(
		ctx context.Context,
		opts ...TransactionCategoryOption,
	) ([]entity.TransactionCategory, error)
	CountTransactionCategories(
		ctx context.Context,
		opts ...TransactionCategoryOption,
	) (int64, error)
	CountTransactionCategoriesByIDs(
		ctx context.Context,
		ids []uuid.UUID,
	) (int64, error)
	CreateTransactionCategories(
		ctx context.Context,
		params []CreateTransactionCategoriesParams,
	) error
	ListTransactionCategoriesByExternalIDs(
		ctx context.Context,
		externalIDs []string,
	) ([]entity.TransactionCategory, error)
	GetTransactionCategory(
		ctx context.Context,
		id uuid.UUID,
	) (*entity.TransactionCategory, error)
}
