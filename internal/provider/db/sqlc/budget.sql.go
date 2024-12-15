// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: budget.sql

package sqlc

import (
	"context"

	"github.com/google/uuid"
)

const createBudget = `-- name: CreateBudget :one
INSERT INTO budgets (amount, user_id)
VALUES ($1, $2)
RETURNING id, amount, created_at, updated_at, user_id
`

type CreateBudgetParams struct {
	Amount float64   `json:"amount"`
	UserID uuid.UUID `json:"user_id"`
}

func (q *Queries) CreateBudget(ctx context.Context, arg CreateBudgetParams) (Budget, error) {
	row := q.db.QueryRow(ctx, createBudget, arg.Amount, arg.UserID)
	var i Budget
	err := row.Scan(
		&i.ID,
		&i.Amount,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserID,
	)
	return i, err
}

type CreateBudgetCategoriesParams struct {
	Amount     float64   `json:"amount"`
	BudgetID   uuid.UUID `json:"budget_id"`
	CategoryID uuid.UUID `json:"category_id"`
}

const deleteBudgetCategoriesByBudgetID = `-- name: DeleteBudgetCategoriesByBudgetID :exec
DELETE FROM budget_categories
WHERE budget_id = $1
`

func (q *Queries) DeleteBudgetCategoriesByBudgetID(ctx context.Context, budgetID uuid.UUID) error {
	_, err := q.db.Exec(ctx, deleteBudgetCategoriesByBudgetID, budgetID)
	return err
}

const getBudgetByUserID = `-- name: GetBudgetByUserID :one
SELECT id, amount, created_at, updated_at, user_id
FROM budgets
WHERE user_id = $1
`

func (q *Queries) GetBudgetByUserID(ctx context.Context, userID uuid.UUID) (Budget, error) {
	row := q.db.QueryRow(ctx, getBudgetByUserID, userID)
	var i Budget
	err := row.Scan(
		&i.ID,
		&i.Amount,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserID,
	)
	return i, err
}

const updateBudget = `-- name: UpdateBudget :exec
UPDATE budgets
SET amount = $1
WHERE user_id = $2
`

type UpdateBudgetParams struct {
	Amount float64   `json:"amount"`
	UserID uuid.UUID `json:"user_id"`
}

func (q *Queries) UpdateBudget(ctx context.Context, arg UpdateBudgetParams) error {
	_, err := q.db.Exec(ctx, updateBudget, arg.Amount, arg.UserID)
	return err
}
