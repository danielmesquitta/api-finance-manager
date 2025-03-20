package server

import (
	"context"
	"net/http"
	"testing"

	"github.com/danielmesquitta/api-finance-manager/internal/app/server/dto"
	"github.com/danielmesquitta/api-finance-manager/internal/app/server/handler"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/oauth/mockoauth"
	"github.com/stretchr/testify/assert"
)

func TestListInstitutions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description            string
		queryParams            map[string]string
		token                  string
		expectedCode           int
		expectedInstitutionIDs []string
	}{
		{
			description:            "fails without token",
			queryParams:            map[string]string{},
			token:                  "",
			expectedCode:           http.StatusBadRequest,
			expectedInstitutionIDs: []string{},
		},
		{
			description:  "lists institutions",
			queryParams:  map[string]string{},
			token:        mockoauth.PremiumTierMockToken,
			expectedCode: http.StatusOK,
			expectedInstitutionIDs: []string{
				"eb3a4329-ba36-4000-9123-748e2c1fdd60",
				"9daf7ad9-2d35-4597-8788-788d7f7f6c98",
				"88f812ab-9bc9-4830-afc6-7ac0ba67b1ec",
				"66a1475d-94d6-4848-b4c1-61a91f8317f3",
				"df5dbd97-89c7-4776-8b3f-7992bc2bb16b",
				"e250139e-0c18-4368-9f7f-5b546740a6f8",
			},
		},
		{
			description: "paginates institutions",
			queryParams: map[string]string{
				handler.QueryParamPage:     "2",
				handler.QueryParamPageSize: "2",
			},
			token:        mockoauth.PremiumTierMockToken,
			expectedCode: http.StatusOK,
			expectedInstitutionIDs: []string{
				"df5dbd97-89c7-4776-8b3f-7992bc2bb16b",
				"9daf7ad9-2d35-4597-8788-788d7f7f6c98",
			},
		},
		{
			description: "searches institutions",
			queryParams: map[string]string{
				handler.QueryParamSearch: "inter",
			},
			token:        mockoauth.PremiumTierMockToken,
			expectedCode: http.StatusOK,
			expectedInstitutionIDs: []string{
				"eb3a4329-ba36-4000-9123-748e2c1fdd60",
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

			var out dto.ListInstitutionsResponse
			statusCode, rawBody, err := app.MakeRequest(
				http.MethodGet,
				"/api/v1/institutions",
				WithQueryParams(test.queryParams),
				WithBearerToken(signInRes.AccessToken),
				WithResponse(&out),
			)
			assert.Nil(t, err)

			assert.Equal(
				t,
				test.expectedCode,
				statusCode,
				rawBody,
			)

			if len(test.expectedInstitutionIDs) == 0 {
				assert.Empty(t, out.Items)
				return
			}

			assert.Len(
				t,
				out.Items,
				len(test.expectedInstitutionIDs),
			)

			institutionIDs := make([]string, len(out.Items))
			for i, institution := range out.Items {
				institutionIDs[i] = institution.ID.String()
			}

			assert.ElementsMatch(
				t,
				test.expectedInstitutionIDs,
				institutionIDs,
			)
		})
	}
}

func TestListUserInstitutions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description            string
		queryParams            map[string]string
		token                  string
		expectedCode           int
		expectedInstitutionIDs []string
	}{
		{
			description:            "fails without token",
			queryParams:            map[string]string{},
			token:                  "",
			expectedCode:           http.StatusBadRequest,
			expectedInstitutionIDs: []string{},
		},
		{
			description:  "lists user institutions",
			queryParams:  map[string]string{},
			token:        mockoauth.PremiumTierMockToken,
			expectedCode: http.StatusOK,
			expectedInstitutionIDs: []string{
				"88f812ab-9bc9-4830-afc6-7ac0ba67b1ec",
				"df5dbd97-89c7-4776-8b3f-7992bc2bb16b",
			},
		},
		{
			description: "paginates institutions",
			queryParams: map[string]string{
				handler.QueryParamPage:     "1",
				handler.QueryParamPageSize: "1",
			},
			token:        mockoauth.PremiumTierMockToken,
			expectedCode: http.StatusOK,
			expectedInstitutionIDs: []string{
				"88f812ab-9bc9-4830-afc6-7ac0ba67b1ec",
			},
		},
		{
			description: "searches user institutions",
			queryParams: map[string]string{
				handler.QueryParamSearch: "nubank",
			},
			token:        mockoauth.PremiumTierMockToken,
			expectedCode: http.StatusOK,
			expectedInstitutionIDs: []string{
				"88f812ab-9bc9-4830-afc6-7ac0ba67b1ec",
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

			var out dto.ListInstitutionsResponse
			statusCode, rawBody, err := app.MakeRequest(
				http.MethodGet,
				"/api/v1/users/institutions",
				WithQueryParams(test.queryParams),
				WithBearerToken(signInRes.AccessToken),
				WithResponse(&out),
			)
			assert.Nil(t, err)

			assert.Equal(
				t,
				test.expectedCode,
				statusCode,
				rawBody,
			)

			if len(test.expectedInstitutionIDs) == 0 {
				assert.Empty(t, out.Items)
				return
			}

			assert.Len(
				t,
				out.Items,
				len(test.expectedInstitutionIDs),
			)

			institutionIDs := make([]string, len(out.Items))
			for i, institution := range out.Items {
				institutionIDs[i] = institution.ID.String()
			}

			assert.ElementsMatch(
				t,
				test.expectedInstitutionIDs,
				institutionIDs,
			)
		})
	}
}
