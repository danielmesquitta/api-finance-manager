package pgrepo

import (
	"context"
	"database/sql"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/db"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/db/sqlc"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

type UserPgRepo struct {
	q *db.Queries
}

func NewUserPgRepo(q *db.Queries) *UserPgRepo {
	return &UserPgRepo{
		q: q,
	}
}

func (r *UserPgRepo) CreateUser(
	ctx context.Context,
	params repo.CreateUserParams,
) (*entity.User, error) {
	sqlcParams := sqlc.CreateUserParams{}
	if err := copier.Copy(&sqlcParams, params); err != nil {
		return nil, err
	}

	tx := r.q.UseTx(ctx)
	user, err := tx.CreateUser(ctx, sqlcParams)
	if err != nil {
		return nil, err
	}

	result := entity.User{}
	if err := copier.Copy(&result, user); err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *UserPgRepo) GetUserByEmail(
	ctx context.Context,
	email string,
) (*entity.User, error) {
	user, err := r.q.GetUserByEmail(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	result := entity.User{}
	if err := copier.Copy(&result, user); err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *UserPgRepo) GetUserByID(
	ctx context.Context,
	id uuid.UUID,
) (*entity.User, error) {
	user, err := r.q.GetUserByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	result := entity.User{}
	if err := copier.Copy(&result, user); err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *UserPgRepo) UpdateUser(
	ctx context.Context,
	params repo.UpdateUserParams,
) (*entity.User, error) {
	sqlcParams := sqlc.UpdateUserParams{}
	if err := copier.Copy(&sqlcParams, params); err != nil {
		return nil, err
	}

	tx := r.q.UseTx(ctx)
	user, err := tx.UpdateUser(ctx, sqlcParams)
	if err != nil {
		return nil, err
	}

	result := entity.User{}
	if err := copier.Copy(&result, user); err != nil {
		return nil, err
	}

	return &result, nil
}

var _ repo.UserRepo = &UserPgRepo{}
