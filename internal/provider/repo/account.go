package repo

import (
	"context"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/google/uuid"
)

type AccountOptions struct {
	Limit                uint                 `json:"limit"`
	Offset               uint                 `json:"offset"`
	UserIDs              []uuid.UUID          `json:"user_id"`
	ExternalIDs          []string             `json:"external_ids"`
	UserTiers            []entity.Tier        `json:"user_tiers"`
	Types                []entity.AccountType `json:"types"`
	IsSubscriptionActive *bool                `json:"is_subscription_active"`
}

type AccountOption func(*AccountOptions)

func WithAccountPagination(
	limit uint,
	offset uint,
) AccountOption {
	return func(o *AccountOptions) {
		o.Limit = limit
		o.Offset = offset
	}
}

func WithAccountUserIDs(
	userIDs []uuid.UUID,
) AccountOption {
	return func(o *AccountOptions) {
		o.UserIDs = userIDs
	}
}

func WithAccountExternalIDs(
	externalIDs []string,
) AccountOption {
	return func(o *AccountOptions) {
		o.ExternalIDs = externalIDs
	}
}

func WithAccountUserTiers(
	userTiers []entity.Tier,
) AccountOption {
	return func(o *AccountOptions) {
		o.UserTiers = userTiers
	}
}

func WithAccountSubscriptionActive(
	isSubscriptionActive bool,
) AccountOption {
	return func(o *AccountOptions) {
		o.IsSubscriptionActive = &isSubscriptionActive
	}
}

func WithAccountTypes(
	types []entity.AccountType,
) AccountOption {
	return func(o *AccountOptions) {
		o.Types = types
	}
}

type AccountRepo interface {
	ListAccounts(
		ctx context.Context,
		opts ...AccountOption,
	) ([]entity.Account, error)
	ListFullAccounts(
		ctx context.Context,
		opts ...AccountOption,
	) ([]entity.FullAccount, error)
	CreateAccounts(
		ctx context.Context,
		params []CreateAccountsParams,
	) error
}
