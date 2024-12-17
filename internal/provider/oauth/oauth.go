package oauth

import (
	"context"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
)

type Provider interface {
	GetUser(ctx context.Context, token string) (*entity.User, error)
}
