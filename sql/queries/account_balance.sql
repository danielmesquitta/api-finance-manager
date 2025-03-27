-- name: CreateAccountBalances :copyfrom
INSERT INTO account_balances (amount, account_id)
VALUES ($1, $2);