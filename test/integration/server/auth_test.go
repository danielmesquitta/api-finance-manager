package server

import (
	"context"
	"net/http"
	"testing"

	"github.com/danielmesquitta/api-finance-manager/internal/app/server/dto"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/usecase/auth"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/oauth/mockoauth"
	"github.com/stretchr/testify/assert"
)

func TestSignIn(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description   string
		body          *dto.SignInRequest
		token         string
		expectedCode  int
		expectedEmail string
	}{
		{
			description: "signs in with registered user",
			body: &dto.SignInRequest{
				SignInUseCaseInput: auth.SignInUseCaseInput{
					Provider: entity.ProviderMock,
				},
			},
			token:         mockoauth.PremiumTierMockToken,
			expectedCode:  http.StatusOK,
			expectedEmail: mockoauth.Users[mockoauth.PremiumTierMockToken].Email,
		},
		{
			description: "signs in with unregistered user",
			body: &dto.SignInRequest{
				SignInUseCaseInput: auth.SignInUseCaseInput{
					Provider: entity.ProviderMock,
				},
			},
			token:         mockoauth.UnregisteredUserMockToken,
			expectedCode:  http.StatusOK,
			expectedEmail: mockoauth.Users[mockoauth.UnregisteredUserMockToken].Email,
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

			var actual dto.SignInResponse
			statusCode, rawBody, err := app.MakeRequest(
				http.MethodPost,
				"/api/v1/auth/sign-in",
				WithBody(test.body),
				WithToken(test.token),
				WithResponse(&actual),
			)
			assert.Nil(t, err)

			assert.Equal(
				t,
				test.expectedCode,
				statusCode,
				rawBody,
			)

			assert.NotEmpty(t, actual.AccessToken)
			assert.NotEmpty(t, actual.RefreshToken)
			assert.Equal(
				t,
				test.expectedEmail,
				actual.User.Email,
			)
		})
	}
}

func TestRefreshToken(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description   string
		token         string
		expectedCode  int
		expectedEmail string
	}{
		{
			description:   "refreshes token",
			token:         mockoauth.PremiumTierMockToken,
			expectedCode:  http.StatusOK,
			expectedEmail: mockoauth.Users[mockoauth.PremiumTierMockToken].Email,
		},
		{
			description:   "fails without access token",
			token:         "",
			expectedCode:  http.StatusBadRequest,
			expectedEmail: "",
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

			signInRes := &dto.SignInResponse{}
			if test.token != "" {
				signInRes = app.SignIn(test.token)
			}

			var actual dto.SignInResponse
			statusCode, rawBody, err := app.MakeRequest(
				http.MethodPost,
				"/api/v1/auth/refresh",
				WithBearerToken(signInRes.RefreshToken),
				WithResponse(&actual),
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

			assert.NotEmpty(t, actual.AccessToken)
			assert.NotEmpty(t, actual.RefreshToken)
			assert.Equal(
				t,
				test.expectedEmail,
				actual.User.Email,
			)
		})
	}
}
