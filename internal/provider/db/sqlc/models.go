// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package sqlc

import (
	"time"

	"github.com/google/uuid"
)

type Account struct {
	ID                uuid.UUID  `json:"id"`
	ExternalID        string     `json:"external_id"`
	Name              string     `json:"name"`
	Type              string     `json:"type"`
	CreatedAt         time.Time  `json:"created_at"`
	DeletedAt         *time.Time `json:"deleted_at"`
	UserInstitutionID uuid.UUID  `json:"user_institution_id"`
}

type AccountBalance struct {
	ID        uuid.UUID  `json:"id"`
	Amount    int64      `json:"amount"`
	CreatedAt time.Time  `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	AccountID uuid.UUID  `json:"account_id"`
}

type AiChat struct {
	ID        uuid.UUID  `json:"id"`
	Title     *string    `json:"title"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	UserID    uuid.UUID  `json:"user_id"`
}

type AiChatAnswer struct {
	ID              uuid.UUID  `json:"id"`
	Message         string     `json:"message"`
	Rating          *string    `json:"rating"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	DeletedAt       *time.Time `json:"deleted_at"`
	AiChatMessageID uuid.UUID  `json:"ai_chat_message_id"`
}

type AiChatMessage struct {
	ID        uuid.UUID  `json:"id"`
	Message   string     `json:"message"`
	CreatedAt time.Time  `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	AiChatID  uuid.UUID  `json:"ai_chat_id"`
}

type Budget struct {
	ID        uuid.UUID  `json:"id"`
	Amount    int64      `json:"amount"`
	Date      time.Time  `json:"date"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	UserID    uuid.UUID  `json:"user_id"`
}

type BudgetCategory struct {
	ID         uuid.UUID  `json:"id"`
	Amount     int64      `json:"amount"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at"`
	BudgetID   uuid.UUID  `json:"budget_id"`
	CategoryID uuid.UUID  `json:"category_id"`
}

type Feedback struct {
	ID        uuid.UUID  `json:"id"`
	Message   string     `json:"message"`
	CreatedAt time.Time  `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	UserID    *uuid.UUID `json:"user_id"`
}

type Institution struct {
	ID         uuid.UUID  `json:"id"`
	ExternalID string     `json:"external_id"`
	Name       string     `json:"name"`
	Logo       *string    `json:"logo"`
	CreatedAt  time.Time  `json:"created_at"`
	DeletedAt  *time.Time `json:"deleted_at"`
}

type PaymentMethod struct {
	ID         uuid.UUID  `json:"id"`
	ExternalID string     `json:"external_id"`
	Name       string     `json:"name"`
	CreatedAt  time.Time  `json:"created_at"`
	DeletedAt  *time.Time `json:"deleted_at"`
}

type Transaction struct {
	ID              uuid.UUID  `json:"id"`
	ExternalID      *string    `json:"external_id"`
	Name            string     `json:"name"`
	Amount          int64      `json:"amount"`
	IsIgnored       bool       `json:"is_ignored"`
	Date            time.Time  `json:"date"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	DeletedAt       *time.Time `json:"deleted_at"`
	PaymentMethodID uuid.UUID  `json:"payment_method_id"`
	UserID          uuid.UUID  `json:"user_id"`
	CategoryID      uuid.UUID  `json:"category_id"`
	AccountID       *uuid.UUID `json:"account_id"`
	InstitutionID   *uuid.UUID `json:"institution_id"`
}

type TransactionCategory struct {
	ID         uuid.UUID  `json:"id"`
	ExternalID string     `json:"external_id"`
	Name       string     `json:"name"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at"`
}

type User struct {
	ID                    uuid.UUID  `json:"id"`
	Name                  string     `json:"name"`
	Email                 string     `json:"email"`
	Tier                  string     `json:"tier"`
	Avatar                *string    `json:"avatar"`
	SubscriptionExpiresAt *time.Time `json:"subscription_expires_at"`
	SynchronizedAt        *time.Time `json:"synchronized_at"`
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at"`
	DeletedAt             *time.Time `json:"deleted_at"`
}

type UserAuthProvider struct {
	ID            uuid.UUID  `json:"id"`
	ExternalID    string     `json:"external_id"`
	Provider      string     `json:"provider"`
	VerifiedEmail bool       `json:"verified_email"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at"`
	UserID        uuid.UUID  `json:"user_id"`
}

type UserInstitution struct {
	ID            uuid.UUID  `json:"id"`
	ExternalID    string     `json:"external_id"`
	CreatedAt     time.Time  `json:"created_at"`
	DeletedAt     *time.Time `json:"deleted_at"`
	UserID        uuid.UUID  `json:"user_id"`
	InstitutionID uuid.UUID  `json:"institution_id"`
}
