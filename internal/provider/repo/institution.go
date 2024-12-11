package repo

import (
	"context"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
)

type InstitutionRepo interface {
	ListInstitutions(
		ctx context.Context,
	) ([]entity.Institution, error)
	CreateManyInstitutions(
		ctx context.Context,
		params []CreateManyInstitutionsParams,
	) error
}
