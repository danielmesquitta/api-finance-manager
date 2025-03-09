package restapi

import (
	"context"
	"net/http"
	"testing"

	"github.com/danielmesquitta/api-finance-manager/internal/app/restapi/dto"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/oauth/mockoauth"
	"github.com/stretchr/testify/assert"
)

func TestListTransactionsRoute(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description            string
		queryParams            map[string]string
		token                  string
		expectedCode           int
		expectedTransactionIDs []string
	}{
		{
			description:            "Fail to list transactions without token",
			queryParams:            map[string]string{},
			token:                  "",
			expectedCode:           http.StatusBadRequest,
			expectedTransactionIDs: []string{},
		},
		{
			description:  "List all transactions",
			queryParams:  map[string]string{},
			token:        mockoauth.DefaultMockToken,
			expectedCode: http.StatusOK,
			expectedTransactionIDs: []string{
				"61edb9dd-c137-4e7e-8361-bc78a7ef864b",
				"571c215a-1ee2-4b1a-a316-3ffdd971340d",
				"26657ab2-19cc-47a0-8af6-160f12737e14",
				"a5707214-415e-4c0c-b8da-e1e225365151",
				"f274a4be-1150-4542-896d-88239378b828",
				"8df11353-b6ec-42c5-9fec-84ae140d85cb",
				"27d30d16-c585-49da-8370-bdd77c278295",
				"f68cadfa-b54c-4e37-857c-51db6bb0c465",
				"eb3c0fc8-77bd-4130-83b6-af815d1a2956",
				"79260d65-66bb-476e-85db-1fce518b6aae",
				"cad1e583-f48c-460f-8a46-a3a86abbb2fa",
				"18d326f3-13f2-43c3-ab33-920bc9caefb2",
				"3319a062-50cd-4ea3-afbe-8edc18b21686",
				"f8309ff6-f457-485e-abd5-6c8df4f20ceb",
				"37204747-eedf-4407-8618-ce0e24c9a36a",
			},
		},
		{
			description: "Filter transactions by payment method id",
			queryParams: map[string]string{
				"payment_method_ids": "2158b0b6-844f-44b6-b487-282d0c1b045c",
			},
			token:        mockoauth.DefaultMockToken,
			expectedCode: http.StatusOK,
			expectedTransactionIDs: []string{
				"61edb9dd-c137-4e7e-8361-bc78a7ef864b",
				"a5707214-415e-4c0c-b8da-e1e225365151",
				"27d30d16-c585-49da-8370-bdd77c278295",
				"f68cadfa-b54c-4e37-857c-51db6bb0c465",
				"79260d65-66bb-476e-85db-1fce518b6aae",
				"18d326f3-13f2-43c3-ab33-920bc9caefb2",
				"f8309ff6-f457-485e-abd5-6c8df4f20ceb",
			},
		},
		{
			description: "Filter transactions by category id",
			queryParams: map[string]string{
				"category_ids": "02701aac-b8db-4c7e-834c-6d4f4eab3399,03bd0abc-7186-4eb3-9871-e4f624c535b8",
			},
			token:        mockoauth.DefaultMockToken,
			expectedCode: http.StatusOK,
			expectedTransactionIDs: []string{
				"61edb9dd-c137-4e7e-8361-bc78a7ef864b",
				"26657ab2-19cc-47a0-8af6-160f12737e14",
			},
		},
		{
			description: "Filter transactions by institution id",
			queryParams: map[string]string{
				"institution_ids": "1202269c-ed03-4dfe-bbcd-c61d615a17b5",
			},
			token:        mockoauth.DefaultMockToken,
			expectedCode: http.StatusOK,
			expectedTransactionIDs: []string{
				"27d30d16-c585-49da-8370-bdd77c278295",
			},
		},
		{
			description: "Filter transactions by is expense",
			queryParams: map[string]string{
				"is_expense": "TRUE",
			},
			token:        mockoauth.DefaultMockToken,
			expectedCode: http.StatusOK,
			expectedTransactionIDs: []string{
				"61edb9dd-c137-4e7e-8361-bc78a7ef864b",
				"a5707214-415e-4c0c-b8da-e1e225365151",
				"f68cadfa-b54c-4e37-857c-51db6bb0c465",
				"eb3c0fc8-77bd-4130-83b6-af815d1a2956",
				"79260d65-66bb-476e-85db-1fce518b6aae",
				"18d326f3-13f2-43c3-ab33-920bc9caefb2",
			},
		},
		{
			description: "Filter transactions by is income",
			queryParams: map[string]string{
				"is_income": "TRUE",
			},
			token:        mockoauth.DefaultMockToken,
			expectedCode: http.StatusOK,
			expectedTransactionIDs: []string{
				"571c215a-1ee2-4b1a-a316-3ffdd971340d",
				"26657ab2-19cc-47a0-8af6-160f12737e14",
				"f274a4be-1150-4542-896d-88239378b828",
				"8df11353-b6ec-42c5-9fec-84ae140d85cb",
				"27d30d16-c585-49da-8370-bdd77c278295",
				"cad1e583-f48c-460f-8a46-a3a86abbb2fa",
				"3319a062-50cd-4ea3-afbe-8edc18b21686",
				"f8309ff6-f457-485e-abd5-6c8df4f20ceb",
				"37204747-eedf-4407-8618-ce0e24c9a36a",
			},
		},
		{
			description: "Filter transactions by is ignored",
			queryParams: map[string]string{
				"is_ignored": "TRUE",
			},
			token:        mockoauth.DefaultMockToken,
			expectedCode: http.StatusOK,
			expectedTransactionIDs: []string{
				"79260d65-66bb-476e-85db-1fce518b6aae",
			},
		},
		{
			description: "Filter transactions by date period",
			queryParams: map[string]string{
				"start_date": "2024-03-09T03:18:28.211Z",
				"end_date":   "2024-06-11T23:18:28.211Z",
			},
			token:        mockoauth.DefaultMockToken,
			expectedCode: http.StatusOK,
			expectedTransactionIDs: []string{
				"a5707214-415e-4c0c-b8da-e1e225365151",
				"26657ab2-19cc-47a0-8af6-160f12737e14",
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
				accessToken = app.SignIn(test.token)
			}

			var out dto.ListTransactionsResponse
			var errRes dto.ErrorResponse
			statusCode, rawBody, err := app.MakeRequest(
				http.MethodGet,
				"/api/v1/transactions",
				WithQueryParams(test.queryParams),
				WithBearerToken(accessToken),
				WithResponse(&out),
				WithError(&errRes),
			)
			assert.Nil(t, err)

			assert.Equal(
				t,
				test.expectedCode,
				statusCode,
			)

			if len(test.expectedTransactionIDs) == 0 {
				assert.Empty(t, out.Items, rawBody)
				return
			}

			assert.Len(
				t,
				out.Items,
				len(test.expectedTransactionIDs),
				rawBody,
			)

			transactionIDs := make([]string, len(out.Items))
			for i, transaction := range out.Items {
				transactionIDs[i] = transaction.ID.String()
			}

			assert.ElementsMatch(
				t,
				test.expectedTransactionIDs,
				transactionIDs,
				rawBody,
			)
		})
	}
}
