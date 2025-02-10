-- name: CreateAccountBalances :copyfrom
INSERT INTO account_balances (amount, user_id, account_id)
VALUES ($1, $2, $3);
-- name: GetUserBalanceOnDate :one
WITH balances_on_date AS (
  SELECT DISTINCT ON (account_id) account_id,
    amount
  FROM account_balances
  WHERE user_id = $1
    AND created_at > sqlc.arg(date)::timestamptz
  ORDER BY account_id,
    created_at ASC
)
SELECT COALESCE(SUM(amount), 0)::bigint AS total_balance
FROM balances_on_date;