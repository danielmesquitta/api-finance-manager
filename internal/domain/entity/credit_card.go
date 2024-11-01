package entity

import (
	"time"

	"github.com/google/uuid"
)

type CreditCard struct {
	ID             uuid.UUID `json:"id,omitempty"`
	Level          string    `json:"level,omitempty"`
	Brand          string    `json:"brand,omitempty"`
	Limit          int64     `json:"limit,omitempty"`
	AvailableLimit int64     `json:"available_limit,omitempty"`
	CreatedAt      time.Time `json:"created_at,omitempty"`
	UpdatedAt      time.Time `json:"updated_at,omitempty"`
	AccountID      uuid.UUID `json:"account_id,omitempty"`
	Account        *Account  `json:"account,omitempty"`
}
