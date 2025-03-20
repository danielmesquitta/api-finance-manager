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

func TestListTransactionCategories(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description                    string
		queryParams                    map[string]string
		token                          string
		expectedCode                   int
		expectedTransactionCategoryIDs []string
	}{
		{
			description:                    "fails without tokenout token",
			queryParams:                    map[string]string{},
			token:                          "",
			expectedCode:                   http.StatusBadRequest,
			expectedTransactionCategoryIDs: []string{},
		},
		{
			description: "lists transaction categories",
			queryParams: map[string]string{
				handler.QueryParamPageSize: "100",
			},
			token:        mockoauth.PremiumTierMockToken,
			expectedCode: http.StatusOK,
			expectedTransactionCategoryIDs: []string{
				"30996ce1-e273-4cf3-9eb6-df0e1df03bba",
				"470850a0-781b-4892-8e96-5590d16c9fe1",
				"56a3416d-e633-4c60-a3e8-f4649474c14b",
				"a5a179a6-37be-43c6-931b-388d7f928f76",
				"373b150b-94bd-44b2-abdd-2aab14e74fad",
				"5cdb75ae-f27a-4ebf-9135-ec46273cdeea",
				"2a226707-1f75-4276-9697-3e7aac3c7db6",
				"896d5ff8-1534-4d4f-aa1f-53e385097f74",
				"e21fb416-a9a3-4a16-9347-a28cc65076a0",
				"e03c511c-56dc-41cc-a4b5-082e461c83ea",
				"c7a297df-3d62-4f67-a994-ed86ac440053",
				"ed80ba2a-1b70-40b1-b14c-ff63797dd58e",
				"58b19781-d512-43bd-ac1e-d7b0050eedaa",
				"7ce3e515-e83e-4a9a-8ea5-f1017600c71f",
				"58693bd4-24d7-4fb7-9b26-a73849241933",
				"f0b18149-1909-4534-9a3b-9a75f72304ee",
				"059efe62-9a56-414b-bc8e-65caf03f12e4",
				"d590d095-3588-4f16-82af-4d79651b1a86",
				"e9b42238-9c12-4a79-b2c8-1e426373c008",
				"fef6c0f2-281d-4b62-9031-317a778426c9",
				"84b266ed-d64d-49f8-bb86-c6f9cc4cf45a",
				"40086b51-ac58-47c7-9f14-684346af9012",
			},
		},
		{
			description: "searches transaction categories",
			queryParams: map[string]string{
				handler.QueryParamSearch: "Emprestimos",
			},
			token:        mockoauth.PremiumTierMockToken,
			expectedCode: http.StatusOK,
			expectedTransactionCategoryIDs: []string{
				"470850a0-781b-4892-8e96-5590d16c9fe1",
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

			var out dto.ListTransactionCategoriesResponse
			statusCode, rawBody, err := app.MakeRequest(
				http.MethodGet,
				"/api/v1/transactions/categories",
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

			if len(test.expectedTransactionCategoryIDs) == 0 {
				assert.Empty(t, out.Items)
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
