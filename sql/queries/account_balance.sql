-- name: CreateAccountBalances :copyfrom
INSERT INTO account_balances (amount, account_id)
VALUES ($1, $2);
-- name: GetUserBalanceOnDate :one
SELECT COALESCE(SUM(ab.amount), 0)::bigint AS total_balance
FROM accounts a
  JOIN LATERAL (
    SELECT ab.amount
    FROM account_balances ab
    WHERE ab.account_id = a.id
      AND ab.created_at <= sqlc.arg(date)::timestamptz
      AND ab.deleted_at IS NULL
    ORDER BY ab.created_at DESC
    LIMIT 1
  ) ab ON TRUE
  JOIN user_institutions ui ON a.user_institution_id = ui.id
WHERE a.type = 'BANK'
  AND ui.user_id = $1
  AND a.deleted_at IS NULL
  AND ui.deleted_at IS NULL;