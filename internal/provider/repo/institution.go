package repo

import (
	"context"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/google/uuid"
)

type InstitutionOptions struct {
	Limit   uint        `json:"-"`
	Offset  uint        `json:"-"`
	Search  string      `json:"search"`
	UserIDs []uuid.UUID `json:"-"`
}

type InstitutionRepo interface {
	ListInstitutions(
		ctx context.Context,
		opts ...InstitutionOptions,
	) ([]entity.Institution, error)
	CountInstitutions(
		ctx context.Context,
		opts ...InstitutionOptions,
	) (int64, error)
	CreateInstitutions(
		ctx context.Context,
		params []CreateInstitutionsParams,
	) error
	GetInstitutionByExternalID(
		ctx context.Context,
		externalID string,
	) (*entity.Institution, error)
}
