// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: account_balance.sql

package sqlc

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type CreateAccountBalancesParams struct {
	Amount    int64     `json:"amount"`
	UserID    uuid.UUID `json:"user_id"`
	AccountID uuid.UUID `json:"account_id"`
}

const getUserBalanceOnDate = `-- name: GetUserBalanceOnDate :one
SELECT COALESCE(SUM(ab.amount), 0)::bigint AS total_balance
FROM accounts a
  JOIN LATERAL (
    SELECT ab.amount
    FROM account_balances ab
    WHERE ab.account_id = a.id
      AND ab.created_at <= $2::timestamptz
    ORDER BY ab.created_at DESC
    LIMIT 1
  ) ab ON TRUE
  JOIN user_institutions ui ON a.user_institution_id = ui.id
WHERE a.type = 'BANK'
  AND ui.user_id = $1
`

type GetUserBalanceOnDateParams struct {
	UserID uuid.UUID `json:"user_id"`
	Date   time.Time `json:"date"`
}

func (q *Queries) GetUserBalanceOnDate(ctx context.Context, arg GetUserBalanceOnDateParams) (int64, error) {
	row := q.db.QueryRow(ctx, getUserBalanceOnDate, arg.UserID, arg.Date)
	var total_balance int64
	err := row.Scan(&total_balance)
	return total_balance, err
}
