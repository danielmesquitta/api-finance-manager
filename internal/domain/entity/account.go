package entity

import (
	"time"

	"github.com/google/uuid"
)

type Account struct {
	ID         uuid.UUID   `json:"id,omitempty"`
	ExternalID string      `json:"external_id,omitempty"`
	Name       string      `json:"name,omitempty"`
	Type       AccountType `json:"type,omitempty"`
	CreatedAt  time.Time   `json:"created_at,omitempty"`
	UpdatedAt  time.Time   `json:"updated_at,omitempty"`

	UserID        string        `json:"user_id,omitempty"`
	User          *User         `json:"user,omitempty"`
	InstitutionID string        `json:"institution_id,omitempty"`
	Institution   *Institution  `json:"institution,omitempty"`
	Transactions  []Transaction `json:"transactions,omitempty"`
}
