// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: account.sql

package sqlc

import (
	"context"

	"github.com/google/uuid"
)

type CreateAccountsParams struct {
	ExternalID    string    `json:"external_id"`
	Name          string    `json:"name"`
	Type          string    `json:"type"`
	UserID        uuid.UUID `json:"user_id"`
	InstitutionID uuid.UUID `json:"institution_id"`
}

const listAccountsByUserID = `-- name: ListAccountsByUserID :many
SELECT id, external_id, name, type, created_at, updated_at, deleted_at, user_id, institution_id
FROM accounts
WHERE user_id = $1
  AND deleted_at IS NULL
`

func (q *Queries) ListAccountsByUserID(ctx context.Context, userID uuid.UUID) ([]Account, error) {
	rows, err := q.db.Query(ctx, listAccountsByUserID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Account
	for rows.Next() {
		var i Account
		if err := rows.Scan(
			&i.ID,
			&i.ExternalID,
			&i.Name,
			&i.Type,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
			&i.UserID,
			&i.InstitutionID,
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
