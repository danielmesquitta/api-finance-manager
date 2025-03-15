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

func TestSignIn(t *testing.T) {
	t.Parallel()

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
			expectedCode:       http.StatusOK,
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
			expectedCode:       http.StatusOK,
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
				rawBody,
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

func TestRefreshToken(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description        string
		token              string
		expectedCode       int
		expectedUserAuthID string
	}{
		{
			description:        "Refresh token",
			token:              mockoauth.DefaultMockToken,
			expectedCode:       http.StatusOK,
			expectedUserAuthID: mockoauth.Users[mockoauth.DefaultMockToken].AuthID,
		},
		{
			description:        "Fail to refresh token without access token",
			token:              "",
			expectedCode:       http.StatusBadRequest,
			expectedUserAuthID: "",
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

			refreshToken := ""
			if test.token != "" {
				signInRes := app.SignIn(test.token)
				refreshToken = signInRes.RefreshToken
			}

			var out dto.SignInResponse
			statusCode, rawBody, err := app.MakeRequest(
				http.MethodPost,
				"/api/v1/auth/refresh",
				WithBearerToken(refreshToken),
				WithResponse(&out),
			)
			assert.Nil(t, err)

			assert.Equal(
				t,
				test.expectedCode,
				statusCode,
				rawBody,
			)

			if test.expectedCode != http.StatusOK {
				return
			}

			assert.NotEmpty(t, out.AccessToken)
			assert.NotEmpty(t, out.RefreshToken)
			assert.Equal(
				t,
				test.expectedUserAuthID,
				out.User.AuthID,
			)
		})
	}
}
