package pgrepo

import (
	"context"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/db"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

type AccountPgRepo struct {
	q *db.Queries
}

func NewAccountPgRepo(q *db.Queries) *AccountPgRepo {
	return &AccountPgRepo{
		q: q,
	}
}

func (r *AccountPgRepo) ListAccountsByUserID(
	ctx context.Context,
	userID uuid.UUID,
) ([]entity.Account, error) {
	accounts, err := r.q.ListAccountsByUserID(ctx, userID)
	if err != nil {
		return nil, errs.New(err)
	}

	results := []entity.Account{}
	if err := copier.Copy(&results, accounts); err != nil {
		return nil, errs.New(err)
	}

	return results, nil
}

var _ repo.AccountRepo = &AccountPgRepo{}
