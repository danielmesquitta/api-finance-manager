-- name: CreateUserInstitution :one
INSERT INTO user_institutions (external_id, user_id, institution_id)
VALUES ($1, $2, $3)
RETURNING *;
-- name: GetUserInstitutionByExternalID :one
SELECT *
FROM user_institutions
WHERE external_id = $1
  AND deleted_at IS NULL;