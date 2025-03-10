-- name: GetUserByID :one
SELECT *
FROM users
WHERE id = $1;
-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = $1;
-- name: CreateUser :one
INSERT INTO users (
    auth_id,
    open_finance_id,
    provider,
    name,
    email,
    verified_email,
    tier,
    avatar,
    subscription_expires_at
  )
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;
-- name: UpdateUser :one
UPDATE users
SET auth_id = $2,
  open_finance_id = $3,
  provider = $4,
  name = $5,
  email = $6,
  verified_email = $7,
  tier = $8,
  avatar = $9,
  subscription_expires_at = $10,
  synchronized_at = $11
WHERE id = $1
RETURNING *;
-- name: UpdateUserSynchronizedAt :exec
UPDATE users
SET synchronized_at = $2
WHERE id = $1;
-- name: ListUsers :many
SELECT *
FROM users
WHERE deleted_at IS NULL;
-- name: DeleteUser :exec
UPDATE users
SET deleted_at = NOW(),
  name = $2,
  email = $3
WHERE id = $1;
-- name: DestroyUser :exec
DELETE FROM users
WHERE id = $1;