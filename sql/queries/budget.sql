-- name: GetBudgetByUserID :one
SELECT *
FROM budgets
WHERE user_id = $1;
-- name: GetBudgetWithCategoriesByUserID :many
SELECT sqlc.embed(budgets),
  sqlc.embed(budget_categories),
  sqlc.embed(categories)
FROM budgets
  JOIN budget_categories ON budgets.id = budget_categories.budget_id
  JOIN categories ON budget_categories.category_id = categories.id
WHERE user_id = $1;
-- name: CreateBudget :one
INSERT INTO budgets (amount, user_id)
VALUES ($1, $2)
RETURNING *;
-- name: UpdateBudget :exec
UPDATE budgets
SET amount = $1
WHERE user_id = $2;
-- name: CreateBudgetCategories :copyfrom
INSERT INTO budget_categories (amount, budget_id, category_id)
VALUES ($1, $2, $3);
-- name: DeleteBudgetCategoriesByBudgetID :exec
DELETE FROM budget_categories
WHERE budget_id = $1;
-- name: DeleteBudgetByID :exec
DELETE FROM budgets
WHERE id = $1;