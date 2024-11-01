package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID                    uuid.UUID `json:"id"`
	Name                  string    `json:"name"`
	Email                 string    `json:"email"`
	Tier                  Tier      `json:"tier"`
	SubscriptionExpiresAt time.Time `json:"subscription_expires_at"`
	SynchronizedAt        time.Time `json:"synchronized_at"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`

	Budget       *Budget       `json:"budget"`
	Transactions []Transaction `json:"transactions"`
	Accounts     []Account     `json:"accounts"`
}
