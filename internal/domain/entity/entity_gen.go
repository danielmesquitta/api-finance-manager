//nolint
//go:build !codeanalysis
// +build !codeanalysis

package entity

import (
	"github.com/google/uuid"
	"time"
)

type Tier string

const (
	TierTrial Tier = "TRIAL"
	TierPro   Tier = "PRO"
)

type Provider string

const (
	ProviderGoogle Provider = "GOOGLE"
	ProviderApple  Provider = "APPLE"
)

type User struct {
	ID                    uuid.UUID     `json:"id,omitempty"`
	ExternalId            string        `json:"external_id,omitempty"`
	Provider              Provider      `json:"provider,omitempty"`
	Name                  string        `json:"name,omitempty"`
	Email                 string        `json:"email,omitempty"`
	VerifiedEmail         bool          `json:"verified_email,omitempty"`
	Tier                  Tier          `json:"tier,omitempty"`
	Avatar                *string       `json:"avatar,omitempty"`
	SubscriptionExpiresAt time.Time     `json:"subscription_expires_at,omitempty"`
	SynchronizedAt        time.Time     `json:"synchronized_at,omitempty"`
	CreatedAt             time.Time     `json:"created_at,omitempty"`
	UpdatedAt             time.Time     `json:"updated_at,omitempty"`
	Budget                *Budget       `json:"budget,omitempty"`
	Transactions          []Transaction `json:"transactions,omitempty"`
	Investments           []Investment  `json:"investments,omitempty"`
	Accounts              []Account     `json:"accounts,omitempty"`
}

type PaymentMethod string

const (
	PaymentMethodPix          PaymentMethod = "PIX"
	PaymentMethodBoleto       PaymentMethod = "BOLETO"
	PaymentMethodCreditCard   PaymentMethod = "CREDIT_CARD"
	PaymentMethodDebitCard    PaymentMethod = "DEBIT_CARD"
	PaymentMethodTransference PaymentMethod = "TRANSFERENCE"
	PaymentMethodUnknown      PaymentMethod = "UNKNOWN"
)

type Transaction struct {
	ID            uuid.UUID     `json:"id,omitempty"`
	ExternalId    string        `json:"external_id,omitempty"`
	Name          string        `json:"name,omitempty"`
	Description   *string       `json:"description,omitempty"`
	Amount        int64         `json:"amount,omitempty"`
	PaymentMethod PaymentMethod `json:"payment_method,omitempty"`
	IsIgnored     bool          `json:"is_ignored,omitempty"`
	Date          time.Time     `json:"date,omitempty"`
	CreatedAt     time.Time     `json:"created_at,omitempty"`
	UpdatedAt     time.Time     `json:"updated_at,omitempty"`
	User          User          `json:"user,omitempty"`
	UserId        uuid.UUID     `json:"user_id,omitempty"`
	Account       *Account      `json:"account,omitempty"`
	AccountId     *uuid.UUID    `json:"account_id,omitempty"`
	Category      *Category     `json:"category,omitempty"`
	CategoryId    *uuid.UUID    `json:"category_id,omitempty"`
}

type InvestmentType string

const (
	InvestmentTypeFixedIncome InvestmentType = "FIXED_INCOME"
	InvestmentTypeUnknown     InvestmentType = "UNKNOWN"
)

type InvestmentRateType string

const (
	InvestmentRateTypeCdi      InvestmentRateType = "CDI"
	InvestmentRateTypeSelic    InvestmentRateType = "SELIC"
	InvestmentRateTypeIpca     InvestmentRateType = "IPCA"
	InvestmentRateTypePrefixed InvestmentRateType = "PREFIXED"
	InvestmentRateTypeUnknown  InvestmentRateType = "UNKNOWN"
)

type Investment struct {
	ID         uuid.UUID          `json:"id,omitempty"`
	ExternalId string             `json:"external_id,omitempty"`
	Name       string             `json:"name,omitempty"`
	Amount     int64              `json:"amount,omitempty"`
	Type       InvestmentType     `json:"type,omitempty"`
	Rate       int64              `json:"rate,omitempty"`
	RateType   InvestmentRateType `json:"rateType,omitempty"`
	CreatedAt  time.Time          `json:"created_at,omitempty"`
	UpdatedAt  time.Time          `json:"updated_at,omitempty"`
	User       User               `json:"user,omitempty"`
	UserId     uuid.UUID          `json:"user_id,omitempty"`
}

type Category struct {
	ID               uuid.UUID        `json:"id,omitempty"`
	ExternalId       string           `json:"external_id,omitempty"`
	Name             string           `json:"name,omitempty"`
	CreatedAt        time.Time        `json:"created_at,omitempty"`
	UpdatedAt        time.Time        `json:"updated_at,omitempty"`
	Transactions     []Transaction    `json:"transactions,omitempty"`
	BudgetCategories []BudgetCategory `json:"budget_categories,omitempty"`
}

type AccountType string

const (
	AccountTypeBank    AccountType = "BANK"
	AccountTypeCredit  AccountType = "CREDIT"
	AccountTypeUnknown AccountType = "UNKNOWN"
)

type Account struct {
	ID            uuid.UUID     `json:"id,omitempty"`
	ExternalId    string        `json:"external_id,omitempty"`
	Name          string        `json:"name,omitempty"`
	Balance       int64         `json:"balance,omitempty"`
	Type          AccountType   `json:"type,omitempty"`
	CreatedAt     time.Time     `json:"created_at,omitempty"`
	UpdatedAt     time.Time     `json:"updated_at,omitempty"`
	User          User          `json:"user,omitempty"`
	UserId        uuid.UUID     `json:"user_id,omitempty"`
	Institution   Institution   `json:"institution,omitempty"`
	InstitutionId uuid.UUID     `json:"institution_id,omitempty"`
	CreditCard    *CreditCard   `json:"CreditCard,omitempty"`
	Transactions  []Transaction `json:"transactions,omitempty"`
}

type CreditCard struct {
	ID             uuid.UUID `json:"id,omitempty"`
	Level          string    `json:"level,omitempty"`
	Brand          string    `json:"brand,omitempty"`
	Limit          int64     `json:"limit,omitempty"`
	AvailableLimit int64     `json:"available_limit,omitempty"`
	CreatedAt      time.Time `json:"created_at,omitempty"`
	UpdatedAt      time.Time `json:"updated_at,omitempty"`
	Account        Account   `json:"account,omitempty"`
	AccountId      uuid.UUID `json:"account_id,omitempty"`
}

type Institution struct {
	ID         uuid.UUID `json:"id,omitempty"`
	ExternalId string    `json:"external_id,omitempty"`
	Name       string    `json:"name,omitempty"`
	ImageUrl   *string   `json:"image_url,omitempty"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
	UpdatedAt  time.Time `json:"updated_at,omitempty"`
	Accounts   []Account `json:"accounts,omitempty"`
}

type Budget struct {
	ID               uuid.UUID        `json:"id,omitempty"`
	Amount           int64            `json:"amount,omitempty"`
	CreatedAt        time.Time        `json:"created_at,omitempty"`
	UpdatedAt        time.Time        `json:"updated_at,omitempty"`
	User             User             `json:"user,omitempty"`
	UserId           uuid.UUID        `json:"user_id,omitempty"`
	BudgetCategories []BudgetCategory `json:"budget_categories,omitempty"`
}

type BudgetCategory struct {
	ID         uuid.UUID `json:"id,omitempty"`
	Amount     int64     `json:"amount,omitempty"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
	UpdatedAt  time.Time `json:"updated_at,omitempty"`
	Budget     Budget    `json:"budget,omitempty"`
	BudgetId   uuid.UUID `json:"budget_id,omitempty"`
	Category   Category  `json:"category,omitempty"`
	CategoryId uuid.UUID `json:"category_id,omitempty"`
}
