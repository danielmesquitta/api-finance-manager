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
    subscription_expires_at
  )
VALUES ($1, $2, $3, $4)
RETURNING *;
-- name: UpdateUser :one
UPDATE users
SET name = $2,
  email = $3,
  tier = $4,
  subscription_expires_at = $5,
  synchronized_at = $6
WHERE id = $1
RETURNING *;