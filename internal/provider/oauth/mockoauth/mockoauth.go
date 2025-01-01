package mockoauth

import (
	"context"

	"github.com/danielmesquitta/api-finance-manager/internal/config"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/oauth"
)

type MockOAuth struct {
	e *config.Env

	users map[string]*entity.User
}

const ValidUserToken = "valid_user_token"

func NewMockOAuth(
	e *config.Env,
) *MockOAuth {
	if e.Environment == config.EnvironmentProduction {
		panic("mock oauth not allowed in production")
	}

	avatar := "https://avatars.githubusercontent.com/u/60039311"
	user := entity.User{
		ExternalID:    "6c2342aa-bdac-4efe-a31b-3a018072cff9",
		Name:          "John Doe",
		Email:         "johndoe@email.com",
		Avatar:        &avatar,
		Provider:      string(entity.ProviderMock),
		VerifiedEmail: true,
	}

	return &MockOAuth{
		e: e,
		users: map[string]*entity.User{
			ValidUserToken: &user,
		},
	}
}

func (m *MockOAuth) GetUser(
	ctx context.Context,
	token string,
) (*entity.User, error) {
	if user, ok := m.users[token]; ok {
		return user, nil
	}

	return nil, errs.ErrUnauthorized
}

var _ oauth.Provider = &MockOAuth{}
