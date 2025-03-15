// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
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
WHERE budget_categories.budget_id = $1
  AND budget_categories.deleted_at IS NULL
`

func (q *Queries) DeleteBudgetCategories(ctx context.Context, budgetID uuid.UUID) error {
	_, err := q.db.Exec(ctx, deleteBudgetCategories, budgetID)
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
ORDER BY date DESC
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

const getBudgetCategory = `-- name: GetBudgetCategory :one
SELECT bc.id, bc.amount, bc.created_at, bc.updated_at, bc.deleted_at, bc.budget_id, bc.category_id,
  tc.id, tc.external_id, tc.name, tc.created_at, tc.updated_at, tc.deleted_at
FROM budget_categories bc
  JOIN transaction_categories tc ON bc.category_id = tc.id
  JOIN budgets b ON bc.budget_id = b.id
WHERE b.user_id = $1
  AND b.date <= $2
  AND bc.deleted_at IS NULL
LIMIT 1
`

type GetBudgetCategoryParams struct {
	UserID uuid.UUID `json:"user_id"`
	Date   time.Time `json:"date"`
}

type GetBudgetCategoryRow struct {
	BudgetCategory      BudgetCategory      `json:"budget_category"`
	TransactionCategory TransactionCategory `json:"transaction_category"`
}

func (q *Queries) GetBudgetCategory(ctx context.Context, arg GetBudgetCategoryParams) (GetBudgetCategoryRow, error) {
	row := q.db.QueryRow(ctx, getBudgetCategory, arg.UserID, arg.Date)
	var i GetBudgetCategoryRow
	err := row.Scan(
		&i.BudgetCategory.ID,
		&i.BudgetCategory.Amount,
		&i.BudgetCategory.CreatedAt,
		&i.BudgetCategory.UpdatedAt,
		&i.BudgetCategory.DeletedAt,
		&i.BudgetCategory.BudgetID,
		&i.BudgetCategory.CategoryID,
		&i.TransactionCategory.ID,
		&i.TransactionCategory.ExternalID,
		&i.TransactionCategory.Name,
		&i.TransactionCategory.CreatedAt,
		&i.TransactionCategory.UpdatedAt,
		&i.TransactionCategory.DeletedAt,
	)
	return i, err
}

const listBudgetCategories = `-- name: ListBudgetCategories :many
SELECT bc.id, bc.amount, bc.created_at, bc.updated_at, bc.deleted_at, bc.budget_id, bc.category_id,
  tc.id, tc.external_id, tc.name, tc.created_at, tc.updated_at, tc.deleted_at
FROM budget_categories bc
  JOIN transaction_categories tc ON bc.category_id = tc.id
WHERE budget_id = $1
  AND bc.deleted_at IS NULL
ORDER BY tc.name ASC
`

type ListBudgetCategoriesRow struct {
	BudgetCategory      BudgetCategory      `json:"budget_category"`
	TransactionCategory TransactionCategory `json:"transaction_category"`
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
			&i.TransactionCategory.ID,
			&i.TransactionCategory.ExternalID,
			&i.TransactionCategory.Name,
			&i.TransactionCategory.CreatedAt,
			&i.TransactionCategory.UpdatedAt,
			&i.TransactionCategory.DeletedAt,
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
