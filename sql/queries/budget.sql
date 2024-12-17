-- name: GetBudget :one
SELECT *
FROM budgets
WHERE user_id = $1
  AND date <= $2
ORDER BY date ASC
LIMIT 1;
-- name: GetBudgetCategories :many
SELECT sqlc.embed(budget_categories),
  sqlc.embed(categories)
FROM budget_categories
  JOIN categories ON budget_categories.category_id = categories.id
WHERE budget_id = $1
ORDER BY categories.name ASC;
-- name: CreateBudget :one
INSERT INTO budgets (amount, date, user_id)
VALUES ($1, $2, $3)
RETURNING *;
-- name: UpdateBudget :exec
UPDATE budgets
SET amount = $1
WHERE user_id = $2
  AND date = $3;
-- name: CreateBudgetCategories :copyfrom
INSERT INTO budget_categories (amount, budget_id, category_id)
VALUES ($1, $2, $3);
-- name: DeleteBudgetCategories :exec
DELETE FROM budget_categories USING budgets
WHERE budget_categories.budget_id = budgets.id
  AND budgets.user_id = $1;
-- name: DeleteBudgets :exec
DELETE FROM budgets
WHERE user_id = $1;