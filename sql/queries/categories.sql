-- name: ListCategories :many
SELECT *
FROM categories;
-- name: CreateManyCategories :copyfrom
INSERT INTO categories (external_id, name)
VALUES ($1, $2);