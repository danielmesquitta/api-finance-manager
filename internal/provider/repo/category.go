package repo

import (
	"context"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/google/uuid"
)

type CategoryOptions struct {
	Limit  uint   `json:"-"`
	Offset uint   `json:"-"`
	Search string `json:"search"`
}

type CategoryOption func(*CategoryOptions)

func WithCategoriesPagination(limit uint, offset uint) CategoryOption {
	return func(o *CategoryOptions) {
		o.Limit = limit
		o.Offset = offset
	}
}

func WithCategoriesSearch(search string) CategoryOption {
	return func(o *CategoryOptions) {
		o.Search = search
	}
}

type CategoryRepo interface {
	ListCategories(
		ctx context.Context,
		opts ...CategoryOption,
	) ([]entity.Category, error)
	CountCategories(
		ctx context.Context,
		opts ...CategoryOption,
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
	GetCategory(
		ctx context.Context,
		id uuid.UUID,
	) (*entity.Category, error)
}
