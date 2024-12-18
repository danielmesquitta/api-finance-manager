package openfinance

import (
	"context"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
)

type Client interface {
	ListInstitutions(
		ctx context.Context,
	) ([]entity.Institution, error)
	ListCategories(
		ctx context.Context,
	) ([]entity.Category, error)
}
