package restapi

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/danielmesquitta/api-finance-manager/internal/app/restapi/dto"
	"github.com/danielmesquitta/api-finance-manager/internal/app/restapi/handler"
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
			description:  "fails without token",
			token:        "",
			expectedCode: http.StatusBadRequest,
			shouldCreate: false,
		},
		{
			description:  "should not create a new AI chat if the user's latest one is empty",
			token:        mockoauth.PremiumTierMockToken,
			expectedCode: http.StatusCreated,
			shouldCreate: false,
		},
		{
			description:  "creates a new AI chat",
			token:        mockoauth.TrialTierMockToken,
			expectedCode: http.StatusCreated,
			shouldCreate: true,
		},
		{
			description:  "should not create a new AI chat for a free tier user",
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

func TestGenerateAIChatMessage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description  string
		token        string
		expectedCode int
	}{
		{
			description:  "fails without token",
			token:        "",
			expectedCode: http.StatusBadRequest,
		},
		{
			description:  "should not create a new AI chat if the user's latest one is empty",
			token:        mockoauth.PremiumTierMockToken,
			expectedCode: http.StatusCreated,
		},
		{
			description:  "creates a new AI chat",
			token:        mockoauth.TrialTierMockToken,
			expectedCode: http.StatusCreated,
		},
		{
			description:  "should not create a new AI chat for a free tier user",
			token:        mockoauth.FreeTierMockToken,
			expectedCode: http.StatusBadRequest,
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
			description:  "fails without token",
			token:        "",
			aiChatID:     "e1c73c22-7d52-43e2-80a8-63ce6da99e53",
			expectedCode: http.StatusBadRequest,
		},
		{
			description:  "fails with free tier",
			token:        mockoauth.FreeTierMockToken,
			aiChatID:     "df2017de-e019-4d14-b540-b31aafddffb8",
			expectedCode: http.StatusBadRequest,
			body: dto.UpdateAIChatRequest{
				UpdateAIChatInput: usecase.UpdateAIChatInput{
					Title: "Foo bar",
				},
			},
		},
		{
			description:  "fails with non-existing",
			token:        mockoauth.PremiumTierMockToken,
			aiChatID:     "5fde4a75-f4df-415e-86bb-d7e24d488e36",
			expectedCode: http.StatusNotFound,
			body: dto.UpdateAIChatRequest{
				UpdateAIChatInput: usecase.UpdateAIChatInput{
					Title: "Non-existing AI chat",
				},
			},
		},
		func() Test {
			title := "Lorem ipsum"
			return Test{
				description:  "updates ai chat",
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

			actualAIChat, err := app.db.GetAIChatByID(
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

func TestDeleteAIChat(t *testing.T) {
	t.Parallel()

	type Test struct {
		description  string
		token        string
		expectedCode int
		aiChatID     string
	}

	tests := []Test{
		{
			description:  "fails without token",
			token:        "",
			expectedCode: http.StatusBadRequest,
			aiChatID:     "df2017de-e019-4d14-b540-b31aafddffb8",
		},
		{
			description:  "fails with non-existing",
			token:        mockoauth.PremiumTierMockToken,
			expectedCode: http.StatusNotFound,
			aiChatID:     "b8f7cc16-157a-48a8-8c04-287754599e3e",
		},
		{
			description:  "deletes ai chat",
			token:        mockoauth.PremiumTierMockToken,
			expectedCode: http.StatusNoContent,
			aiChatID:     "df2017de-e019-4d14-b540-b31aafddffb8",
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
				http.MethodDelete,
				"/api/v1/ai-chats/"+test.aiChatID,
				WithBearerToken(signInRes.AccessToken),
			)
			assert.Nil(t, err)

			assert.Equal(
				t,
				test.expectedCode,
				statusCode,
				rawBody,
			)

			actualAIChat, err := app.db.GetLatestUserDeletedAIChat(
				ctx,
				signInRes.User.ID,
			)
			assert.Nil(t, err)

			if test.expectedCode != http.StatusNoContent {
				assert.NotEqual(t, test.aiChatID, actualAIChat.ID.String())
				return
			}

			assert.NotNil(t, actualAIChat.DeletedAt)
			assert.Equal(t, test.aiChatID, actualAIChat.ID.String())
		})
	}
}

func TestListAIChatMessagesAndAnswers(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description              string
		aiChatID                 string
		queryParams              map[string]string
		token                    string
		expectedCode             int
		expectedAIChatMessageIDs []string
	}{
		{
			description:              "fails without token",
			aiChatID:                 "9945780b-c3f8-4464-a83d-e063d2faf93d",
			queryParams:              map[string]string{},
			token:                    "",
			expectedCode:             http.StatusBadRequest,
			expectedAIChatMessageIDs: []string{},
		},
		{
			description:  "lists ai chat messages",
			aiChatID:     "9945780b-c3f8-4464-a83d-e063d2faf93d",
			queryParams:  map[string]string{},
			token:        mockoauth.PremiumTierMockToken,
			expectedCode: http.StatusOK,
			expectedAIChatMessageIDs: []string{
				"f0ff7683-1d70-4d8d-bb83-281add648dbe",
				"952ec6a8-f550-4de3-bb0b-6dc2e937b8e6",
				"967253d1-965f-4ae7-90b1-5990e103dcac",
				"c68f4fe7-9cde-411f-86ef-b62e0daaa58f",
			},
		},
		{
			description: "paginates ai chat messages",
			aiChatID:    "9945780b-c3f8-4464-a83d-e063d2faf93d",
			queryParams: map[string]string{
				handler.QueryParamPageSize: "2",
			},
			token:        mockoauth.PremiumTierMockToken,
			expectedCode: http.StatusOK,
			expectedAIChatMessageIDs: []string{
				"c68f4fe7-9cde-411f-86ef-b62e0daaa58f",
				"952ec6a8-f550-4de3-bb0b-6dc2e937b8e6",
			},
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

			var out dto.ListAIChatMessagesAndAnswersResponse
			var errRes dto.ErrorResponse
			statusCode, rawBody, err := app.MakeRequest(
				http.MethodGet,
				fmt.Sprintf("/api/v1/ai-chats/%s/messages", test.aiChatID),
				WithQueryParams(test.queryParams),
				WithBearerToken(signInRes.AccessToken),
				WithResponse(&out),
				WithError(&errRes),
			)
			assert.Nil(t, err)

			assert.Equal(
				t,
				test.expectedCode,
				statusCode,
				rawBody,
			)

			if len(test.expectedAIChatMessageIDs) == 0 {
				assert.Empty(t, out.Items)
				return
			}

			assert.Len(
				t,
				out.Items,
				len(test.expectedAIChatMessageIDs),
			)

			aiChatMessageIDs := make([]string, len(out.Items))
			for i, aiChatMessage := range out.Items {
				aiChatMessageIDs[i] = aiChatMessage.ID.String()
			}

			assert.ElementsMatch(
				t,
				test.expectedAIChatMessageIDs,
				aiChatMessageIDs,
			)
		})
	}
}
