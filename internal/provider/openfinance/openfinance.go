package openfinance

import (
	"context"
	"time"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/google/uuid"
)

type ListInstitutionsOptions struct {
	Types  []string `json:"types,omitempty"`
	Search string   `json:"search,omitempty"`
}

type ListInstitutionsOption func(*ListInstitutionsOptions)

func WithInstitutionTypes(types []string) ListInstitutionsOption {
	return func(o *ListInstitutionsOptions) {
		o.Types = types
	}
}

func WithInstitutionSearch(search string) ListInstitutionsOption {
	return func(o *ListInstitutionsOptions) {
		o.Search = search
	}
}

type ListTransactionsOptions struct {
	StartDate time.Time `json:"startDate,omitempty"`
	EndDate   time.Time `json:"endDate,omitempty"`
}

type ListTransactionsOption func(*ListTransactionsOptions)

func WithTransactionStartDate(startDate time.Time) ListTransactionsOption {
	return func(o *ListTransactionsOptions) {
		o.StartDate = startDate
	}
}

func WithTransactionEndDate(endDate time.Time) ListTransactionsOption {
	return func(o *ListTransactionsOptions) {
		o.EndDate = endDate
	}
}

type Transaction struct {
	entity.Transaction
	CategoryExternalID string
}
type Client interface {
	ListInstitutions(
		ctx context.Context,
		options ...ListInstitutionsOption,
	) ([]entity.Institution, error)
	ListCategories(
		ctx context.Context,
	) ([]entity.Category, error)
	ListTransactions(
		ctx context.Context,
		accountID uuid.UUID,
		options ...ListTransactionsOption,
	) ([]Transaction, error)
}
