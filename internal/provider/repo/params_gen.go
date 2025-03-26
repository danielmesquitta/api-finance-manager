//nolint
//go:build !codeanalysis
// +build !codeanalysis

package repo

import (
	"time"

	"github.com/google/uuid"
)

type CreateAccountsParams struct {
	ID                uuid.UUID `json:"id"`
	ExternalID        string    `json:"external_id"`
	Name              string    `json:"name"`
	Type              string    `json:"type"`
	UserInstitutionID uuid.UUID `json:"user_institution_id"`
}

type CreateAccountBalancesParams struct {
	Amount    int64     `json:"amount"`
	AccountID uuid.UUID `json:"account_id"`
}

type GetUserBalanceOnDateParams struct {
	UserID uuid.UUID `json:"user_id"`
	Date   time.Time `json:"date"`
}

type ListAIChatMessagesAndAnswersParams struct {
	AiChatID uuid.UUID `json:"ai_chat_id"`
	Limit    int32     `json:"limit"`
	Offset   int32     `json:"offset"`
}

type UpdateAIChatParams struct {
	ID    uuid.UUID `json:"id"`
	Title *string   `json:"title"`
}

type CreateAIChatAnswerParams struct {
	Message         string    `json:"message"`
	AiChatMessageID uuid.UUID `json:"ai_chat_message_id"`
}

type UpdateAIChatAnswerParams struct {
	ID     uuid.UUID `json:"id"`
	Rating *string   `json:"rating"`
}

type CreateAIChatMessageParams struct {
	Message  string    `json:"message"`
	AiChatID uuid.UUID `json:"ai_chat_id"`
}

type CreateBudgetParams struct {
	Amount int64     `json:"amount"`
	Date   time.Time `json:"date"`
	UserID uuid.UUID `json:"user_id"`
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

type CreateBudgetCategoriesParams struct {
	Amount     int64     `json:"amount"`
	BudgetID   uuid.UUID `json:"budget_id"`
	CategoryID uuid.UUID `json:"category_id"`
}

type GetBudgetCategoryParams struct {
	UserID uuid.UUID `json:"user_id"`
	Date   time.Time `json:"date"`
}

type CreateFeedbackParams struct {
	Message string     `json:"message"`
	UserID  *uuid.UUID `json:"user_id"`
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

type CreateTransactionParams struct {
	Name            string    `json:"name"`
	Amount          int64     `json:"amount"`
	PaymentMethodID uuid.UUID `json:"payment_method_id"`
	Date            time.Time `json:"date"`
	UserID          uuid.UUID `json:"user_id"`
	CategoryID      uuid.UUID `json:"category_id"`
}

type CreateTransactionsParams struct {
	ExternalID      *string    `json:"external_id"`
	Name            string     `json:"name"`
	Amount          int64      `json:"amount"`
	PaymentMethodID uuid.UUID  `json:"payment_method_id"`
	Date            time.Time  `json:"date"`
	UserID          uuid.UUID  `json:"user_id"`
	AccountID       *uuid.UUID `json:"account_id"`
	InstitutionID   *uuid.UUID `json:"institution_id"`
	CategoryID      uuid.UUID  `json:"category_id"`
	IsIgnored       bool       `json:"is_ignored"`
}

type UpdateTransactionParams struct {
	ID              uuid.UUID  `json:"id"`
	Name            string     `json:"name"`
	Amount          int64      `json:"amount"`
	PaymentMethodID uuid.UUID  `json:"payment_method_id"`
	Date            time.Time  `json:"date"`
	AccountID       *uuid.UUID `json:"account_id"`
	InstitutionID   *uuid.UUID `json:"institution_id"`
	CategoryID      uuid.UUID  `json:"category_id"`
	UserID          uuid.UUID  `json:"user_id"`
}

type CreateTransactionCategoriesParams struct {
	ExternalID string `json:"external_id"`
	Name       string `json:"name"`
}

type CreateUserParams struct {
	Name                  string     `json:"name"`
	Email                 string     `json:"email"`
	Tier                  string     `json:"tier"`
	Avatar                *string    `json:"avatar"`
	SubscriptionExpiresAt *time.Time `json:"subscription_expires_at"`
}

type DeleteUserParams struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
}

type UpdateUserParams struct {
	ID                    uuid.UUID  `json:"id"`
	Name                  string     `json:"name"`
	Email                 string     `json:"email"`
	Tier                  string     `json:"tier"`
	Avatar                *string    `json:"avatar"`
	SubscriptionExpiresAt *time.Time `json:"subscription_expires_at"`
	SynchronizedAt        *time.Time `json:"synchronized_at"`
}

type UpdateUserSynchronizedAtParams struct {
	ID             uuid.UUID  `json:"id"`
	SynchronizedAt *time.Time `json:"synchronized_at"`
}

type CreateUserAuthProviderParams struct {
	ExternalID    string    `json:"external_id"`
	Provider      string    `json:"provider"`
	VerifiedEmail bool      `json:"verified_email"`
	UserID        uuid.UUID `json:"user_id"`
}

type GetUserAuthProviderParams struct {
	UserID   uuid.UUID `json:"user_id"`
	Provider string    `json:"provider"`
}

type UpdateUserAuthProviderParams struct {
	ID            uuid.UUID `json:"id"`
	VerifiedEmail bool      `json:"verified_email"`
}

type CreateUserInstitutionParams struct {
	ExternalID    string    `json:"external_id"`
	UserID        uuid.UUID `json:"user_id"`
	InstitutionID uuid.UUID `json:"institution_id"`
}
