-- name: CreateAccounts :copyfrom
INSERT INTO accounts (
    id,
    external_id,
    name,
    type,
    user_institution_id
  )
VALUES ($1, $2, $3, $4, $5);