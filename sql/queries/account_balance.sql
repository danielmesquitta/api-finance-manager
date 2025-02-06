-- name: CreateAccounts :exec
INSERT INTO account_balances (amount, account_id)
VALUES ($1, $2);