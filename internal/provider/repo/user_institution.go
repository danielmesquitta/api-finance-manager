package repo

import (
	"context"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
)

type UserInstitutionRepo interface {
	CreateUserInstitution(
		ctx context.Context,
		params CreateUserInstitutionParams,
	) (*entity.UserInstitution, error)
	GetUserInstitutionByExternalID(
		ctx context.Context,
		externalID string,
	) (*entity.UserInstitution, error)
}
