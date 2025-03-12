package restapi

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/danielmesquitta/api-finance-manager/internal/app/restapi/dto"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/usecase"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/db/sqlc"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/oauth/mockoauth"
	"github.com/google/uuid"
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
			description: "List transactions",
			queryParams: map[string]string{
				"page_size": "10",
			},
			token:        mockoauth.DefaultMockToken,
			expectedCode: http.StatusOK,
			expectedTransactionIDs: []string{
				"7e004b7e-bbc8-4030-af19-bba0354582c4",
				"6e238d84-2e41-477b-bce2-713272b006ee",
				"aff780e3-5140-4e4b-a3e3-9f407c07f165",
				"4ec1fab1-fc1c-42d4-a5e7-419af725e14e",
				"8cf37ab8-06fd-4756-897f-1c0be9a4cb32",
				"ee7560a7-5f74-4537-a6a9-f2a6068cdc58",
				"a9038c41-c458-43e7-8c2f-abec2cce5f15",
				"ed0f5518-8580-4da1-bdb7-baa223f611eb",
				"a65d5ecc-eb91-4659-a5b9-2862080d17c8",
				"ecbe6fee-6de7-422c-9863-8cbe788ede19",
			},
		},
		{
			description: "Filter transactions by payment method id",
			queryParams: map[string]string{
				"payment_method_ids": "5d140153-c072-42ce-b19c-c5c9b528dba4",
				"page_size":          "5",
			},
			token:        mockoauth.DefaultMockToken,
			expectedCode: http.StatusOK,
			expectedTransactionIDs: []string{
				"366882dd-a154-4206-9028-fd098c9523e2",
				"d083729b-c99d-46b6-a794-0264f6da9a56",
				"03fe4c66-631d-4c68-a1f6-b19fb0ea79a3",
				"71c230a9-e588-404c-b4c1-99c9cfb29d34",
				"8af66ada-4b5d-4148-8edc-9118c7550743",
			},
		},
		{
			description: "Filter transactions by category id",
			queryParams: map[string]string{
				"category_ids": "059efe62-9a56-414b-bc8e-65caf03f12e4,ed80ba2a-1b70-40b1-b14c-ff63797dd58e",
			},
			token:        mockoauth.DefaultMockToken,
			expectedCode: http.StatusOK,
			expectedTransactionIDs: []string{
				"1ca742b6-a68e-4f0f-850e-a7c33874e0d5",
				"9a8d60e2-6cef-4867-ac92-550780e2260e",
				"5fb7bf3b-f472-4912-b692-4195dc4abd32",
				"71cb5cc5-5510-4502-849d-18ee132042b3",
				"134cb489-e1d6-4dec-822e-265598d90bb6",
			},
		},
		{
			description: "Filter transactions by institution id",
			queryParams: map[string]string{
				"institution_ids": "88f812ab-9bc9-4830-afc6-7ac0ba67b1ec",
				"page_size":       "5",
			},
			token:        mockoauth.DefaultMockToken,
			expectedCode: http.StatusOK,
			expectedTransactionIDs: []string{
				"7e004b7e-bbc8-4030-af19-bba0354582c4",
				"aff780e3-5140-4e4b-a3e3-9f407c07f165",
				"6e238d84-2e41-477b-bce2-713272b006ee",
				"4ec1fab1-fc1c-42d4-a5e7-419af725e14e",
				"8cf37ab8-06fd-4756-897f-1c0be9a4cb32",
			},
		},
		{
			description: "Filter transactions by is expense",
			queryParams: map[string]string{
				"is_expense": "TRUE",
				"page_size":  "5",
			},
			token:        mockoauth.DefaultMockToken,
			expectedCode: http.StatusOK,
			expectedTransactionIDs: []string{
				"db6d71ef-d8e0-4249-b6bb-33275015a4cc",
				"c48bacbf-8587-4ce1-8d75-98d0684ce48b",
				"70e62f2a-83c0-4f56-afb5-ef4c9b75fb6e",
				"2f1a74ef-159b-43dd-9d49-f4baa2d2eae8",
				"a2c13da7-3fe4-4347-8f30-371488d403ed",
			},
		},
		{
			description: "Filter transactions by is income",
			queryParams: map[string]string{
				"is_income": "TRUE",
				"page_size": "5",
			},
			token:        mockoauth.DefaultMockToken,
			expectedCode: http.StatusOK,
			expectedTransactionIDs: []string{
				"7e004b7e-bbc8-4030-af19-bba0354582c4",
				"6e238d84-2e41-477b-bce2-713272b006ee",
				"aff780e3-5140-4e4b-a3e3-9f407c07f165",
				"4ec1fab1-fc1c-42d4-a5e7-419af725e14e",
				"8cf37ab8-06fd-4756-897f-1c0be9a4cb32",
			},
		},
		{
			description: "Filter transactions by is ignored",
			queryParams: map[string]string{
				"is_ignored": "TRUE",
				"page_size":  "5",
			},
			token:        mockoauth.DefaultMockToken,
			expectedCode: http.StatusOK,
			expectedTransactionIDs: []string{
				"422357cb-655a-46db-903e-46af7206ae10",
				"cf52af82-aa64-4568-85e1-6860d652e0a7",
				"7316fedd-490e-4895-ac0e-915d97eaf49e",
				"4c2ffde0-29c7-4cfe-898f-b9ac00baf912",
				"bc459c69-d49e-4cee-8afb-12260ac32185",
			},
		},
		{
			description: "Filter transactions by date period",
			queryParams: map[string]string{
				"start_date": "2024-11-29T00:00:00.000-03:00",
				"end_date":   "2024-11-30T23:59:59.999-03:00",
				"page_size":  "5",
			},
			token:        mockoauth.DefaultMockToken,
			expectedCode: http.StatusOK,
			expectedTransactionIDs: []string{
				"9a8d60e2-6cef-4867-ac92-550780e2260e",
				"15e0db20-10a8-4cc1-984e-0c8d07b9300c",
				"58ee3351-cd16-49ba-bc0a-3be69865f0cf",
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
				rawBody,
			)

			if len(test.expectedTransactionIDs) == 0 {
				assert.Empty(t, out.Items)
				return
			}

			assert.Len(
				t,
				out.Items,
				len(test.expectedTransactionIDs),
			)

			transactionIDs := make([]string, len(out.Items))
			for i, transaction := range out.Items {
				transactionIDs[i] = transaction.ID.String()
			}

			assert.ElementsMatch(
				t,
				test.expectedTransactionIDs,
				transactionIDs,
			)
		})
	}
}

func TestGetTransactionRoute(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description           string
		token                 string
		transactionID         string
		expectedCode          int
		expectedTransactionID string
	}{
		{
			description:   "Fail to list transactions without token",
			token:         "",
			transactionID: "7e004b7e-bbc8-4030-af19-bba0354582c4",
			expectedCode:  http.StatusBadRequest,
		},
		{
			description:           "Get transactions",
			token:                 mockoauth.DefaultMockToken,
			transactionID:         "7e004b7e-bbc8-4030-af19-bba0354582c4",
			expectedCode:          http.StatusOK,
			expectedTransactionID: "7e004b7e-bbc8-4030-af19-bba0354582c4",
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

			var out dto.GetTransactionResponse
			statusCode, rawBody, err := app.MakeRequest(
				http.MethodGet,
				"/api/v1/transactions/"+test.transactionID,
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

			if test.expectedTransactionID != "" {
				assert.Equal(
					t,
					test.expectedTransactionID,
					out.Transaction.ID.String(),
				)
			}
		})
	}
}

func TestUpdateTransactionRoute(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description   string
		token         string
		transactionID string
		body          dto.UpdateTransactionRequest
		expectedCode  int
	}{
		{
			description:   "Fail to update transaction without token",
			token:         "",
			transactionID: "7e004b7e-bbc8-4030-af19-bba0354582c4",
			expectedCode:  http.StatusBadRequest,
		},
		{
			description:   "Update transaction",
			token:         mockoauth.DefaultMockToken,
			transactionID: "7e004b7e-bbc8-4030-af19-bba0354582c4",
			expectedCode:  http.StatusNoContent,
			body: dto.UpdateTransactionRequest{
				UpdateTransactionInput: usecase.UpdateTransactionInput{
					Name:   "Foo bar",
					Amount: 5436,
					Date:   time.Now(),
					PaymentMethodID: uuid.MustParse(
						"fc7adfa0-259c-430e-99f5-bef5281add10",
					),
					AccountID: uuid.MustParse(
						"ac4d82a0-9eff-4936-8a2e-8d12591c9d00",
					),
					InstitutionID: uuid.MustParse(
						"df5dbd97-89c7-4776-8b3f-7992bc2bb16b",
					),
					CategoryID: uuid.MustParse(
						"059efe62-9a56-414b-bc8e-65caf03f12e4",
					),
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
				http.MethodPut,
				"/api/v1/transactions/"+test.transactionID,
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

			if test.expectedCode != http.StatusNoContent {
				return
			}

			transaction, err := app.db.GetTransaction(
				ctx,
				sqlc.GetTransactionParams{
					ID:     uuid.MustParse(test.transactionID),
					UserID: user.ID,
				},
			)
			assert.Nil(t, err)

			expectedTransaction := map[string]any{
				"Name":            test.body.Name,
				"Amount":          test.body.Amount,
				"Date":            test.body.Date.Format(time.RFC3339),
				"PaymentMethodID": test.body.PaymentMethodID,
				"CategoryID":      test.body.CategoryID,
				"AccountID":       test.body.AccountID,
				"InstitutionID":   test.body.InstitutionID,
			}

			actualTransaction := map[string]any{
				"Name":            transaction.Name,
				"Amount":          transaction.Amount,
				"Date":            transaction.Date.Format(time.RFC3339),
				"PaymentMethodID": transaction.PaymentMethodID,
				"CategoryID":      transaction.CategoryID,
				"AccountID":       *transaction.AccountID,
				"InstitutionID":   *transaction.InstitutionID,
			}

			assert.Equal(t, expectedTransaction, actualTransaction)
		})
	}
}

func TestCreateTransactionRoute(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description  string
		token        string
		body         dto.CreateTransactionRequest
		expectedCode int
	}{
		{
			description:  "Fail to update transaction without token",
			token:        "",
			expectedCode: http.StatusBadRequest,
		},
		{
			description:  "Create transaction",
			token:        mockoauth.DefaultMockToken,
			expectedCode: http.StatusCreated,
			body: func() dto.CreateTransactionRequest {
				categoryID := uuid.MustParse(
					"373b150b-94bd-44b2-abdd-2aab14e74fad",
				)
				return dto.CreateTransactionRequest{
					CreateTransactionInput: usecase.CreateTransactionInput{
						Name:   "Foo bar",
						Amount: 5436,
						PaymentMethodID: uuid.MustParse(
							"5d140153-c072-42ce-b19c-c5c9b528dba4",
						),
						Date:       time.Now(),
						CategoryID: &categoryID,
					},
				}
			}(),
		},
		{
			description:  "Create transaction without category",
			token:        mockoauth.DefaultMockToken,
			expectedCode: http.StatusCreated,
			body: dto.CreateTransactionRequest{
				CreateTransactionInput: usecase.CreateTransactionInput{
					Name:   "Foo bar",
					Amount: -6543,
					PaymentMethodID: uuid.MustParse(
						"fc7adfa0-259c-430e-99f5-bef5281add10",
					),
					Date: time.Now(),
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
				"/api/v1/transactions",
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

			transaction, err := app.db.GetLatestTransactionByUserID(
				ctx,
				user.ID.String(),
			)
			assert.Nil(t, err)

			defaultCategoryID := uuid.MustParse(
				"40086b51-ac58-47c7-9f14-684346af9012",
			)

			expectedCategoryID := test.body.CategoryID
			if expectedCategoryID == nil {
				expectedCategoryID = &defaultCategoryID
			}

			expectedTransaction := map[string]any{
				"Name":            test.body.Name,
				"Amount":          test.body.Amount,
				"PaymentMethodID": test.body.PaymentMethodID,
				"Date":            test.body.Date.Format(time.RFC3339),
				"CategoryID":      *expectedCategoryID,
			}

			actualTransaction := map[string]any{
				"Name":            transaction.Name,
				"Amount":          transaction.Amount,
				"PaymentMethodID": transaction.PaymentMethodID,
				"Date":            transaction.Date.Format(time.RFC3339),
				"CategoryID":      transaction.CategoryID,
			}

			assert.Equal(t, expectedTransaction, actualTransaction)
		})
	}
}
