// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: institution.sql

package sqlc

import (
	"context"
)

type CreateInstitutionsParams struct {
	ExternalID string  `json:"external_id"`
	Name       string  `json:"name"`
	Logo       *string `json:"logo"`
}

const getInstitutionByExternalID = `-- name: GetInstitutionByExternalID :one
SELECT id, external_id, name, logo, created_at, deleted_at
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
		&i.DeletedAt,
	)
	return i, err
}
