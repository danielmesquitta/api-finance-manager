// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: budget.sql

package sqlc

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createBudget = `-- name: CreateBudget :one
INSERT INTO budgets (amount, date, user_id)
VALUES ($1, $2, $3)
RETURNING id, amount, date, created_at, updated_at, deleted_at, user_id
`

type CreateBudgetParams struct {
	Amount int64     `json:"amount"`
	Date   time.Time `json:"date"`
	UserID uuid.UUID `json:"user_id"`
}

func (q *Queries) CreateBudget(ctx context.Context, arg CreateBudgetParams) (Budget, error) {
	row := q.db.QueryRow(ctx, createBudget, arg.Amount, arg.Date, arg.UserID)
	var i Budget
	err := row.Scan(
		&i.ID,
		&i.Amount,
		&i.Date,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.UserID,
	)
	return i, err
}

type CreateBudgetCategoriesParams struct {
	Amount     int64     `json:"amount"`
	BudgetID   uuid.UUID `json:"budget_id"`
	CategoryID uuid.UUID `json:"category_id"`
}

const deleteBudgetCategories = `-- name: DeleteBudgetCategories :exec
UPDATE budget_categories
SET deleted_at = NOW()
FROM budgets
WHERE budget_categories.budget_id = budgets.id
  AND budgets.user_id = $1
  AND deleted_at IS NULL
`

func (q *Queries) DeleteBudgetCategories(ctx context.Context, userID uuid.UUID) error {
	_, err := q.db.Exec(ctx, deleteBudgetCategories, userID)
	return err
}

const deleteBudgets = `-- name: DeleteBudgets :exec
UPDATE budgets
SET deleted_at = NOW()
WHERE user_id = $1
  AND deleted_at IS NULL
`

func (q *Queries) DeleteBudgets(ctx context.Context, userID uuid.UUID) error {
	_, err := q.db.Exec(ctx, deleteBudgets, userID)
	return err
}

const getBudget = `-- name: GetBudget :one
SELECT id, amount, date, created_at, updated_at, deleted_at, user_id
FROM budgets
WHERE user_id = $1
  AND date <= $2
  AND deleted_at IS NULL
ORDER BY date ASC
LIMIT 1
`

type GetBudgetParams struct {
	UserID uuid.UUID `json:"user_id"`
	Date   time.Time `json:"date"`
}

func (q *Queries) GetBudget(ctx context.Context, arg GetBudgetParams) (Budget, error) {
	row := q.db.QueryRow(ctx, getBudget, arg.UserID, arg.Date)
	var i Budget
	err := row.Scan(
		&i.ID,
		&i.Amount,
		&i.Date,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.UserID,
	)
	return i, err
}

const listBudgetCategories = `-- name: ListBudgetCategories :many
SELECT budget_categories.id, budget_categories.amount, budget_categories.created_at, budget_categories.updated_at, budget_categories.deleted_at, budget_categories.budget_id, budget_categories.category_id,
  categories.id, categories.external_id, categories.name, categories.created_at, categories.updated_at, categories.deleted_at
FROM budget_categories
  JOIN categories ON budget_categories.category_id = categories.id
WHERE budget_id = $1
  AND deleted_at IS NULL
ORDER BY categories.name ASC
`

type ListBudgetCategoriesRow struct {
	BudgetCategory BudgetCategory `json:"budget_category"`
	Category       Category       `json:"category"`
}

func (q *Queries) ListBudgetCategories(ctx context.Context, budgetID uuid.UUID) ([]ListBudgetCategoriesRow, error) {
	rows, err := q.db.Query(ctx, listBudgetCategories, budgetID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListBudgetCategoriesRow
	for rows.Next() {
		var i ListBudgetCategoriesRow
		if err := rows.Scan(
			&i.BudgetCategory.ID,
			&i.BudgetCategory.Amount,
			&i.BudgetCategory.CreatedAt,
			&i.BudgetCategory.UpdatedAt,
			&i.BudgetCategory.DeletedAt,
			&i.BudgetCategory.BudgetID,
			&i.BudgetCategory.CategoryID,
			&i.Category.ID,
			&i.Category.ExternalID,
			&i.Category.Name,
			&i.Category.CreatedAt,
			&i.Category.UpdatedAt,
			&i.Category.DeletedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateBudget = `-- name: UpdateBudget :exec
UPDATE budgets
SET amount = $1
WHERE user_id = $2
  AND date = $3
  AND deleted_at IS NULL
`

type UpdateBudgetParams struct {
	Amount int64     `json:"amount"`
	UserID uuid.UUID `json:"user_id"`
	Date   time.Time `json:"date"`
}

func (q *Queries) UpdateBudget(ctx context.Context, arg UpdateBudgetParams) error {
	_, err := q.db.Exec(ctx, updateBudget, arg.Amount, arg.UserID, arg.Date)
	return err
}
