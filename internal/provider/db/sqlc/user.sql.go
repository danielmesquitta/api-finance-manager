// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: user.sql

package sqlc

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (
    external_id,
    provider,
    name,
    email,
    verified_email,
    tier,
    avatar,
    subscription_expires_at
  )
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id, external_id, provider, name, email, verified_email, tier, avatar, subscription_expires_at, synchronized_at, created_at, updated_at, deleted_at
`

type CreateUserParams struct {
	ExternalID            string      `json:"external_id"`
	Provider              string      `json:"provider"`
	Name                  string      `json:"name"`
	Email                 string      `json:"email"`
	VerifiedEmail         bool        `json:"verified_email"`
	Tier                  string      `json:"tier"`
	Avatar                pgtype.Text `json:"avatar"`
	SubscriptionExpiresAt *time.Time  `json:"subscription_expires_at"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, createUser,
		arg.ExternalID,
		arg.Provider,
		arg.Name,
		arg.Email,
		arg.VerifiedEmail,
		arg.Tier,
		arg.Avatar,
		arg.SubscriptionExpiresAt,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.ExternalID,
		&i.Provider,
		&i.Name,
		&i.Email,
		&i.VerifiedEmail,
		&i.Tier,
		&i.Avatar,
		&i.SubscriptionExpiresAt,
		&i.SynchronizedAt,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, external_id, provider, name, email, verified_email, tier, avatar, subscription_expires_at, synchronized_at, created_at, updated_at, deleted_at
FROM users
WHERE email = $1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRow(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.ExternalID,
		&i.Provider,
		&i.Name,
		&i.Email,
		&i.VerifiedEmail,
		&i.Tier,
		&i.Avatar,
		&i.SubscriptionExpiresAt,
		&i.SynchronizedAt,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const getUserByID = `-- name: GetUserByID :one
SELECT id, external_id, provider, name, email, verified_email, tier, avatar, subscription_expires_at, synchronized_at, created_at, updated_at, deleted_at
FROM users
WHERE id = $1
`

func (q *Queries) GetUserByID(ctx context.Context, id uuid.UUID) (User, error) {
	row := q.db.QueryRow(ctx, getUserByID, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.ExternalID,
		&i.Provider,
		&i.Name,
		&i.Email,
		&i.VerifiedEmail,
		&i.Tier,
		&i.Avatar,
		&i.SubscriptionExpiresAt,
		&i.SynchronizedAt,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const listPremiumActiveUsersWithAccounts = `-- name: ListPremiumActiveUsersWithAccounts :many
SELECT users.id, users.external_id, users.provider, users.name, users.email, users.verified_email, users.tier, users.avatar, users.subscription_expires_at, users.synchronized_at, users.created_at, users.updated_at, users.deleted_at,
  accounts.id, accounts.external_id, accounts.name, accounts.type, accounts.created_at, accounts.updated_at, accounts.deleted_at, accounts.user_id, accounts.institution_id
FROM users
  JOIN accounts ON accounts.user_id = users.id
WHERE tier IN ('PREMIUM', 'TRIAL')
  AND subscription_expires_at > NOW()
  AND users.deleted_at IS NULL
  AND accounts.deleted_at IS NULL
`

type ListPremiumActiveUsersWithAccountsRow struct {
	User    User    `json:"user"`
	Account Account `json:"account"`
}

func (q *Queries) ListPremiumActiveUsersWithAccounts(ctx context.Context) ([]ListPremiumActiveUsersWithAccountsRow, error) {
	rows, err := q.db.Query(ctx, listPremiumActiveUsersWithAccounts)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListPremiumActiveUsersWithAccountsRow
	for rows.Next() {
		var i ListPremiumActiveUsersWithAccountsRow
		if err := rows.Scan(
			&i.User.ID,
			&i.User.ExternalID,
			&i.User.Provider,
			&i.User.Name,
			&i.User.Email,
			&i.User.VerifiedEmail,
			&i.User.Tier,
			&i.User.Avatar,
			&i.User.SubscriptionExpiresAt,
			&i.User.SynchronizedAt,
			&i.User.CreatedAt,
			&i.User.UpdatedAt,
			&i.User.DeletedAt,
			&i.Account.ID,
			&i.Account.ExternalID,
			&i.Account.Name,
			&i.Account.Type,
			&i.Account.CreatedAt,
			&i.Account.UpdatedAt,
			&i.Account.DeletedAt,
			&i.Account.UserID,
			&i.Account.InstitutionID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listUsers = `-- name: ListUsers :many
SELECT id, external_id, provider, name, email, verified_email, tier, avatar, subscription_expires_at, synchronized_at, created_at, updated_at, deleted_at
FROM users
WHERE deleted_at IS NULL
`

func (q *Queries) ListUsers(ctx context.Context) ([]User, error) {
	rows, err := q.db.Query(ctx, listUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.ExternalID,
			&i.Provider,
			&i.Name,
			&i.Email,
			&i.VerifiedEmail,
			&i.Tier,
			&i.Avatar,
			&i.SubscriptionExpiresAt,
			&i.SynchronizedAt,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateUser = `-- name: UpdateUser :one
UPDATE users
SET external_id = $2,
  provider = $3,
  name = $4,
  email = $5,
  verified_email = $6,
  tier = $7,
  avatar = $8,
  subscription_expires_at = $9,
  synchronized_at = $10
WHERE id = $1
RETURNING id, external_id, provider, name, email, verified_email, tier, avatar, subscription_expires_at, synchronized_at, created_at, updated_at, deleted_at
`

type UpdateUserParams struct {
	ID                    uuid.UUID   `json:"id"`
	ExternalID            string      `json:"external_id"`
	Provider              string      `json:"provider"`
	Name                  string      `json:"name"`
	Email                 string      `json:"email"`
	VerifiedEmail         bool        `json:"verified_email"`
	Tier                  string      `json:"tier"`
	Avatar                pgtype.Text `json:"avatar"`
	SubscriptionExpiresAt *time.Time  `json:"subscription_expires_at"`
	SynchronizedAt        *time.Time  `json:"synchronized_at"`
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, updateUser,
		arg.ID,
		arg.ExternalID,
		arg.Provider,
		arg.Name,
		arg.Email,
		arg.VerifiedEmail,
		arg.Tier,
		arg.Avatar,
		arg.SubscriptionExpiresAt,
		arg.SynchronizedAt,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.ExternalID,
		&i.Provider,
		&i.Name,
		&i.Email,
		&i.VerifiedEmail,
		&i.Tier,
		&i.Avatar,
		&i.SubscriptionExpiresAt,
		&i.SynchronizedAt,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return i, err
}
