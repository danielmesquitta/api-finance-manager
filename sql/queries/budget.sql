-- name: GetBudgetByUserID :one
SELECT *
FROM budgets
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