-- name: GetBudget :one
SELECT *
FROM budgets
WHERE user_id = $1
  AND date <= $2
  AND deleted_at IS NULL
ORDER BY date DESC
LIMIT 1;
-- name: GetBudgetCategory :one
SELECT sqlc.embed(budget_categories),
  sqlc.embed(transaction_categories)
FROM budget_categories
  JOIN transaction_categories ON budget_categories.category_id = transaction_categories.id
  AND transaction_categories.deleted_at IS NULL
  JOIN budgets ON budget_categories.budget_id = budgets.id
  AND budgets.deleted_at IS NULL
WHERE budgets.user_id = $1
  AND budgets.date <= $2
  AND budget_categories.deleted_at IS NULL
LIMIT 1;
-- name: ListBudgetCategories :many
SELECT sqlc.embed(budget_categories),
  sqlc.embed(transaction_categories)
FROM budget_categories
  JOIN transaction_categories ON budget_categories.category_id = transaction_categories.id
WHERE budget_id = $1
  AND budget_categories.deleted_at IS NULL
ORDER BY transaction_categories.name ASC;
-- name: CreateBudget :one
INSERT INTO budgets (amount, date, user_id)
VALUES ($1, $2, $3)
RETURNING *;
-- name: UpdateBudget :exec
UPDATE budgets
SET amount = $1
WHERE user_id = $2
  AND date = $3
  AND deleted_at IS NULL;
-- name: CreateBudgetCategories :copyfrom
INSERT INTO budget_categories (amount, budget_id, category_id)
VALUES ($1, $2, $3);
-- name: DeleteBudgetCategories :exec
UPDATE budget_categories
SET deleted_at = NOW()
FROM budgets
WHERE budget_categories.budget_id = budgets.id
  AND budgets.user_id = $1
  AND budget_categories.deleted_at IS NULL;
-- name: DeleteBudgets :exec
UPDATE budgets
SET deleted_at = NOW()
WHERE user_id = $1
  AND deleted_at IS NULL;