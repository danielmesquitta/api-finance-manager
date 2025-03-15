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

type UserInstitutionRepo struct {
	db *db.DB
}

func NewUserInstitutionRepo(db *db.DB) *UserInstitutionRepo {
	return &UserInstitutionRepo{
		db: db,
	}
}

func (r *UserInstitutionRepo) CreateUserInstitution(
	ctx context.Context,
	params repo.CreateUserInstitutionParams,
) (*entity.UserInstitution, error) {
	dbParams := sqlc.CreateUserInstitutionParams{}
	if err := copier.Copy(&dbParams, params); err != nil {
		return nil, errs.New(err)
	}

	tx := r.db.UseTx(ctx)
	user, err := tx.CreateUserInstitution(ctx, dbParams)
	if err != nil {
		return nil, errs.New(err)
	}

	result := entity.UserInstitution{}
	if err := copier.Copy(&result, user); err != nil {
		return nil, errs.New(err)
	}

	return &result, nil
}

func (r *UserInstitutionRepo) GetUserInstitutionByExternalID(
	ctx context.Context,
	externalID string,
) (*entity.UserInstitution, error) {
	user, err := r.db.GetUserInstitutionByExternalID(ctx, externalID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, errs.New(err)
	}

	result := entity.UserInstitution{}
	if err := copier.Copy(&result, user); err != nil {
		return nil, errs.New(err)
	}

	return &result, nil
}

var _ repo.UserInstitutionRepo = (*UserInstitutionRepo)(nil)
