// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: users.sql

package sqlc

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (
    name,
    email,
    tier,
    avatar,
    subscription_expires_at
  )
VALUES ($1, $2, $3, $4, $5)
RETURNING id, name, email, tier, avatar, subscription_expires_at, synchronized_at, created_at, updated_at
`

type CreateUserParams struct {
	Name                  string      `json:"name"`
	Email                 string      `json:"email"`
	Tier                  Tier        `json:"tier"`
	Avatar                pgtype.Text `json:"avatar"`
	SubscriptionExpiresAt time.Time   `json:"subscription_expires_at"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, createUser,
		arg.Name,
		arg.Email,
		arg.Tier,
		arg.Avatar,
		arg.SubscriptionExpiresAt,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.Tier,
		&i.Avatar,
		&i.SubscriptionExpiresAt,
		&i.SynchronizedAt,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, name, email, tier, avatar, subscription_expires_at, synchronized_at, created_at, updated_at
FROM users
WHERE email = $1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRow(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.Tier,
		&i.Avatar,
		&i.SubscriptionExpiresAt,
		&i.SynchronizedAt,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getUserByID = `-- name: GetUserByID :one
SELECT id, name, email, tier, avatar, subscription_expires_at, synchronized_at, created_at, updated_at
FROM users
WHERE id = $1
`

func (q *Queries) GetUserByID(ctx context.Context, id uuid.UUID) (User, error) {
	row := q.db.QueryRow(ctx, getUserByID, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.Tier,
		&i.Avatar,
		&i.SubscriptionExpiresAt,
		&i.SynchronizedAt,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateUser = `-- name: UpdateUser :one
UPDATE users
SET name = $2,
  email = $3,
  tier = $4,
  avatar = $5,
  subscription_expires_at = $6,
  synchronized_at = $7,
  updated_at = NOW()
WHERE id = $1
RETURNING id, name, email, tier, avatar, subscription_expires_at, synchronized_at, created_at, updated_at
`

type UpdateUserParams struct {
	ID                    uuid.UUID   `json:"id"`
	Name                  string      `json:"name"`
	Email                 string      `json:"email"`
	Tier                  Tier        `json:"tier"`
	Avatar                pgtype.Text `json:"avatar"`
	SubscriptionExpiresAt time.Time   `json:"subscription_expires_at"`
	SynchronizedAt        time.Time   `json:"synchronized_at"`
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, updateUser,
		arg.ID,
		arg.Name,
		arg.Email,
		arg.Tier,
		arg.Avatar,
		arg.SubscriptionExpiresAt,
		arg.SynchronizedAt,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.Tier,
		&i.Avatar,
		&i.SubscriptionExpiresAt,
		&i.SynchronizedAt,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
