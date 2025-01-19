package repo

import (
	"context"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/google/uuid"
)

type ListCategoriesOptions struct {
	Limit  uint   `json:"-"`
	Offset uint   `json:"-"`
	Search string `json:"search"`
}

type ListCategoriesOption func(*ListCategoriesOptions)

func WithCategoriesPagination(limit uint, offset uint) ListCategoriesOption {
	return func(o *ListCategoriesOptions) {
		o.Limit = limit
		o.Offset = offset
	}
}

func WithCategoriesSearch(search string) ListCategoriesOption {
	return func(o *ListCategoriesOptions) {
		o.Search = search
	}
}

type CategoryRepo interface {
	ListCategories(
		ctx context.Context,
		opts ...ListCategoriesOption,
	) ([]entity.Category, error)
	CountCategories(
		ctx context.Context,
		opts ...ListCategoriesOption,
	) (int64, error)
	CountCategoriesByIDs(ctx context.Context, ids []uuid.UUID) (int64, error)
	CreateCategories(
		ctx context.Context,
		params []CreateCategoriesParams,
	) error
	ListCategoriesByExternalIDs(
		ctx context.Context,
		externalIDs []string,
	) ([]entity.Category, error)
}
