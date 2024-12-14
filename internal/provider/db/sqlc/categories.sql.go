// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: categories.sql

package sqlc

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const countSearchCategories = `-- name: CountSearchCategories :one
SELECT COUNT(*)
FROM categories
WHERE (
    $1::string = ''
    OR levenshtein(unaccent(name), unaccent($1::string)) <= 2
  )
`

func (q *Queries) CountSearchCategories(ctx context.Context, search string) (int64, error) {
	row := q.db.QueryRow(ctx, countSearchCategories, search)
	var count int64
	err := row.Scan(&count)
	return count, err
}

type CreateManyCategoriesParams struct {
	ExternalID string `json:"external_id"`
	Name       string `json:"name"`
}

const listCategories = `-- name: ListCategories :many
SELECT id, external_id, name, created_at, updated_at
FROM categories
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

const searchCategories = `-- name: SearchCategories :many
SELECT id, external_id, name, created_at, updated_at,
  levenshtein(unaccent(name), unaccent($3)) AS distance
FROM categories
WHERE (
    $3 = ''
    OR levenshtein(unaccent(name), unaccent($3)) <= 2
  )
ORDER BY distance,
  name
LIMIT $1 OFFSET $2
`

type SearchCategoriesParams struct {
	Limit  int32  `json:"limit"`
	Offset int32  `json:"offset"`
	Search string `json:"search"`
}

type SearchCategoriesRow struct {
	ID         uuid.UUID `json:"id"`
	ExternalID string    `json:"external_id"`
	Name       string    `json:"name"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Distance   int32     `json:"distance"`
}

func (q *Queries) SearchCategories(ctx context.Context, arg SearchCategoriesParams) ([]SearchCategoriesRow, error) {
	rows, err := q.db.Query(ctx, searchCategories, arg.Limit, arg.Offset, arg.Search)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []SearchCategoriesRow
	for rows.Next() {
		var i SearchCategoriesRow
		if err := rows.Scan(
			&i.ID,
			&i.ExternalID,
			&i.Name,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Distance,
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
