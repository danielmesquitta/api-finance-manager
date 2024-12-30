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
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

type UserPgRepo struct {
	db *db.DB
}

func NewUserPgRepo(db *db.DB) *UserPgRepo {
	return &UserPgRepo{
		db: db,
	}
}

func (r *UserPgRepo) CreateUser(
	ctx context.Context,
	params repo.CreateUserParams,
) (*entity.User, error) {
	dbParams := sqlc.CreateUserParams{}
	if err := copier.Copy(&dbParams, params); err != nil {
		return nil, errs.New(err)
	}

	tx := r.db.UseTx(ctx)
	user, err := tx.CreateUser(ctx, dbParams)
	if err != nil {
		return nil, errs.New(err)
	}

	result := entity.User{}
	if err := copier.Copy(&result, user); err != nil {
		return nil, errs.New(err)
	}

	return &result, nil
}

func (r *UserPgRepo) GetUserByEmail(
	ctx context.Context,
	email string,
) (*entity.User, error) {
	user, err := r.db.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, errs.New(err)
	}

	result := entity.User{}
	if err := copier.Copy(&result, user); err != nil {
		return nil, errs.New(err)
	}

	return &result, nil
}

func (r *UserPgRepo) GetUserByID(
	ctx context.Context,
	id uuid.UUID,
) (*entity.User, error) {
	user, err := r.db.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, errs.New(err)
	}

	result := entity.User{}
	if err := copier.Copy(&result, user); err != nil {
		return nil, errs.New(err)
	}

	return &result, nil
}

func (r *UserPgRepo) UpdateUser(
	ctx context.Context,
	params repo.UpdateUserParams,
) (*entity.User, error) {
	dbParams := sqlc.UpdateUserParams{}
	if err := copier.Copy(&dbParams, params); err != nil {
		return nil, errs.New(err)
	}

	tx := r.db.UseTx(ctx)
	user, err := tx.UpdateUser(ctx, dbParams)
	if err != nil {
		return nil, errs.New(err)
	}

	result := entity.User{}
	if err := copier.Copy(&result, user); err != nil {
		return nil, errs.New(err)
	}

	return &result, nil
}

func (r *UserPgRepo) ListUsers(
	ctx context.Context,
) ([]entity.User, error) {
	users, err := r.db.ListUsers(ctx)
	if err != nil {
		return nil, errs.New(err)
	}

	results := []entity.User{}
	if err := copier.Copy(&results, users); err != nil {
		return nil, errs.New(err)
	}

	return results, nil
}

func (r *UserPgRepo) ListPremiumActiveUsersWithAccounts(
	ctx context.Context,
) ([]entity.User, []entity.Account, error) {
	rows, err := r.db.ListPremiumActiveUsersWithAccounts(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, nil
		}
		return nil, nil, errs.New(err)
	}

	uniqueUsers := map[uuid.UUID]entity.User{}
	uniqueAccounts := map[uuid.UUID]entity.Account{}
	for _, row := range rows {
		user := entity.User{}
		if err := copier.Copy(&user, row.User); err != nil {
			return nil, nil, errs.New(err)
		}
		uniqueUsers[user.ID] = user

		account := entity.Account{}
		if err := copier.Copy(&account, row.Account); err != nil {
			return nil, nil, errs.New(err)
		}
		uniqueAccounts[account.ID] = account
	}

	users := []entity.User{}
	for _, user := range uniqueUsers {
		users = append(users, user)
	}

	accounts := []entity.Account{}
	for _, account := range uniqueAccounts {
		accounts = append(accounts, account)
	}

	return users, accounts, nil
}

var _ repo.UserRepo = &UserPgRepo{}
