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
    external_id,
    provider,
    name,
    email,
    verified_email,
    tier,
    avatar
  )
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;
-- name: UpdateUser :one
UPDATE users
SET external_id = $2,
  provider = $3,
  name = $4,
  email = $5,
  verified_email = $6,
  tier = $7,
  avatar = $8,
  subscription_expires_at = $9,
  synchronized_at = $10
WHERE id = $1
RETURNING *;
-- name: ListUsers :many
SELECT *
FROM users;
-- name: ListPremiumActiveUsersWithAccounts :many
SELECT sqlc.embed(users),
  sqlc.embed(accounts)
FROM users
  JOIN accounts ON accounts.user_id = users.id
WHERE tier IN ('PREMIUM', 'TRIAL')
  AND subscription_expires_at > NOW();