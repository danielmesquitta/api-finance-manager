-- name: CreatePaymentMethods :copyfrom
INSERT INTO payment_methods (external_id, name)
VALUES ($1, $2);
-- name: GetPaymentMethod :one
SELECT *
FROM payment_methods
WHERE id = $1
  AND deleted_at IS NULL;
