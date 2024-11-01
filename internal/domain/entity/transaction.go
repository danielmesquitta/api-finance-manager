package entity

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	ID            uuid.UUID     `json:"id"`
	ExternalID    string        `json:"external_id"`
	Name          string        `json:"name"`
	Description   *string       `json:"description"`
	Amount        int64         `json:"amount"`
	PaymentMethod PaymentMethod `json:"payment_method"`
	IsIgnored     bool          `json:"is_ignored"`
	Date          time.Time     `json:"date"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`

	UserID     string    `json:"user_id"`
	User       *User     `json:"user"`
	AccountID  *string   `json:"account_id"`
	Account    *Account  `json:"account"`
	CategoryID *string   `json:"category_id"`
	Category   *Category `json:"category"`
}
