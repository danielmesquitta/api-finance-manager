package entity

import (
	"time"

	"github.com/google/uuid"
)

type Institution struct {
	ID         uuid.UUID `json:"id"`
	ExternalID string    `json:"external_id"`
	Name       string    `json:"name"`
	ImageURL   *string   `json:"image_url"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`

	Accounts []Account `json:"accounts"`
}
