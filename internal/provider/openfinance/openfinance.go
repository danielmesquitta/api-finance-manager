package openfinance

import (
	"context"
	"time"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
)

type InstitutionOptions struct {
	Types  []string `json:"types,omitzero"`
	Search string   `json:"search,omitzero"`
}

type InstitutionOption func(*InstitutionOptions)

func WithInstitutionTypes(types []string) InstitutionOption {
	return func(o *InstitutionOptions) {
		o.Types = types
	}
}

func WithInstitutionSearch(search string) InstitutionOption {
	return func(o *InstitutionOptions) {
		o.Search = search
	}
}

type TransactionOptions struct {
	StartDate time.Time `json:"startDate,omitzero"`
	EndDate   time.Time `json:"endDate,omitzero"`
}

type TransactionOption func(*TransactionOptions)

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

type Transaction struct {
	entity.Transaction
	CategoryExternalID      string
	PaymentMethodExternalID string
}

type Account struct {
	entity.Account
	Balance int64
}

type Client interface {
	ListInstitutions(
		ctx context.Context,
		options ...InstitutionOption,
	) ([]entity.Institution, error)
	ListTransactionCategories(
		ctx context.Context,
	) ([]entity.TransactionCategory, error)
	GetParentCategoryExternalID(
		externalCategoryID string,
		categoriesByExternalID map[string]entity.TransactionCategory,
	) (parentExternalID string, ok bool)
	ListTransactions(
		ctx context.Context,
		accountID string,
		options ...TransactionOption,
	) ([]Transaction, error)
	ListAccounts(
		ctx context.Context,
		connectionID string,
	) ([]Account, error)
}
