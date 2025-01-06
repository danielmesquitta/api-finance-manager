package openfinance

import (
	"context"
	"time"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
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
	GetParentCategoryExternalID(
		externalCategoryID string,
		categoriesByExternalID map[string]entity.Category,
	) (parentExternalID string, ok bool)
	ListTransactions(
		ctx context.Context,
		accountID string,
		options ...ListTransactionsOption,
	) ([]Transaction, error)
	ListAccounts(
		ctx context.Context,
		connectionID string,
	) ([]entity.Account, error)
}
