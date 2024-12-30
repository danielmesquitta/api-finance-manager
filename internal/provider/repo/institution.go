package repo

import (
	"context"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
)

type InstitutionRepo interface {
	ListInstitutions(
		ctx context.Context,
	) ([]entity.Institution, error)
	CreateInstitutions(
		ctx context.Context,
		params []CreateInstitutionsParams,
	) error
	GetInstitutionByExternalID(
		ctx context.Context,
		externalID string,
	) (*entity.Institution, error)
}
