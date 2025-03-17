package repo

import (
	"context"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/google/uuid"
)

type TransactionCategoryOptions struct {
	Limit  uint        `json:"-"`
	Offset uint        `json:"-"`
	Search string      `json:"search"`
	IDs    []uuid.UUID `json:"category_ids"`
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

func WithTransactionCategoryIDs(
	ids []uuid.UUID,
) TransactionCategoryOption {
	return func(o *TransactionCategoryOptions) {
		o.IDs = ids
	}
}

type TransactionCategoryRepo interface {
	ListTransactionCategories(
		ctx context.Context,
		opts ...TransactionCategoryOption,
	) ([]entity.TransactionCategory, error)
	CountTransactionCategories(
		ctx context.Context,
		opts ...TransactionCategoryOption,
	) (int64, error)
	CreateTransactionCategories(
		ctx context.Context,
		params []CreateTransactionCategoriesParams,
	) error
	GetTransactionCategoryByID(
		ctx context.Context,
		id uuid.UUID,
	) (*entity.TransactionCategory, error)
	GetDefaultTransactionCategory(
		ctx context.Context,
	) (*entity.TransactionCategory, error)
}
