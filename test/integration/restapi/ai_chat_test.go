package restapi

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/danielmesquitta/api-finance-manager/internal/app/restapi/dto"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/usecase"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/oauth/mockoauth"
	"github.com/google/uuid"
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

			signInRes := &dto.SignInResponse{}
			if test.token != "" {
				signInRes = app.SignIn(test.token)
			}

			statusCode, rawBody, err := app.MakeRequest(
				http.MethodPost,
				"/api/v1/ai-chats",
				WithBearerToken(signInRes.AccessToken),
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
				signInRes.User.ID,
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

func TestUpdateAIChat(t *testing.T) {
	t.Parallel()

	type Test struct {
		description    string
		token          string
		aiChatID       string
		body           dto.UpdateAIChatRequest
		expectedCode   int
		expectedAIChat *entity.AIChat
	}

	tests := []Test{
		{
			description:  "Fail to update ai chat without token",
			token:        "",
			aiChatID:     "e1c73c22-7d52-43e2-80a8-63ce6da99e53",
			expectedCode: http.StatusBadRequest,
		},
		{
			description:  "Fail to update AI chat with free tier",
			token:        mockoauth.FreeTierMockToken,
			aiChatID:     "df2017de-e019-4d14-b540-b31aafddffb8",
			expectedCode: http.StatusBadRequest,
			body: dto.UpdateAIChatRequest{
				UpdateAIChatInput: usecase.UpdateAIChatInput{
					Title: "Foo bar",
				},
			},
		},
		func() Test {
			title := "Foo bar"
			return Test{
				description:  "AI chat update",
				token:        mockoauth.PremiumTierMockToken,
				aiChatID:     "df2017de-e019-4d14-b540-b31aafddffb8",
				expectedCode: http.StatusNoContent,
				body: dto.UpdateAIChatRequest{
					UpdateAIChatInput: usecase.UpdateAIChatInput{
						Title: title,
					},
				},
				expectedAIChat: &entity.AIChat{
					Title: &title,
				},
			}
		}(),
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

			signInRes := &dto.SignInResponse{}
			if test.token != "" {
				signInRes = app.SignIn(test.token)
			}

			statusCode, rawBody, err := app.MakeRequest(
				http.MethodPut,
				"/api/v1/ai-chats/"+test.aiChatID,
				WithBearerToken(signInRes.AccessToken),
				WithBody(test.body),
			)
			assert.Nil(t, err)

			assert.Equal(
				t,
				test.expectedCode,
				statusCode,
				rawBody,
			)

			if test.expectedCode != http.StatusNoContent {
				return
			}

			actualAIChat, err := app.db.GetAIChat(
				ctx,
				uuid.MustParse(test.aiChatID),
			)
			assert.Nil(t, err)

			assert.Equal(
				t,
				test.expectedAIChat.Title,
				actualAIChat.Title,
			)
		})
	}
}
