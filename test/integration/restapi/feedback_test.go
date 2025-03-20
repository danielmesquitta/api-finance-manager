package restapi

import (
	"context"
	"net/http"
	"testing"

	"github.com/danielmesquitta/api-finance-manager/internal/app/restapi/dto"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/usecase/feedback"
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
			description:  "fails without token",
			token:        "",
			expectedCode: http.StatusBadRequest,
		},
		{
			description:  "creates feedback",
			token:        mockoauth.PremiumTierMockToken,
			expectedCode: http.StatusCreated,
			body: dto.CreateFeedbackRequest{
				CreateFeedbackUseCaseInput: feedback.CreateFeedbackUseCaseInput{
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

			signInRes := &dto.SignInResponse{}
			if test.token != "" {
				signInRes = app.SignIn(test.token)
			}

			statusCode, rawBody, err := app.MakeRequest(
				http.MethodPost,
				"/api/v1/feedbacks",
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

			if test.expectedCode != http.StatusCreated {
				return
			}

			feedback, err := app.db.GetLatestFeedbackByUserID(
				ctx,
				signInRes.User.ID.String(),
			)
			assert.Nil(t, err)

			assert.Equal(t, test.body.Message, feedback.Message)
			assert.Equal(
				t,
				signInRes.User.ID.String(),
				feedback.UserID.String(),
			)
		})
	}
}
