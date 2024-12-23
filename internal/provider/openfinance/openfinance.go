package openfinance

import (
	"context"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
)

type ListInstitutionsParams struct {
	Types  []string `json:"types,omitempty"`
	Search string   `json:"search,omitempty"`
}

type Client interface {
	ListInstitutions(
		ctx context.Context,
		params ListInstitutionsParams,
	) ([]entity.Institution, error)
	ListCategories(
		ctx context.Context,
	) ([]entity.Category, error)
}
