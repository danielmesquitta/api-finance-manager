//nolint
//go:build !codeanalysis
// +build !codeanalysis

package repo

import (
	"time"

	"github.com/google/uuid"
)

type CreateAccountsParams struct {
	ExternalID    string    `json:"external_id"`
	Name          string    `json:"name"`
	Type          string    `json:"type"`
	UserID        uuid.UUID `json:"user_id"`
	InstitutionID uuid.UUID `json:"institution_id"`
}

type CreateAccountBalancesParams struct {
	Amount    int64     `json:"amount"`
	UserID    uuid.UUID `json:"user_id"`
	AccountID uuid.UUID `json:"account_id"`
}

type GetUserBalanceOnDateParams struct {
	UserID uuid.UUID `json:"user_id"`
	Date   time.Time `json:"date"`
}

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

type GetBudgetCategoryParams struct {
	UserID uuid.UUID `json:"user_id"`
	Date   time.Time `json:"date"`
}

type UpdateBudgetParams struct {
	Amount int64     `json:"amount"`
	UserID uuid.UUID `json:"user_id"`
	Date   time.Time `json:"date"`
}

type CreateInstitutionsParams struct {
	ExternalID string  `json:"external_id"`
	Name       string  `json:"name"`
	Logo       *string `json:"logo"`
}

type CreatePaymentMethodsParams struct {
	ExternalID string `json:"external_id"`
	Name       string `json:"name"`
}

type CreateTransactionsParams struct {
	ExternalID      string     `json:"external_id"`
	Name            string     `json:"name"`
	Amount          int64      `json:"amount"`
	PaymentMethodID uuid.UUID  `json:"payment_method_id"`
	Date            time.Time  `json:"date"`
	UserID          uuid.UUID  `json:"user_id"`
	AccountID       *uuid.UUID `json:"account_id"`
	InstitutionID   *uuid.UUID `json:"institution_id"`
	CategoryID      *uuid.UUID `json:"category_id"`
}

type GetTransactionParams struct {
	ID     uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"user_id"`
}

type UpdateTransactionParams struct {
	ID              uuid.UUID  `json:"id"`
	Name            string     `json:"name"`
	Amount          int64      `json:"amount"`
	PaymentMethodID uuid.UUID  `json:"payment_method_id"`
	Date            time.Time  `json:"date"`
	AccountID       *uuid.UUID `json:"account_id"`
	InstitutionID   *uuid.UUID `json:"institution_id"`
	CategoryID      *uuid.UUID `json:"category_id"`
	UserID          uuid.UUID  `json:"user_id"`
}

type CreateTransactionCategoriesParams struct {
	ExternalID string `json:"external_id"`
	Name       string `json:"name"`
}

type CreateUserParams struct {
	AuthID                string     `json:"auth_id"`
	OpenFinanceID         *string    `json:"open_finance_id"`
	Provider              string     `json:"provider"`
	Name                  string     `json:"name"`
	Email                 string     `json:"email"`
	VerifiedEmail         bool       `json:"verified_email"`
	Tier                  string     `json:"tier"`
	Avatar                *string    `json:"avatar"`
	SubscriptionExpiresAt *time.Time `json:"subscription_expires_at"`
}

type UpdateUserParams struct {
	ID                    uuid.UUID  `json:"id"`
	AuthID                string     `json:"auth_id"`
	OpenFinanceID         *string    `json:"open_finance_id"`
	Provider              string     `json:"provider"`
	Name                  string     `json:"name"`
	Email                 string     `json:"email"`
	VerifiedEmail         bool       `json:"verified_email"`
	Tier                  string     `json:"tier"`
	Avatar                *string    `json:"avatar"`
	SubscriptionExpiresAt *time.Time `json:"subscription_expires_at"`
	SynchronizedAt        *time.Time `json:"synchronized_at"`
}

type UpdateUserSynchronizedAtParams struct {
	ID             uuid.UUID  `json:"id"`
	SynchronizedAt *time.Time `json:"synchronized_at"`
}
