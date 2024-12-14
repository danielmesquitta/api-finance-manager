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
}
