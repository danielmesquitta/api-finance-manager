// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: transaction.sql

package sqlc

import (
	"time"

	"github.com/google/uuid"
)

type CreateTransactionsParams struct {
	ExternalID    string     `json:"external_id"`
	Name          string     `json:"name"`
	Amount        int64      `json:"amount"`
	PaymentMethod string     `json:"payment_method"`
	Date          time.Time  `json:"date"`
	UserID        uuid.UUID  `json:"user_id"`
	AccountID     *uuid.UUID `json:"account_id"`
	InstitutionID *uuid.UUID `json:"institution_id"`
	CategoryID    *uuid.UUID `json:"category_id"`
}