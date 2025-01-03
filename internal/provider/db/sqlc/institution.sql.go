// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: institution.sql

package sqlc

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

type CreateInstitutionsParams struct {
	ExternalID string      `json:"external_id"`
	Name       string      `json:"name"`
	Logo       pgtype.Text `json:"logo"`
}

const getInstitutionByExternalID = `-- name: GetInstitutionByExternalID :one
SELECT id, external_id, name, logo, created_at, updated_at, deleted_at
FROM institutions
WHERE external_id = $1
  AND deleted_at IS NULL
`

func (q *Queries) GetInstitutionByExternalID(ctx context.Context, externalID string) (Institution, error) {
	row := q.db.QueryRow(ctx, getInstitutionByExternalID, externalID)
	var i Institution
	err := row.Scan(
		&i.ID,
		&i.ExternalID,
		&i.Name,
		&i.Logo,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const listInstitutions = `-- name: ListInstitutions :many
SELECT id, external_id, name, logo, created_at, updated_at, deleted_at
FROM institutions
WHERE deleted_at IS NULL
`

func (q *Queries) ListInstitutions(ctx context.Context) ([]Institution, error) {
	rows, err := q.db.Query(ctx, listInstitutions)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Institution
	for rows.Next() {
		var i Institution
		if err := rows.Scan(
			&i.ID,
			&i.ExternalID,
			&i.Name,
			&i.Logo,
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
