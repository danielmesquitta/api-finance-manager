package repo

import (
	"context"
	"time"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/google/uuid"
)

type TransactionOptions struct {
	Limit           uint      `json:"limit"`
	Offset          uint      `json:"offset"`
	Search          string    `json:"search"`
	StartDate       time.Time `json:"start_date"`
	EndDate         time.Time `json:"end_date"`
	CategoryID      uuid.UUID `json:"category_id"`
	InstitutionID   uuid.UUID `json:"institution_id"`
	IsExpense       bool      `json:"is_expense"`
	IsIncome        bool      `json:"is_income"`
	PaymentMethodID uuid.UUID `json:"payment_method"`
	IsIgnored       *bool     `json:"is_ignored"`
}

type TransactionOption func(*TransactionOptions)

func WithTransactionPagination(
	limit uint,
	offset uint,
) TransactionOption {
	return func(o *TransactionOptions) {
		o.Limit = limit
		o.Offset = offset
	}
}

func WithTransactionSearch(search string) TransactionOption {
	return func(o *TransactionOptions) {
		o.Search = search
	}
}

func WithTransactionDateAfter(startDate time.Time) TransactionOption {
	return func(o *TransactionOptions) {
		o.StartDate = startDate
	}
}

func WithTransactionDateBefore(endDate time.Time) TransactionOption {
	return func(o *TransactionOptions) {
		o.EndDate = endDate
	}
}

func WithTransactionCategory(categoryID uuid.UUID) TransactionOption {
	return func(o *TransactionOptions) {
		o.CategoryID = categoryID
	}
}

func WithTransactionInstitution(
	institutionID uuid.UUID,
) TransactionOption {
	return func(o *TransactionOptions) {
		o.InstitutionID = institutionID
	}
}

func WithTransactionIsExpense(isExpense bool) TransactionOption {
	return func(o *TransactionOptions) {
		o.IsExpense = isExpense
	}
}

func WithTransactionIsIncome(isIncome bool) TransactionOption {
	return func(o *TransactionOptions) {
		o.IsIncome = isIncome
	}
}

func WithTransactionPaymentMethod(
	paymentMethodID uuid.UUID,
) TransactionOption {
	return func(o *TransactionOptions) {
		o.PaymentMethodID = paymentMethodID
	}
}

func WithTransactionIsIgnored(
	isIgnored bool,
) TransactionOption {
	return func(o *TransactionOptions) {
		o.IsIgnored = &isIgnored
	}
}

type TransactionRepo interface {
	CountTransactions(
		ctx context.Context,
		userID uuid.UUID,
		opts ...TransactionOption,
	) (int64, error)
	CreateTransactions(
		ctx context.Context,
		params []CreateTransactionsParams,
	) error
	GetTransaction(
		ctx context.Context,
		params GetTransactionParams,
	) (*entity.FullTransaction, error)
	ListTransactions(
		ctx context.Context,
		userID uuid.UUID,
		opts ...TransactionOption,
	) ([]entity.Transaction, error)
	ListTransactionsWithCategoriesAndInstitutions(
		ctx context.Context,
		userID uuid.UUID,
		opts ...TransactionOption,
	) ([]entity.FullTransaction, error)
	SumTransactions(
		ctx context.Context,
		userID uuid.UUID,
		opts ...TransactionOption,
	) (int64, error)
	SumTransactionsByCategory(
		ctx context.Context,
		userID uuid.UUID,
		opts ...TransactionOption,
	) (map[uuid.UUID]int64, error)
	UpdateTransaction(
		ctx context.Context,
		params UpdateTransactionParams,
	) error
}
