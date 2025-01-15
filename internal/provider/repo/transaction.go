package repo

import (
	"context"
	"time"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/google/uuid"
)

type ListTransactionsOptions struct {
	Limit           uint      `json:"-"`
	Offset          uint      `json:"-"`
	Search          string    `json:"search"`
	StartDate       time.Time `json:"start_date"`
	EndDate         time.Time `json:"end_date"`
	CategoryID      uuid.UUID `json:"category_id"`
	InstitutionID   uuid.UUID `json:"institution_id"`
	IsExpense       bool      `json:"is_expense"`
	IsIncome        bool      `json:"is_income"`
	PaymentMethodID uuid.UUID `json:"payment_method"`
}

type ListTransactionsOption func(*ListTransactionsOptions)

func WithTransactionsPagination(
	limit uint,
	offset uint,
) ListTransactionsOption {
	return func(o *ListTransactionsOptions) {
		o.Limit = limit
		o.Offset = offset
	}
}

func WithTransactionsSearch(search string) ListTransactionsOption {
	return func(o *ListTransactionsOptions) {
		o.Search = search
	}
}

func WithTransactionDateAfter(startDate time.Time) ListTransactionsOption {
	return func(o *ListTransactionsOptions) {
		o.StartDate = startDate
	}
}

func WithTransactionDateBefore(endDate time.Time) ListTransactionsOption {
	return func(o *ListTransactionsOptions) {
		o.EndDate = endDate
	}
}

func WithTransactionCategory(categoryID uuid.UUID) ListTransactionsOption {
	return func(o *ListTransactionsOptions) {
		o.CategoryID = categoryID
	}
}

func WithTransactionInstitution(
	institutionID uuid.UUID,
) ListTransactionsOption {
	return func(o *ListTransactionsOptions) {
		o.InstitutionID = institutionID
	}
}

func WithTransactionIsExpense(isExpense bool) ListTransactionsOption {
	return func(o *ListTransactionsOptions) {
		o.IsExpense = isExpense
	}
}

func WithTransactionIsIncome(isIncome bool) ListTransactionsOption {
	return func(o *ListTransactionsOptions) {
		o.IsIncome = isIncome
	}
}

func WithTransactionPaymentMethod(
	paymentMethodID uuid.UUID,
) ListTransactionsOption {
	return func(o *ListTransactionsOptions) {
		o.PaymentMethodID = paymentMethodID
	}
}

type TransactionRepo interface {
	ListTransactions(
		ctx context.Context,
		userID uuid.UUID,
		opts ...ListTransactionsOption,
	) ([]entity.Transaction, error)
	ListTransactionsWithCategoriesAndInstitutions(
		ctx context.Context,
		userID uuid.UUID,
		opts ...ListTransactionsOption,
	) ([]entity.TransactionWithCategoryAndInstitution, error)
	CountTransactions(
		ctx context.Context,
		userID uuid.UUID,
		opts ...ListTransactionsOption,
	) (int64, error)
	CreateTransactions(
		ctx context.Context,
		params []CreateTransactionsParams,
	) error
}
