package repo

import (
	"context"
	"time"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/google/uuid"
)

type CreateUserParams struct {
	Name                  string      `json:"name,omitempty"`
	Email                 string      `json:"email,omitempty"`
	Tier                  entity.Tier `json:"tier,omitempty"`
	Avatar                *string     `json:"avatar,omitempty"`
	SubscriptionExpiresAt time.Time   `json:"subscription_expires_at,omitempty"`
}

type UpdateUserParams struct {
	ID                    uuid.UUID   `json:"id,omitempty"`
	Name                  string      `json:"name,omitempty"`
	Email                 string      `json:"email,omitempty"`
	Tier                  entity.Tier `json:"tier,omitempty"`
	Avatar                *string     `json:"avatar,omitempty"`
	SubscriptionExpiresAt time.Time   `json:"subscription_expires_at,omitempty"`
	SynchronizedAt        time.Time   `json:"synchronized_at,omitempty"`
}

type UserRepo interface {
	CreateUser(
		ctx context.Context,
		params CreateUserParams,
	) (*entity.User, error)
	UpdateUser(
		ctx context.Context,
		params UpdateUserParams,
	) (*entity.User, error)
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
}
