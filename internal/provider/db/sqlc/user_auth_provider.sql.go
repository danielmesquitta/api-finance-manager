// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: user_auth_provider.sql

package sqlc

import (
	"context"

	"github.com/google/uuid"
)

const createUserAuthProvider = `-- name: CreateUserAuthProvider :exec
INSERT INTO user_auth_providers (external_id, provider, verified_email, user_id)
VALUES ($1, $2, $3, $4)
`

type CreateUserAuthProviderParams struct {
	ExternalID    string    `json:"external_id"`
	Provider      string    `json:"provider"`
	VerifiedEmail bool      `json:"verified_email"`
	UserID        uuid.UUID `json:"user_id"`
}

func (q *Queries) CreateUserAuthProvider(ctx context.Context, arg CreateUserAuthProviderParams) error {
	_, err := q.db.Exec(ctx, createUserAuthProvider,
		arg.ExternalID,
		arg.Provider,
		arg.VerifiedEmail,
		arg.UserID,
	)
	return err
}

const getUserAuthProvider = `-- name: GetUserAuthProvider :one
SELECT id, external_id, provider, verified_email, created_at, updated_at, deleted_at, user_id
FROM user_auth_providers
WHERE user_id = $1
  AND provider = $2
  AND deleted_at IS NULL
`

type GetUserAuthProviderParams struct {
	UserID   uuid.UUID `json:"user_id"`
	Provider string    `json:"provider"`
}

func (q *Queries) GetUserAuthProvider(ctx context.Context, arg GetUserAuthProviderParams) (UserAuthProvider, error) {
	row := q.db.QueryRow(ctx, getUserAuthProvider, arg.UserID, arg.Provider)
	var i UserAuthProvider
	err := row.Scan(
		&i.ID,
		&i.ExternalID,
		&i.Provider,
		&i.VerifiedEmail,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.UserID,
	)
	return i, err
}

const updateUserAuthProvider = `-- name: UpdateUserAuthProvider :exec
UPDATE user_auth_providers
SET verified_email = $2
WHERE id = $1
`

type UpdateUserAuthProviderParams struct {
	ID            uuid.UUID `json:"id"`
	VerifiedEmail bool      `json:"verified_email"`
}

func (q *Queries) UpdateUserAuthProvider(ctx context.Context, arg UpdateUserAuthProviderParams) error {
	_, err := q.db.Exec(ctx, updateUserAuthProvider, arg.ID, arg.VerifiedEmail)
	return err
}
