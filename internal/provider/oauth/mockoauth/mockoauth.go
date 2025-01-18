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

	users map[string]*entity.User
}

const MockToken = "mock_token"

func NewMockOAuth(
	e *config.Env,
) *MockOAuth {
	if e.Environment == config.EnvironmentProduction {
		panic("mock oauth not allowed in production")
	}

	avatar := "https://avatars.githubusercontent.com/u/60039311"
	subscriptionExpiresAt := time.Now().AddDate(1, 0, 0)
	user := entity.User{
		ExternalID:            "6c2342aa-bdac-4efe-a31b-3a018072cff9",
		Name:                  "John Doe",
		Email:                 "johndoe@email.com",
		Avatar:                &avatar,
		Provider:              string(entity.ProviderMock),
		Tier:                  string(entity.TierPremium),
		SubscriptionExpiresAt: &subscriptionExpiresAt,
		VerifiedEmail:         true,
	}

	return &MockOAuth{
		e: e,
		users: map[string]*entity.User{
			MockToken: &user,
		},
	}
}

func (m *MockOAuth) GetUser(
	ctx context.Context,
	token string,
) (*entity.User, error) {
	return m.users[MockToken], nil
}

var _ oauth.Provider = &MockOAuth{}
