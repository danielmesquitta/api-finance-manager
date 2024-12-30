package pgrepo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/db"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/db/sqlc"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
	"github.com/jinzhu/copier"
)

type InstitutionPgRepo struct {
	db *db.DB
}

func NewInstitutionPgRepo(db *db.DB) *InstitutionPgRepo {
	return &InstitutionPgRepo{
		db: db,
	}
}

func (r *InstitutionPgRepo) ListInstitutions(
	ctx context.Context,
) ([]entity.Institution, error) {
	institutions, err := r.db.ListInstitutions(ctx)
	if err != nil {
		return nil, errs.New(err)
	}

	results := []entity.Institution{}
	if err := copier.Copy(&results, institutions); err != nil {
		return nil, errs.New(err)
	}

	return results, nil
}

func (r *InstitutionPgRepo) CreateInstitutions(
	ctx context.Context,
	params []repo.CreateInstitutionsParams,
) error {
	dbParams := make([]sqlc.CreateInstitutionsParams, len(params))
	if err := copier.Copy(&dbParams, params); err != nil {
		return errs.New(err)
	}

	tx := r.db.UseTx(ctx)
	_, err := tx.CreateInstitutions(ctx, dbParams)
	if err != nil {
		return errs.New(err)
	}

	return nil
}

func (r *InstitutionPgRepo) GetInstitutionByExternalID(
	ctx context.Context,
	externalID string,
) (*entity.Institution, error) {
	institution, err := r.db.GetInstitutionByExternalID(ctx, externalID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, errs.New(err)
	}

	result := entity.Institution{}
	if err := copier.Copy(&result, institution); err != nil {
		return nil, errs.New(err)
	}

	return &result, nil
}

var _ repo.InstitutionRepo = &InstitutionPgRepo{}
