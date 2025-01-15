-- name: CreatePaymentMethods :copyfrom
INSERT INTO payment_methods (external_id, name)
VALUES ($1, $2);