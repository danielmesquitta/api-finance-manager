package repo

import (
	"context"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/google/uuid"
)

type ListInstitutionsOptions struct {
	Limit  uint      `json:"-"`
	Offset uint      `json:"-"`
	Search string    `json:"search"`
	UserID uuid.UUID `json:"-"`
}

type ListInstitutionsOption func(*ListInstitutionsOptions)

func WithInstitutionsPagination(
	limit uint,
	offset uint,
) ListInstitutionsOption {
	return func(o *ListInstitutionsOptions) {
		o.Limit = limit
		o.Offset = offset
	}
}

func WithInstitutionsSearch(search string) ListInstitutionsOption {
	return func(o *ListInstitutionsOptions) {
		o.Search = search
	}
}

func WithUser(userID uuid.UUID) ListInstitutionsOption {
	return func(o *ListInstitutionsOptions) {
		o.UserID = userID
	}
}

type InstitutionRepo interface {
	ListInstitutions(
		ctx context.Context,
		opts ...ListInstitutionsOption,
	) ([]entity.Institution, error)
	CountInstitutions(
		ctx context.Context,
		opts ...ListInstitutionsOption,
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
