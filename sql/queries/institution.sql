-- name: ListInstitutions :many
SELECT *
FROM institutions;
-- name: GetInstitutionByExternalID :one
SELECT *
FROM institutions
WHERE external_id = $1;
-- name: CreateInstitutions :copyfrom
INSERT INTO institutions (external_id, name, logo)
VALUES ($1, $2, $3);