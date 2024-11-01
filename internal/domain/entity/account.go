package entity

import (
	"time"

	"github.com/google/uuid"
)

type Account struct {
	ID         uuid.UUID   `json:"id"`
	ExternalID string      `json:"external_id"`
	Name       string      `json:"name"`
	Type       AccountType `json:"type"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`

	UserID        string        `json:"user_id"`
	User          *User         `json:"user"`
	InstitutionID string        `json:"institution_id"`
	Institution   *Institution  `json:"institution"`
	Transactions  []Transaction `json:"transactions"`
}
