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
	db *db.DB
}

func NewAccountPgRepo(db *db.DB) *AccountPgRepo {
	return &AccountPgRepo{
		db: db,
	}
}

func (r *AccountPgRepo) ListAccountsByUserID(
	ctx context.Context,
	userID uuid.UUID,
) ([]entity.Account, error) {
	accounts, err := r.db.ListAccountsByUserID(ctx, userID)
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
