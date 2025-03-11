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
				"c2368239-8286-41e5-a905-ac919a551699",
				"a770232a-0feb-46f1-bf77-96f938a58ec9",
				"ce7c4efd-c74a-4eb8-b8b3-7c02b704aa5e",
				"fcab2309-178f-4701-960f-147c3904388b",
				"1202269c-ed03-4dfe-bbcd-c61d615a17b5",
				"97fe5c01-6799-4831-b897-ec892c7368f9",
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
				"c2368239-8286-41e5-a905-ac919a551699",
				"a770232a-0feb-46f1-bf77-96f938a58ec9",
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
				"c2368239-8286-41e5-a905-ac919a551699",
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
			)

			if len(test.expectedInstitutionIDs) == 0 {
				assert.Empty(t, out.Items, rawBody)
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
				"ce7c4efd-c74a-4eb8-b8b3-7c02b704aa5e",
				"1202269c-ed03-4dfe-bbcd-c61d615a17b5",
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
				"1202269c-ed03-4dfe-bbcd-c61d615a17b5",
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
				"ce7c4efd-c74a-4eb8-b8b3-7c02b704aa5e",
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
			)

			if len(test.expectedInstitutionIDs) == 0 {
				assert.Empty(t, out.Items, rawBody)
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
