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
	Amount float64   `json:"amount"`
	UserID uuid.UUID `json:"user_id"`
}

type CreateBudgetCategoriesParams struct {
	Amount     float64   `json:"amount"`
	BudgetID   uuid.UUID `json:"budget_id"`
	CategoryID uuid.UUID `json:"category_id"`
}

type UpdateBudgetParams struct {
	Amount float64   `json:"amount"`
	UserID uuid.UUID `json:"user_id"`
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
