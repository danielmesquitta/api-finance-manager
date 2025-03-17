-- name: CreateUserAuthProvider :exec
INSERT INTO user_auth_providers (external_id, provider, verified_email, user_id)
VALUES ($1, $2, $3, $4);
-- name: UpdateUserAuthProvider :exec
UPDATE user_auth_providers
SET verified_email = $2
WHERE id = $1;
-- name: GetUserAuthProvider :one
SELECT *
FROM user_auth_providers
WHERE user_id = $1
  AND provider = $2
  AND deleted_at IS NULL;