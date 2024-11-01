package entity

import (
	"time"

	"github.com/google/uuid"
)

type Institution struct {
	ID         uuid.UUID `json:"id,omitempty"`
	ExternalID string    `json:"external_id,omitempty"`
	Name       string    `json:"name,omitempty"`
	ImageURL   *string   `json:"image_url,omitempty"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
	UpdatedAt  time.Time `json:"updated_at,omitempty"`

	Accounts []Account `json:"accounts,omitempty"`
}
