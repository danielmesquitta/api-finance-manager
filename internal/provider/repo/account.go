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

type AccountRepo interface {
	ListAccounts(
		ctx context.Context,
		opts ...AccountOptions,
	) ([]entity.Account, error)
	ListFullAccounts(
		ctx context.Context,
		opts ...AccountOptions,
	) ([]entity.FullAccount, error)
	CreateAccounts(
		ctx context.Context,
		params []CreateAccountsParams,
	) error
}
