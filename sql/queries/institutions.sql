-- name: ListInstitutions :many
SELECT *
FROM institutions;
-- name: CreateManyInstitutions :copyfrom
INSERT INTO institutions (external_id, name, logo)
VALUES ($1, $2, $3);