package server

import (
	"context"
	"net/http"
	"testing"

	"github.com/danielmesquitta/api-finance-manager/internal/app/server/dto"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/usecase/user"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/oauth/mockoauth"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetProfile(t *testing.T) {
	t.Parallel()

	type Test struct {
		description      string
		token            string
		expectedCode     int
		expectedResponse *dto.GetUserProfileResponse
	}

	tests := []Test{
		{
			description:  "fails without token",
			token:        "",
			expectedCode: http.StatusBadRequest,
		},
		func() Test {
			avatar := "https://avatar.iran.liara.run/public/15"
			return Test{
				description:  "gets premium tier profile",
				token:        mockoauth.PremiumTierMockToken,
				expectedCode: http.StatusOK,
				expectedResponse: &dto.GetUserProfileResponse{
					User: entity.User{
						ID: uuid.MustParse(
							"fdfdc888-da64-4988-8ad3-f739862c4ceb",
						),
						Name:   "John Doe",
						Email:  "johndoe@email.com",
						Tier:   entity.TierPremium,
						Avatar: &avatar,
					},
				},
			}
		}(),
		func() Test {
			avatar := "https://avatar.iran.liara.run/public/82"
			return Test{
				description:  "gets free tier profile",
				token:        mockoauth.FreeTierMockToken,
				expectedCode: http.StatusOK,
				expectedResponse: &dto.GetUserProfileResponse{
					User: entity.User{
						ID: uuid.MustParse(
							"5b4694a9-c810-41a2-bca6-74c3f3850fe7",
						),
						Name:   "Jane Doe",
						Email:  "janedoe@email.com",
						Tier:   entity.TierFree,
						Avatar: &avatar,
					},
				},
			}
		}(),
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

			var actualResponse dto.GetUserProfileResponse
			statusCode, rawBody, err := app.MakeRequest(
				http.MethodGet,
				"/api/v1/users/profile",
				WithBearerToken(signInRes.AccessToken),
				WithResponse(&actualResponse),
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

			assert.Equal(t, test.expectedResponse.ID, actualResponse.ID)
			assert.Equal(t, test.expectedResponse.Name, actualResponse.Name)
			assert.Equal(t, test.expectedResponse.Email, actualResponse.Email)
			assert.Equal(t, test.expectedResponse.Tier, actualResponse.Tier)
			assert.Equal(t, test.expectedResponse.Avatar, actualResponse.Avatar)
		})
	}
}

func TestUpdateProfile(t *testing.T) {
	t.Parallel()

	type Test struct {
		description  string
		token        string
		body         dto.UpdateProfileRequest
		expectedCode int
		expectedUser *entity.User
	}

	tests := []Test{
		{
			description:  "fails without token",
			token:        "",
			expectedCode: http.StatusBadRequest,
		},
		func() Test {
			id := uuid.MustParse("fdfdc888-da64-4988-8ad3-f739862c4ceb")
			name := "Johnathan Doe"
			email := "johnathandoe@gmail.com"
			return Test{
				description:  "updates full user",
				token:        mockoauth.PremiumTierMockToken,
				expectedCode: http.StatusNoContent,
				body: dto.UpdateProfileRequest{
					UpdateUserUseCaseInput: user.UpdateUserUseCaseInput{
						ID:    id,
						Name:  name,
						Email: email,
					},
				},
				expectedUser: &entity.User{
					ID:    id,
					Name:  name,
					Email: email,
				},
			}
		}(),
		func() Test {
			id := uuid.MustParse("fdfdc888-da64-4988-8ad3-f739862c4ceb")
			name := "Johnathan Doe"
			return Test{
				description:  "updates partial user",
				token:        mockoauth.PremiumTierMockToken,
				expectedCode: http.StatusNoContent,
				body: dto.UpdateProfileRequest{
					UpdateUserUseCaseInput: user.UpdateUserUseCaseInput{
						ID:   id,
						Name: name,
					},
				},
				expectedUser: &entity.User{
					ID:    id,
					Name:  name,
					Email: "johndoe@email.com",
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
				"/api/v1/users/profile",
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

			actualUser, err := app.db.GetUserByID(
				ctx,
				signInRes.User.ID,
			)
			assert.Nil(t, err)

			assert.Equal(t, test.expectedUser.ID, actualUser.ID)
			assert.Equal(t, test.expectedUser.Name, actualUser.Name)
			assert.Equal(t, test.expectedUser.Email, actualUser.Email)
		})
	}
}

func TestDeleteProfile(t *testing.T) {
	t.Parallel()

	type Test struct {
		description  string
		token        string
		expectedCode int
	}

	tests := []Test{
		{
			description:  "fails without token",
			token:        "",
			expectedCode: http.StatusBadRequest,
		},
		{
			description:  "deletes user",
			token:        mockoauth.PremiumTierMockToken,
			expectedCode: http.StatusNoContent,
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
				"/api/v1/users/profile",
				WithBearerToken(signInRes.AccessToken),
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

			actualUser, err := app.db.GetLatestDeletedUser(
				ctx,
				signInRes.User.ID,
			)
			assert.Nil(t, err)

			assert.NotNil(t, actualUser.DeletedAt)

			userBeforeDeletion, ok := mockoauth.Users[test.token]
			assert.True(t, ok)
			assert.NotEqual(t, userBeforeDeletion.Email, actualUser.Email)
			assert.NotEqual(t, userBeforeDeletion.Name, actualUser.Name)
		})
	}
}
