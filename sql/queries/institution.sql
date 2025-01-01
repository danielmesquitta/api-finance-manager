-- name: ListInstitutions :many
SELECT *
FROM institutions
WHERE deleted_at IS NULL;
-- name: GetInstitutionByExternalID :one
SELECT *
FROM institutions
WHERE external_id = $1
  AND deleted_at IS NULL;
-- name: CreateInstitutions :copyfrom
INSERT INTO institutions (external_id, name, logo)
VALUES ($1, $2, $3);