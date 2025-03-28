package repo

import (
	"context"
	"time"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/google/uuid"
)

type TransactionOptions struct {
	Limit            uint        `json:"limit"`
	Offset           uint        `json:"offset"`
	Search           string      `json:"search"`
	StartDate        time.Time   `json:"start_date"`
	EndDate          time.Time   `json:"end_date"`
	CategoryIDs      []uuid.UUID `json:"category_ids"`
	InstitutionIDs   []uuid.UUID `json:"institution_ids"`
	PaymentMethodIDs []uuid.UUID `json:"payment_method_ids"`
	IsExpense        bool        `json:"is_expense"`
	IsIncome         bool        `json:"is_income"`
	IsIgnored        *bool       `json:"is_ignored"`
}

type TransactionRepo interface {
	CountTransactions(
		ctx context.Context,
		userID uuid.UUID,
		opts ...TransactionOptions,
	) (int64, error)
	CreateTransaction(
		ctx context.Context,
		params CreateTransactionParams,
	) error
	CreateTransactions(
		ctx context.Context,
		params []CreateTransactionsParams,
	) error
	GetTransactionByID(
		ctx context.Context,
		id uuid.UUID,
	) (*entity.FullTransaction, error)
	ListTransactions(
		ctx context.Context,
		userID uuid.UUID,
		opts ...TransactionOptions,
	) ([]entity.Transaction, error)
	ListFullTransactions(
		ctx context.Context,
		userID uuid.UUID,
		opts ...TransactionOptions,
	) ([]entity.FullTransaction, error)
	SumTransactions(
		ctx context.Context,
		userID uuid.UUID,
		opts ...TransactionOptions,
	) (int64, error)
	SumTransactionsByCategory(
		ctx context.Context,
		userID uuid.UUID,
		opts ...TransactionOptions,
	) (map[uuid.UUID]int64, error)
	UpdateTransaction(
		ctx context.Context,
		params UpdateTransactionParams,
	) error
}
