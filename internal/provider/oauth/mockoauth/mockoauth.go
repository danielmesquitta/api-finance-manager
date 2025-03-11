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

const DefaultMockToken = "mock_token"
const UnregisteredUserMockToken = "unregistered_user_mock_token"

var Users = map[string]*entity.User{
	DefaultMockToken: func() *entity.User {
		avatar := "https://avatar.iran.liara.run/public/15"
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
	UnregisteredUserMockToken: func() *entity.User {
		avatar := "https://avatar.iran.liara.run/public/82"
		return &entity.User{
			AuthID:        "016aecbd-fae5-4ff0-9046-03b7eabf6a5c",
			Name:          "Jane Doe",
			Email:         "janedoe@email.com",
			Avatar:        &avatar,
			Provider:      string(entity.ProviderMock),
			Tier:          string(entity.TierFree),
			VerifiedEmail: true,
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
		return Users[DefaultMockToken], nil
	}

	return user, nil
}

var _ oauth.Provider = (*MockOAuth)(nil)
