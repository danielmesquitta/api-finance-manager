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
    name,
    email,
    tier,
    avatar,
    subscription_expires_at
  )
VALUES ($1, $2, $3, $4, $5)
RETURNING *;
-- name: UpdateUser :one
UPDATE users
SET name = $2,
  email = $3,
  tier = $4,
  avatar = $5,
  subscription_expires_at = $6,
  synchronized_at = $7
WHERE id = $1
RETURNING *;
-- name: UpdateUserSynchronizedAt :exec
UPDATE users
SET synchronized_at = $2
WHERE id = $1;
-- name: DeleteUser :exec
UPDATE users
SET deleted_at = NOW(),
  name = $2,
  email = $3
WHERE id = $1;
-- name: DestroyUser :exec
DELETE FROM users
WHERE id = $1;