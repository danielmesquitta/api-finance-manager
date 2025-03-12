package restapi

import (
	"context"
	"net/http"
	"testing"

	"github.com/danielmesquitta/api-finance-manager/internal/app/restapi/dto"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/oauth/mockoauth"
	"github.com/stretchr/testify/assert"
)

func TestListInstitutionsRoute(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description            string
		queryParams            map[string]string
		token                  string
		expectedCode           int
		expectedInstitutionIDs []string
	}{
		{
			description:            "Fail to list institutions without token",
			queryParams:            map[string]string{},
			token:                  "",
			expectedCode:           http.StatusBadRequest,
			expectedInstitutionIDs: []string{},
		},
		{
			description:  "List all institutions",
			queryParams:  map[string]string{},
			token:        mockoauth.DefaultMockToken,
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
			description: "Paginate institutions",
			queryParams: map[string]string{
				"page":      "2",
				"page_size": "2",
			},
			token:        mockoauth.DefaultMockToken,
			expectedCode: http.StatusOK,
			expectedInstitutionIDs: []string{
				"df5dbd97-89c7-4776-8b3f-7992bc2bb16b",
				"9daf7ad9-2d35-4597-8788-788d7f7f6c98",
			},
		},
		{
			description: "Search institutions",
			queryParams: map[string]string{
				"search": "Int",
			},
			token:        mockoauth.DefaultMockToken,
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

			accessToken := ""
			if test.token != "" {
				signInRes := app.SignIn(test.token)
				accessToken = signInRes.AccessToken
			}

			var out dto.ListInstitutionsResponse
			statusCode, rawBody, err := app.MakeRequest(
				http.MethodGet,
				"/api/v1/institutions",
				WithQueryParams(test.queryParams),
				WithBearerToken(accessToken),
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

func TestListUserInstitutionsRoute(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description            string
		queryParams            map[string]string
		token                  string
		expectedCode           int
		expectedInstitutionIDs []string
	}{
		{
			description:            "Fail to list institutions without token",
			queryParams:            map[string]string{},
			token:                  "",
			expectedCode:           http.StatusBadRequest,
			expectedInstitutionIDs: []string{},
		},
		{
			description:  "List all institutions",
			queryParams:  map[string]string{},
			token:        mockoauth.DefaultMockToken,
			expectedCode: http.StatusOK,
			expectedInstitutionIDs: []string{
				"88f812ab-9bc9-4830-afc6-7ac0ba67b1ec",
				"df5dbd97-89c7-4776-8b3f-7992bc2bb16b",
			},
		},
		{
			description: "Paginate institutions",
			queryParams: map[string]string{
				"page":      "1",
				"page_size": "1",
			},
			token:        mockoauth.DefaultMockToken,
			expectedCode: http.StatusOK,
			expectedInstitutionIDs: []string{
				"88f812ab-9bc9-4830-afc6-7ac0ba67b1ec",
			},
		},
		{
			description: "Search institutions",
			queryParams: map[string]string{
				"search": "nub",
			},
			token:        mockoauth.DefaultMockToken,
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

			accessToken := ""
			if test.token != "" {
				signInRes := app.SignIn(test.token)
				accessToken = signInRes.AccessToken
			}

			var out dto.ListInstitutionsResponse
			statusCode, rawBody, err := app.MakeRequest(
				http.MethodGet,
				"/api/v1/users/institutions",
				WithQueryParams(test.queryParams),
				WithBearerToken(accessToken),
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
