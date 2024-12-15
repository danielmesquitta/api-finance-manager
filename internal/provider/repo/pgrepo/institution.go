package pgrepo

import (
	"context"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/db"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/db/sqlc"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
	"github.com/jinzhu/copier"
)

type InstitutionPgRepo struct {
	q *db.Queries
}

func NewInstitutionPgRepo(q *db.Queries) *InstitutionPgRepo {
	return &InstitutionPgRepo{
		q: q,
	}
}

func (r *InstitutionPgRepo) ListInstitutions(
	ctx context.Context,
) ([]entity.Institution, error) {
	institutions, err := r.q.ListInstitutions(ctx)
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

	tx := r.q.UseTx(ctx)
	_, err := tx.CreateInstitutions(ctx, dbParams)
	if err != nil {
		return errs.New(err)
	}

	return nil
}

var _ repo.InstitutionRepo = &InstitutionPgRepo{}
