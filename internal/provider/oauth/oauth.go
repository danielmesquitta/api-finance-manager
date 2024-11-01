package oauth

import "github.com/danielmesquitta/api-finance-manager/internal/domain/entity"

type Provider interface {
	GetUser(token string) (*entity.User, error)
}
