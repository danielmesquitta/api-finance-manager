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

type TransactionCategoryRepo interface {
	ListTransactionCategories(
		ctx context.Context,
		opts ...TransactionCategoryOptions,
	) ([]entity.TransactionCategory, error)
	CountTransactionCategories(
		ctx context.Context,
		opts ...TransactionCategoryOptions,
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
