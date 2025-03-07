package restapi

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/danielmesquitta/api-finance-manager/internal/app/restapi/dto"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/usecase"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/oauth/mockoauth"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestSignInRoute(t *testing.T) {
	tests := []struct {
		description        string
		body               *dto.SignInRequest
		headers            map[string]string
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
			headers: map[string]string{
				fiber.HeaderContentType:   fiber.MIMEApplicationJSON,
				fiber.HeaderAuthorization: "Bearer " + mockoauth.MockToken,
			},
			expectedCode:       200,
			expectedUserAuthID: mockoauth.Users[mockoauth.MockToken].AuthID,
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			t.Parallel()

			app, cleanUp := Setup(t)
			defer func() {
				err := cleanUp(context.Background())
				assert.Nil(t, err)
			}()

			jsonBody, err := json.Marshal(test.body)
			assert.Nil(t, err)

			req, _ := http.NewRequest(
				http.MethodPost,
				"/api/v1/auth/sign-in",
				bytes.NewReader(jsonBody),
			)

			for key, value := range test.headers {
				req.Header.Add(key, value)
			}

			res, err := app.Test(req, -1)
			assert.Nil(t, err)

			assert.Equal(
				t,
				test.expectedCode,
				res.StatusCode,
			)

			bytesBody, err := io.ReadAll(res.Body)
			assert.Nil(t, err)

			var response dto.SignInResponse
			err = json.Unmarshal(bytesBody, &response)
			assert.Nil(t, err)

			assert.NotEmpty(t, response.AccessToken, string(bytesBody))
			assert.NotEmpty(t, response.RefreshToken, string(bytesBody))
			assert.Equal(
				t,
				test.expectedUserAuthID,
				response.User.AuthID,
				string(bytesBody),
			)
		})
	}
}
