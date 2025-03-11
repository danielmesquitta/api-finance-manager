package restapi

import (
	"context"
	"net/http"
	"testing"

	"github.com/danielmesquitta/api-finance-manager/internal/app/restapi/dto"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/oauth/mockoauth"
	"github.com/stretchr/testify/assert"
)

func TestListTransactionCategoriesRoute(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description                    string
		queryParams                    map[string]string
		token                          string
		expectedCode                   int
		expectedTransactionCategoryIDs []string
	}{
		{
			description:                    "Fail to list transaction categories without token",
			queryParams:                    map[string]string{},
			token:                          "",
			expectedCode:                   http.StatusBadRequest,
			expectedTransactionCategoryIDs: []string{},
		},
		{
			description: "List all transaction categories",
			queryParams: map[string]string{
				"page_size": "100",
			},
			token:        mockoauth.DefaultMockToken,
			expectedCode: http.StatusOK,
			expectedTransactionCategoryIDs: []string{
				"65583cfa-b72d-4fab-9de1-4ca9dfe11a4e",
				"432130c4-cd6f-4856-a82c-54f227b0382b",
				"029fc6cb-edcf-414c-9c81-9dd69c34e629",
				"f1f7c2ee-f797-4014-ae7e-ab40bee5afdd",
				"a910d4f6-2904-4b1e-a76d-aa04515eb966",
				"b1dcc2c7-d889-4df8-b360-8c20904f7e08",
				"c198bbd9-3f06-42cd-b2be-b212404d83fc",
				"03bd0abc-7186-4eb3-9871-e4f624c535b8",
				"9db43714-3025-494f-9578-4feb5a69681e",
				"12deb35c-0ce5-4d23-87a4-2f68fd77f019",
				"70c89492-9977-42b4-a28c-b1e261c59615",
				"02701aac-b8db-4c7e-834c-6d4f4eab3399",
				"42028a90-e209-4853-9bca-d949f3cec9e6",
				"dba1995e-86c9-474f-a482-92afb6f71615",
				"d62f49b0-7aa6-4346-a7e4-bb156b0a99d4",
				"f8ad2a99-f062-4c46-8582-d298003b46c0",
				"0c84d0a3-7336-4089-bc3d-756ce31c679a",
				"1bd7db5b-5b8a-4ac1-82a1-c75e418a25c0",
				"b06e0c42-4053-4fad-b289-be0cfc22502c",
				"5de3003d-f3c5-4d80-b118-80191e59645d",
				"c03e9602-ac21-4cc4-b2f4-ae0dbb4b8dfb",
				"ee1cde91-ac0f-4b9a-b8ed-1726a83f3643",
			},
		},
		{
			description: "Search transaction categories",
			queryParams: map[string]string{
				"search": "Emprestimos",
			},
			token:        mockoauth.DefaultMockToken,
			expectedCode: http.StatusOK,
			expectedTransactionCategoryIDs: []string{
				"432130c4-cd6f-4856-a82c-54f227b0382b",
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

			var out dto.ListTransactionCategoriesResponse
			statusCode, rawBody, err := app.MakeRequest(
				http.MethodGet,
				"/api/v1/transactions/categories",
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

			if len(test.expectedTransactionCategoryIDs) == 0 {
				assert.Empty(t, out.Items, rawBody)
				return
			}

			assert.Len(
				t,
				out.Items,
				len(test.expectedTransactionCategoryIDs),
			)

			transactionCategoryIDs := make([]string, len(out.Items))
			for i, transactionCategory := range out.Items {
				transactionCategoryIDs[i] = transactionCategory.ID.String()
			}

			assert.ElementsMatch(
				t,
				test.expectedTransactionCategoryIDs,
				transactionCategoryIDs,
			)
		})
	}
}
