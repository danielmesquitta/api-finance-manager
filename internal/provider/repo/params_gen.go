//nolint
//go:build !codeanalysis
// +build !codeanalysis

package repo

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type CreateBudgetParams struct {
	Amount int64     `json:"amount"`
	Date   time.Time `json:"date"`
	UserID uuid.UUID `json:"user_id"`
}

type CreateBudgetCategoriesParams struct {
	Amount     int64     `json:"amount"`
	BudgetID   uuid.UUID `json:"budget_id"`
	CategoryID uuid.UUID `json:"category_id"`
}

type GetBudgetParams struct {
	UserID uuid.UUID `json:"user_id"`
	Date   time.Time `json:"date"`
}

type UpdateBudgetParams struct {
	Amount int64     `json:"amount"`
	UserID uuid.UUID `json:"user_id"`
	Date   time.Time `json:"date"`
}

type CreateCategoriesParams struct {
	ExternalID string `json:"external_id"`
	Name       string `json:"name"`
}

type CreateInstitutionsParams struct {
	ExternalID string      `json:"external_id"`
	Name       string      `json:"name"`
	Logo       pgtype.Text `json:"logo"`
}

type CreateUserParams struct {
	ExternalID    string      `json:"external_id"`
	Provider      string      `json:"provider"`
	Name          string      `json:"name"`
	Email         string      `json:"email"`
	VerifiedEmail bool        `json:"verified_email"`
	Tier          string      `json:"tier"`
	Avatar        pgtype.Text `json:"avatar"`
}

type UpdateUserParams struct {
	ID                    uuid.UUID   `json:"id"`
	ExternalID            string      `json:"external_id"`
	Provider              string      `json:"provider"`
	Name                  string      `json:"name"`
	Email                 string      `json:"email"`
	VerifiedEmail         bool        `json:"verified_email"`
	Tier                  string      `json:"tier"`
	Avatar                pgtype.Text `json:"avatar"`
	SubscriptionExpiresAt *time.Time  `json:"subscription_expires_at"`
	SynchronizedAt        *time.Time  `json:"synchronized_at"`
}
