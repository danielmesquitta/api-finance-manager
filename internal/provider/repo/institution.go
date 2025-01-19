package repo

import (
	"context"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/google/uuid"
)

type InstitutionOptions struct {
	Limit  uint      `json:"-"`
	Offset uint      `json:"-"`
	Search string    `json:"search"`
	UserID uuid.UUID `json:"-"`
}

type InstitutionOption func(*InstitutionOptions)

func WithInstitutionsPagination(
	limit uint,
	offset uint,
) InstitutionOption {
	return func(o *InstitutionOptions) {
		o.Limit = limit
		o.Offset = offset
	}
}

func WithInstitutionsSearch(search string) InstitutionOption {
	return func(o *InstitutionOptions) {
		o.Search = search
	}
}

func WithUser(userID uuid.UUID) InstitutionOption {
	return func(o *InstitutionOptions) {
		o.UserID = userID
	}
}

type InstitutionRepo interface {
	ListInstitutions(
		ctx context.Context,
		opts ...InstitutionOption,
	) ([]entity.Institution, error)
	CountInstitutions(
		ctx context.Context,
		opts ...InstitutionOption,
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
