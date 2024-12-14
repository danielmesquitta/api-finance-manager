package repo

import (
	"context"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
)

type CategoryRepo interface {
	ListCategories(
		ctx context.Context,
	) ([]entity.Category, error)
	CreateManyCategories(
		ctx context.Context,
		params []CreateManyCategoriesParams,
	) error
	SearchCategories(
		ctx context.Context,
		arg SearchCategoriesParams,
	) ([]entity.Category, error)
	CountSearchCategories(
		ctx context.Context,
		search string,
	) (int64, error)
}
