package restapi

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/oauth/mockoauth"
	"github.com/stretchr/testify/assert"
)

func TestCreateAIChat(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description  string
		token        string
		expectedCode int
		shouldCreate bool
	}{
		{
			description:  "Fail to create ai chat without token",
			token:        "",
			expectedCode: http.StatusBadRequest,
			shouldCreate: false,
		},
		{
			description:  "Should not create a new AI chat if the user's latest one is empty",
			token:        mockoauth.PremiumTierMockToken,
			expectedCode: http.StatusCreated,
			shouldCreate: false,
		},
		{
			description:  "Create a new AI chat",
			token:        mockoauth.TrialTierMockToken,
			expectedCode: http.StatusCreated,
			shouldCreate: true,
		},
		{
			description:  "Should not create a new AI chat for a free tier user",
			token:        mockoauth.FreeTierMockToken,
			expectedCode: http.StatusBadRequest,
			shouldCreate: false,
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()

			app, cleanUp := NewTestApp(t)
			defer func() {
				err := cleanUp(context.Background())
				assert.Nil(t, err)
			}()

			var (
				accessToken string
				user        *entity.User
			)
			if test.token != "" {
				signInRes := app.SignIn(test.token)
				accessToken, user = signInRes.AccessToken, &signInRes.User
			}

			statusCode, rawBody, err := app.MakeRequest(
				http.MethodPost,
				"/api/v1/ai-chats",
				WithBearerToken(accessToken),
			)
			assert.Nil(t, err)

			assert.Equal(
				t,
				test.expectedCode,
				statusCode,
				rawBody,
			)

			if test.expectedCode != http.StatusCreated {
				return
			}

			aiChat, err := app.db.GetLatestAIChatByUserID(
				ctx,
				user.ID,
			)
			assert.Nil(t, err)

			now := time.Now()
			if test.shouldCreate {
				assert.GreaterOrEqual(
					t,
					aiChat.CreatedAt.Unix()+60,
					now.Unix(),
				)
				assert.Nil(t, aiChat.Title)
			} else {
				assert.GreaterOrEqual(
					t,
					now.Unix(),
					aiChat.CreatedAt.Unix()+60,
				)
			}

		})
	}
}
