package mockoauth

import (
	"context"
	"time"

	"github.com/danielmesquitta/api-finance-manager/internal/config"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/oauth"
)

type MockOAuth struct {
	e *config.Env
}

const MockToken = "mock_token"

var Users = map[string]*entity.User{
	MockToken: func() *entity.User {
		avatar := "https://avatars.githubusercontent.com/u/60039311"
		subscriptionExpiresAt := time.Now().AddDate(0, 1, 0)
		return &entity.User{
			AuthID:                "6c2342aa-bdac-4efe-a31b-3a018072cff9",
			Name:                  "John Doe",
			Email:                 "johndoe@email.com",
			Avatar:                &avatar,
			Provider:              string(entity.ProviderMock),
			Tier:                  string(entity.TierPremium),
			SubscriptionExpiresAt: &subscriptionExpiresAt,
			VerifiedEmail:         true,
		}
	}(),
}

func NewMockOAuth(
	e *config.Env,
) *MockOAuth {
	if e.Environment == config.EnvironmentProduction {
		panic("mock oauth not allowed in production")
	}

	return &MockOAuth{
		e: e,
	}
}

func (m *MockOAuth) GetUser(
	ctx context.Context,
	token string,
) (*entity.User, error) {
	user, ok := Users[token]
	if !ok {
		return Users[MockToken], nil
	}

	return user, nil
}

var _ oauth.Provider = (*MockOAuth)(nil)
