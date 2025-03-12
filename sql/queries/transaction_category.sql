-- name: ListTransactionCategories :many
SELECT *
FROM transaction_categories
WHERE deleted_at IS NULL;
-- name: CountTransactionCategoriesByIDs :one
SELECT COUNT(*)
FROM transaction_categories
WHERE id = ANY(sqlc.arg(ids)::uuid [])
  AND deleted_at IS NULL;
-- name: ListTransactionCategoriesByExternalIDs :many
SELECT *
FROM transaction_categories
WHERE external_id = ANY(sqlc.arg(ids)::text [])
  AND deleted_at IS NULL;
-- name: CreateTransactionCategories :copyfrom
INSERT INTO transaction_categories (external_id, name)
VALUES ($1, $2);
-- name: GetTransactionCategory :one
SELECT *
FROM transaction_categories
WHERE id = $1
  AND deleted_at IS NULL;
-- name: GetDefaultTransactionCategory :one
SELECT *
FROM transaction_categories
WHERE external_id = '99999999'
  AND deleted_at IS NULL;