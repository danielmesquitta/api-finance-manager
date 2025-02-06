-- name: CreateTransactions :copyfrom
INSERT INTO transactions (
    external_id,
    name,
    amount,
    payment_method_id,
    date,
    user_id,
    account_id,
    institution_id,
    category_id
  )
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);
-- name: UpdateTransaction :exec
UPDATE transactions
SET name = $2,
  amount = $3,
  payment_method_id = $4,
  date = $5,
  account_id = $6,
  institution_id = $7,
  category_id = $8
WHERE id = $1
  AND user_id = $9
  AND deleted_at IS NULL;
-- name: GetTransaction :one
SELECT transactions.*,
  transaction_categories.name as category_name,
  institutions.name as institution_name,
  institutions.logo as institution_logo,
  payment_methods.name as payment_method_name
FROM transactions
  LEFT JOIN transaction_categories ON transactions.category_id = transaction_categories.id
  LEFT JOIN institutions ON transactions.institution_id = institutions.id
  LEFT JOIN payment_methods ON transactions.payment_method_id = payment_methods.id
WHERE transactions.id = $1
  AND user_id = $2
  AND transactions.deleted_at IS NULL;