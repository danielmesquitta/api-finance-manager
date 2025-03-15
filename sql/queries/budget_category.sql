-- name: CreateBudgetCategories :copyfrom
INSERT INTO budget_categories (amount, budget_id, category_id)
VALUES ($1, $2, $3);
-- name: DeleteBudgetCategories :exec
UPDATE budget_categories
SET deleted_at = NOW()
WHERE budget_categories.budget_id = $1
  AND budget_categories.deleted_at IS NULL;
-- name: ListBudgetCategories :many
SELECT sqlc.embed(bc),
  sqlc.embed(tc)
FROM budget_categories bc
  JOIN transaction_categories tc ON bc.category_id = tc.id
WHERE budget_id = $1
  AND bc.deleted_at IS NULL
ORDER BY tc.name ASC;
-- name: GetBudgetCategory :one
SELECT sqlc.embed(bc),
  sqlc.embed(tc)
FROM budget_categories bc
  JOIN transaction_categories tc ON bc.category_id = tc.id
  JOIN budgets b ON bc.budget_id = b.id
WHERE b.user_id = $1
  AND b.date <= $2
  AND bc.deleted_at IS NULL
LIMIT 1;