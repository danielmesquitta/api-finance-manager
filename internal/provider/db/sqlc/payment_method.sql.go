// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: payment_method.sql

package sqlc

type CreatePaymentMethodsParams struct {
	ExternalID string `json:"external_id"`
	Name       string `json:"name"`
}
