-- name: ListCategories :many
SELECT *
FROM categories
WHERE deleted_at IS NULL;
-- name: CountCategoriesByIDs :one
SELECT COUNT(*)
FROM categories
WHERE id = ANY(sqlc.arg(ids)::uuid [])
  AND deleted_at IS NULL;
-- name: ListCategoriesByExternalIDs :many
SELECT *
FROM categories
WHERE external_id = ANY(sqlc.arg(ids)::text [])
  AND deleted_at IS NULL;
-- name: CreateCategories :copyfrom
INSERT INTO categories (external_id, name)
VALUES ($1, $2);
-- name: GetCategory :one
SELECT *
FROM categories
WHERE id = $1
  AND deleted_at IS NULL;