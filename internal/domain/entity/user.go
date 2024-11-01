package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID                    uuid.UUID `json:"id,omitempty"`
	Name                  string    `json:"name,omitempty"`
	Email                 string    `json:"email,omitempty"`
	Tier                  Tier      `json:"tier,omitempty"`
	Avatar                *string   `json:"avatar,omitempty"`
	SubscriptionExpiresAt time.Time `json:"subscription_expires_at,omitempty"`
	SynchronizedAt        time.Time `json:"synchronized_at,omitempty"`
	CreatedAt             time.Time `json:"created_at,omitempty"`
	UpdatedAt             time.Time `json:"updated_at,omitempty"`

	Budget       *Budget       `json:"budget,omitempty"`
	Transactions []Transaction `json:"transactions,omitempty"`
	Accounts     []Account     `json:"accounts,omitempty"`
}
