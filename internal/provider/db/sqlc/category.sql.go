// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: category.sql

package sqlc

import (
	"context"

	"github.com/google/uuid"
)

const countCategoriesByIDs = `-- name: CountCategoriesByIDs :one
SELECT COUNT(*)
FROM categories
WHERE id = ANY($1::uuid [])
  AND deleted_at IS NULL
`

func (q *Queries) CountCategoriesByIDs(ctx context.Context, ids []uuid.UUID) (int64, error) {
	row := q.db.QueryRow(ctx, countCategoriesByIDs, ids)
	var count int64
	err := row.Scan(&count)
	return count, err
}

type CreateCategoriesParams struct {
	ExternalID string `json:"external_id"`
	Name       string `json:"name"`
}

const getCategory = `-- name: GetCategory :one
SELECT id, external_id, name, created_at, updated_at, deleted_at
FROM categories
WHERE id = $1
  AND deleted_at IS NULL
`

func (q *Queries) GetCategory(ctx context.Context, id uuid.UUID) (Category, error) {
	row := q.db.QueryRow(ctx, getCategory, id)
	var i Category
	err := row.Scan(
		&i.ID,
		&i.ExternalID,
		&i.Name,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const listCategories = `-- name: ListCategories :many
SELECT id, external_id, name, created_at, updated_at, deleted_at
FROM categories
WHERE deleted_at IS NULL
`

func (q *Queries) ListCategories(ctx context.Context) ([]Category, error) {
	rows, err := q.db.Query(ctx, listCategories)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Category
	for rows.Next() {
		var i Category
		if err := rows.Scan(
			&i.ID,
			&i.ExternalID,
			&i.Name,
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

const listCategoriesByExternalIDs = `-- name: ListCategoriesByExternalIDs :many
SELECT id, external_id, name, created_at, updated_at, deleted_at
FROM categories
WHERE external_id = ANY($1::text [])
  AND deleted_at IS NULL
`

func (q *Queries) ListCategoriesByExternalIDs(ctx context.Context, ids []string) ([]Category, error) {
	rows, err := q.db.Query(ctx, listCategoriesByExternalIDs, ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Category
	for rows.Next() {
		var i Category
		if err := rows.Scan(
			&i.ID,
			&i.ExternalID,
			&i.Name,
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
