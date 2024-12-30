-- name: ListCategories :many
SELECT *
FROM categories;
-- name: CreateCategories :copyfrom
INSERT INTO categories (external_id, name)
VALUES ($1, $2);
-- name: ListCategoriesByExternalIDs :many
SELECT *
FROM categories
WHERE external_id = ANY(sqlc.arg(ids)::text []);