package repo

import (
	"context"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
)

type CategoryRepo interface {
	ListCategories(
		ctx context.Context,
	) ([]entity.Category, error)
	CreateCategories(
		ctx context.Context,
		params []CreateCategoriesParams,
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
