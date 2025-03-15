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

func TestCreateFeedback(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description  string
		token        string
		body         dto.CreateFeedbackRequest
		expectedCode int
	}{
		{
			description:  "Fail to create feedback without token",
			token:        "",
			expectedCode: http.StatusBadRequest,
		},
		{
			description:  "Create feedback",
			token:        mockoauth.DefaultMockToken,
			expectedCode: http.StatusCreated,
			body: dto.CreateFeedbackRequest{
				CreateFeedbackInput: usecase.CreateFeedbackInput{
					Message: "Loren ipsum dolor sit amet",
				},
			},
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
				"/api/v1/feedbacks",
				WithBearerToken(accessToken),
				WithBody(test.body),
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

			feedback, err := app.db.GetLatestFeedbackByUserID(
				ctx,
				user.ID.String(),
			)
			assert.Nil(t, err)

			expectedFeedback := map[string]any{
				"Message": test.body.Message,
				"UserID":  user.ID.String(),
			}

			actualFeedback := map[string]any{
				"Message": feedback.Message,
				"UserID":  feedback.UserID.String(),
			}

			assert.Equal(t, expectedFeedback, actualFeedback)
		})
	}
}
