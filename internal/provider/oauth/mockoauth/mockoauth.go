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

const DefaultMockToken = "default_mock_token"
const FreeTierMockToken = "free_tier_mock_token"
const UnregisteredUserMockToken = "unregistered_user_mock_token"

type User struct {
	*entity.User
	AuthProvider *entity.UserAuthProvider
}

var Users = map[string]*User{
	DefaultMockToken: func() *User {
		avatar := "https://avatar.iran.liara.run/public/15"
		subscriptionExpiresAt := time.Now().AddDate(0, 1, 0)
		return &User{
			User: &entity.User{
				Name:                  "John Doe",
				Email:                 "johndoe@email.com",
				Avatar:                &avatar,
				Tier:                  string(entity.TierPremium),
				SubscriptionExpiresAt: &subscriptionExpiresAt,
			},
			AuthProvider: &entity.UserAuthProvider{
				ExternalID:    "6c2342aa-bdac-4efe-a31b-3a018072cff9",
				Provider:      entity.ProviderMock,
				VerifiedEmail: true,
			},
		}
	}(),
	FreeTierMockToken: func() *User {
		avatar := "https://avatar.iran.liara.run/public/82"
		return &User{
			User: &entity.User{
				Name:   "Jane Doe",
				Email:  "janedoe@email.com",
				Avatar: &avatar,
				Tier:   string(entity.TierFree),
			},
			AuthProvider: &entity.UserAuthProvider{
				ExternalID:    "016aecbd-fae5-4ff0-9046-03b7eabf6a5c",
				Provider:      entity.ProviderMock,
				VerifiedEmail: true,
			},
		}
	}(),
	UnregisteredUserMockToken: func() *User {
		avatar := "https://avatar.iran.liara.run/public/13"
		return &User{
			User: &entity.User{
				Name:   "Joseph Doe",
				Email:  "josephdoe@email.com",
				Avatar: &avatar,
				Tier:   string(entity.TierFree),
			},
			AuthProvider: &entity.UserAuthProvider{
				ExternalID:    "2824923b-2d93-4473-8397-32680bb412b4",
				Provider:      entity.ProviderMock,
				VerifiedEmail: true,
			},
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
) (*entity.User, *entity.UserAuthProvider, error) {
	user, ok := Users[token]
	if !ok {
		user = Users[DefaultMockToken]
	}

	return user.User, user.AuthProvider, nil
}

var _ oauth.Provider = (*MockOAuth)(nil)
