-- name: CreateTransactions :copyfrom
INSERT INTO transactions (
    external_id,
    name,
    amount,
    payment_method,
    date,
    user_id,
    account_id,
    institution_id,
    category_id
  )
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);