-- name: ListAccountsByUserID :many
SELECT *
FROM accounts
WHERE user_id = $1
  AND deleted_at IS NULL;
-- name: CreateAccounts :copyfrom
INSERT INTO accounts (external_id, name, type, user_id, institution_id)
VALUES ($1, $2, $3, $4, $5);