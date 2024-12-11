package mockoauth

import (
	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/oauth"
)

type MockOAuth struct {
	users map[string]*entity.User
}

const ValidUserToken = "valid_user_token"

func NewMockOAuth() *MockOAuth {
	avatar := "https://avatars.githubusercontent.com/u/60039311"
	user := entity.User{
		ExternalID:    "6c2342aa-bdac-4efe-a31b-3a018072cff9",
		Name:          "John Doe",
		Email:         "johndoe@email.com",
		Avatar:        &avatar,
		Provider:      string(entity.ProviderGoogle),
		VerifiedEmail: true,
	}

	return &MockOAuth{
		users: map[string]*entity.User{
			ValidUserToken: &user,
		},
	}
}

func (m *MockOAuth) GetUser(token string) (*entity.User, error) {
	if user, ok := m.users[token]; ok {
		return user, nil
	}

	return nil, errs.New("invalid token")
}

var _ oauth.Provider = &MockOAuth{}
