-- name: ListCategories :many
SELECT *
FROM categories
WHERE deleted_at IS NULL;
-- name: ListCategoriesByExternalIDs :many
SELECT *
FROM categories
WHERE external_id = ANY(sqlc.arg(ids)::text [])
  AND deleted_at IS NULL;
-- name: CreateCategories :copyfrom
INSERT INTO categories (external_id, name)
VALUES ($1, $2);