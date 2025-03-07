package restapi

import (
	"context"
	"net/http"
	"testing"

	"github.com/danielmesquitta/api-finance-manager/internal/app/restapi/dto"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/usecase"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/oauth/mockoauth"
	"github.com/stretchr/testify/assert"
)

func TestSignInRoute(t *testing.T) {
	tests := []struct {
		description        string
		body               *dto.SignInRequest
		token              string
		expectedCode       int
		expectedUserAuthID string
	}{
		{
			description: "Sign in with registered user",
			body: &dto.SignInRequest{
				SignInInput: usecase.SignInInput{
					Provider: entity.ProviderMock,
				},
			},
			token:              mockoauth.DefaultMockToken,
			expectedCode:       200,
			expectedUserAuthID: mockoauth.Users[mockoauth.DefaultMockToken].AuthID,
		},
		{
			description: "Sign in with unregistered user",
			body: &dto.SignInRequest{
				SignInInput: usecase.SignInInput{
					Provider: entity.ProviderMock,
				},
			},
			token:              mockoauth.UnregisteredUserMockToken,
			expectedCode:       200,
			expectedUserAuthID: mockoauth.Users[mockoauth.UnregisteredUserMockToken].AuthID,
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			t.Parallel()

			app, cleanUp := NewTestApp(t)
			defer func() {
				err := cleanUp(context.Background())
				assert.Nil(t, err)
			}()

			var out dto.SignInResponse
			statusCode, rawBody, err := app.MakeRequest(
				http.MethodPost,
				"/api/v1/auth/sign-in",
				WithBody(test.body),
				WithToken(test.token),
				WithResponse(&out),
			)
			assert.Nil(t, err)

			assert.Equal(
				t,
				test.expectedCode,
				statusCode,
			)

			assert.NotEmpty(t, out.AccessToken, rawBody)
			assert.NotEmpty(t, out.RefreshToken, rawBody)
			assert.Equal(
				t,
				test.expectedUserAuthID,
				out.User.AuthID,
				rawBody,
			)
		})
	}
}
