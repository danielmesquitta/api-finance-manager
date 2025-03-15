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

type UserAuthProviderRepo struct {
	db *db.DB
}

func NewUserAuthProviderRepo(db *db.DB) *UserAuthProviderRepo {
	return &UserAuthProviderRepo{
		db: db,
	}
}

func (r *UserAuthProviderRepo) CreateUserAuthProvider(
	ctx context.Context,
	params repo.CreateUserAuthProviderParams,
) error {
	dbParams := sqlc.CreateUserAuthProviderParams{}
	if err := copier.Copy(&dbParams, params); err != nil {
		return errs.New(err)
	}

	tx := r.db.UseTx(ctx)
	if err := tx.CreateUserAuthProvider(ctx, dbParams); err != nil {
		return errs.New(err)
	}

	return nil
}

func (r *UserAuthProviderRepo) UpdateUserAuthProvider(
	ctx context.Context,
	params repo.UpdateUserAuthProviderParams,
) error {
	dbParams := sqlc.UpdateUserAuthProviderParams{}
	if err := copier.Copy(&dbParams, params); err != nil {
		return errs.New(err)
	}

	tx := r.db.UseTx(ctx)
	if err := tx.UpdateUserAuthProvider(ctx, dbParams); err != nil {
		return errs.New(err)
	}

	return nil
}

func (r *UserAuthProviderRepo) GetUserAuthProvider(
	ctx context.Context,
	params repo.GetUserAuthProviderParams,
) (*entity.UserAuthProvider, error) {
	dbParams := sqlc.GetUserAuthProviderParams{}
	if err := copier.Copy(&dbParams, params); err != nil {
		return nil, errs.New(err)
	}

	user, err := r.db.GetUserAuthProvider(ctx, dbParams)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, errs.New(err)
	}

	result := entity.UserAuthProvider{}
	if err := copier.Copy(&result, user); err != nil {
		return nil, errs.New(err)
	}

	return &result, nil
}

var _ repo.UserAuthProviderRepo = (*UserAuthProviderRepo)(nil)
