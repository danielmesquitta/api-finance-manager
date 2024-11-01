package entity

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	ID            uuid.UUID     `json:"id,omitempty"`
	ExternalID    string        `json:"external_id,omitempty"`
	Name          string        `json:"name,omitempty"`
	Description   *string       `json:"description,omitempty"`
	Amount        int64         `json:"amount,omitempty"`
	PaymentMethod PaymentMethod `json:"payment_method,omitempty"`
	IsIgnored     bool          `json:"is_ignored,omitempty"`
	Date          time.Time     `json:"date,omitempty"`
	CreatedAt     time.Time     `json:"created_at,omitempty"`
	UpdatedAt     time.Time     `json:"updated_at,omitempty"`
	UserID        uuid.UUID     `json:"user_id,omitempty"`
	User          *User         `json:"user,omitempty"`
	AccountID     *uuid.UUID    `json:"account_id,omitempty"`
	Account       *Account      `json:"account,omitempty"`
	CategoryID    *uuid.UUID    `json:"category_id,omitempty"`
	Category      *Category     `json:"category,omitempty"`
}
