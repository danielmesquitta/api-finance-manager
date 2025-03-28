package mockoauth

import (
	"context"
	"errors"
	"time"

	"github.com/danielmesquitta/api-finance-manager/internal/config/env"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/ptr"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/oauth"
)

type MockOAuth struct {
	e *env.Env
}

const PremiumTierMockToken = "premium_mock_token"
const TrialTierMockToken = "trial_mock_token"
const FreeTierMockToken = "free_tier_mock_token"
const UnregisteredUserMockToken = "unregistered_user_mock_token"

type User struct {
	*entity.User
	AuthProvider *entity.UserAuthProvider
}

var Users = map[string]*User{
	PremiumTierMockToken: {
		User: &entity.User{
			Name:  "John Doe",
			Email: "johndoe@email.com",
			Avatar: ptr.New(
				"https://avatar.iran.liara.run/public/15",
			),
			Tier:                  entity.TierPremium,
			SubscriptionExpiresAt: ptr.New(time.Now().AddDate(0, 1, 0)),
		},
		AuthProvider: &entity.UserAuthProvider{
			ExternalID:    "6c2342aa-bdac-4efe-a31b-3a018072cff9",
			Provider:      entity.ProviderMock,
			VerifiedEmail: true,
		},
	},
	TrialTierMockToken: {
		User: &entity.User{
			Name:                  "Jennifer Doe",
			Email:                 "jenniferdoe@email.com",
			Tier:                  entity.TierTrial,
			SubscriptionExpiresAt: ptr.New(time.Now().AddDate(0, 1, 0)),
		},
		AuthProvider: &entity.UserAuthProvider{
			ExternalID:    "2a35fa25-2809-40d7-beeb-0d2766171b1d",
			Provider:      entity.ProviderMock,
			VerifiedEmail: true,
		},
	},
	FreeTierMockToken: {
		User: &entity.User{
			Name:   "Jane Doe",
			Email:  "janedoe@email.com",
			Avatar: ptr.New("https://avatar.iran.liara.run/public/82"),
			Tier:   entity.TierFree,
		},
		AuthProvider: &entity.UserAuthProvider{
			ExternalID:    "016aecbd-fae5-4ff0-9046-03b7eabf6a5c",
			Provider:      entity.ProviderMock,
			VerifiedEmail: true,
		},
	},
	UnregisteredUserMockToken: {
		User: &entity.User{
			Name:   "Joseph Doe",
			Email:  "josephdoe@email.com",
			Avatar: ptr.New("https://avatar.iran.liara.run/public/13"),
			Tier:   entity.TierFree,
		},
		AuthProvider: &entity.UserAuthProvider{
			ExternalID:    "2824923b-2d93-4473-8397-32680bb412b4",
			Provider:      entity.ProviderMock,
			VerifiedEmail: true,
		},
	},
}

func NewMockOAuth(
	e *env.Env,
) *MockOAuth {
	if e.Environment == env.EnvironmentProduction {
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
		user, ok = Users[PremiumTierMockToken]
		if !ok || user == nil {
			return nil, nil, errors.New("user not found")
		}
	}

	return user.User, user.AuthProvider, nil
}

var _ oauth.Provider = (*MockOAuth)(nil)
