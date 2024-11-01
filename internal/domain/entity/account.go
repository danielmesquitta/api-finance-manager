package entity

import (
	"time"

	"github.com/google/uuid"
)

type Account struct {
	ID            uuid.UUID     `json:"id,omitempty"`
	ExternalID    string        `json:"external_id,omitempty"`
	Name          string        `json:"name,omitempty"`
	Balance       int64         `json:"balance,omitempty"`
	Type          AccountType   `json:"type,omitempty"`
	CreatedAt     time.Time     `json:"created_at,omitempty"`
	UpdatedAt     time.Time     `json:"updated_at,omitempty"`
	UserID        uuid.UUID     `json:"user_id,omitempty"`
	User          *User         `json:"user,omitempty"`
	InstitutionID uuid.UUID     `json:"institution_id,omitempty"`
	Institution   *Institution  `json:"institution,omitempty"`
	CreditCard    *CreditCard   `json:"credit_card,omitempty"`
	Transactions  []Transaction `json:"transactions,omitempty"`
}
