package repo

import (
	"context"
	"time"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/google/uuid"
)

type CreateUserParams struct {
	Name                  string      `json:"name"`
	Email                 string      `json:"email"`
	Tier                  entity.Tier `json:"tier"`
	SubscriptionExpiresAt time.Time   `json:"subscription_expires_at"`
}

type UserRepo interface {
	CreateUser(ctx context.Context, params CreateUserParams) (*entity.User, error)
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
}
