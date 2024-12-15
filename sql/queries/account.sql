-- name: ListAccountsByUserID :many
SELECT *
FROM accounts
WHERE user_id = $1;