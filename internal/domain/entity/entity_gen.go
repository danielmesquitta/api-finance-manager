//nolint
//go:build !codeanalysis
// +build !codeanalysis

package entity

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID                    uuid.UUID  `json:"id,omitempty"`
	ExternalID            string     `json:"external_id,omitempty"`
	Provider              string     `json:"provider,omitempty"`
	Name                  string     `json:"name,omitempty"`
	Email                 string     `json:"email,omitempty"`
	VerifiedEmail         bool       `json:"verified_email,omitempty"`
	Tier                  string     `json:"tier,omitempty"`
	Avatar                *string    `json:"avatar,omitempty"`
	SubscriptionExpiresAt *time.Time `json:"subscription_expires_at,omitempty"`
	SynchronizedAt        *time.Time `json:"synchronized_at,omitempty"`
	CreatedAt             time.Time  `json:"created_at,omitempty"`
	UpdatedAt             time.Time  `json:"updated_at,omitempty"`
}

type Transaction struct {
	ID            uuid.UUID  `json:"id,omitempty"`
	ExternalID    string     `json:"external_id,omitempty"`
	Name          string     `json:"name,omitempty"`
	Description   *string    `json:"description,omitempty"`
	Amount        int64      `json:"amount,omitempty"`
	PaymentMethod string     `json:"payment_method,omitempty"`
	IsIgnored     bool       `json:"is_ignored,omitempty"`
	Date          time.Time  `json:"date,omitempty"`
	CreatedAt     time.Time  `json:"created_at,omitempty"`
	UpdatedAt     time.Time  `json:"updated_at,omitempty"`
	UserID        uuid.UUID  `json:"user_id,omitempty"`
	AccountID     *uuid.UUID `json:"account_id,omitempty"`
	CategoryID    *uuid.UUID `json:"category_id,omitempty"`
}

type Investment struct {
	ID         uuid.UUID `json:"id,omitempty"`
	ExternalID string    `json:"external_id,omitempty"`
	Name       string    `json:"name,omitempty"`
	Amount     int64     `json:"amount,omitempty"`
	Type       string    `json:"type,omitempty"`
	Rate       int64     `json:"rate,omitempty"`
	RateType   string    `json:"rateType,omitempty"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
	UpdatedAt  time.Time `json:"updated_at,omitempty"`
	UserID     uuid.UUID `json:"user_id,omitempty"`
}

type Category struct {
	ID         uuid.UUID `json:"id,omitempty"`
	ExternalID string    `json:"external_id,omitempty"`
	Name       string    `json:"name,omitempty"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
	UpdatedAt  time.Time `json:"updated_at,omitempty"`
}

type Account struct {
	ID            uuid.UUID `json:"id,omitempty"`
	ExternalID    string    `json:"external_id,omitempty"`
	Name          string    `json:"name,omitempty"`
	Balance       int64     `json:"balance,omitempty"`
	Type          string    `json:"type,omitempty"`
	CreatedAt     time.Time `json:"created_at,omitempty"`
	UpdatedAt     time.Time `json:"updated_at,omitempty"`
	UserID        uuid.UUID `json:"user_id,omitempty"`
	InstitutionID uuid.UUID `json:"institution_id,omitempty"`
}

type CreditCard struct {
	ID             uuid.UUID `json:"id,omitempty"`
	Level          string    `json:"level,omitempty"`
	Brand          string    `json:"brand,omitempty"`
	Limit          int64     `json:"limit,omitempty"`
	AvailableLimit int64     `json:"available_limit,omitempty"`
	CreatedAt      time.Time `json:"created_at,omitempty"`
	UpdatedAt      time.Time `json:"updated_at,omitempty"`
	AccountID      uuid.UUID `json:"account_id,omitempty"`
}

type Institution struct {
	ID         uuid.UUID `json:"id,omitempty"`
	ExternalID string    `json:"external_id,omitempty"`
	Name       string    `json:"name,omitempty"`
	Logo       *string   `json:"logo,omitempty"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
	UpdatedAt  time.Time `json:"updated_at,omitempty"`
}

type Budget struct {
	ID        uuid.UUID `json:"id,omitempty"`
	Amount    int64     `json:"amount,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	UserID    uuid.UUID `json:"user_id,omitempty"`
}

type BudgetCategory struct {
	ID         uuid.UUID `json:"id,omitempty"`
	Amount     int64     `json:"amount,omitempty"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
	UpdatedAt  time.Time `json:"updated_at,omitempty"`
	BudgetID   uuid.UUID `json:"budget_id,omitempty"`
	CategoryID uuid.UUID `json:"category_id,omitempty"`
}
