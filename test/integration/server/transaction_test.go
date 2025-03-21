package server

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/danielmesquitta/api-finance-manager/internal/app/server/dto"
	"github.com/danielmesquitta/api-finance-manager/internal/app/server/handler"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/usecase/transaction"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/oauth/mockoauth"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestListTransactions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description            string
		queryParams            map[string]string
		token                  string
		expectedCode           int
		expectedTransactionIDs []string
	}{
		{
			description:            "fails without token",
			queryParams:            map[string]string{},
			token:                  "",
			expectedCode:           http.StatusBadRequest,
			expectedTransactionIDs: []string{},
		},
		{
			description: "lists transactions",
			queryParams: map[string]string{
				handler.QueryParamPageSize: "10",
			},
			token:        mockoauth.PremiumTierMockToken,
			expectedCode: http.StatusOK,
			expectedTransactionIDs: []string{
				"cac59381-300c-4f63-be76-7a2f654cd480",
				"e827eec2-1fc4-4976-a3ae-86294a0fc338",
				"5a509fd8-a0bf-43d0-b142-60a803b34141",
				"2b5fe271-4da7-4332-9226-170073e07d2e",
				"c7e63e56-b4e6-4423-af2c-bf5a1c529519",
				"59de5384-83fa-4624-bb63-b0187d2b094f",
				"511592ea-1bc6-4b88-afa9-5a9b0174913a",
				"e9af104e-d43c-4d24-a7ab-617cf8efaede",
				"2d706106-bf58-48d4-9553-d203f06364c2",
				"aecb4079-923d-48d1-8d61-09af5c9f0b00",
			},
		},
		{
			description: "filters transactions by payment method id",
			queryParams: map[string]string{
				handler.QueryParamPaymentMethodIDs: "5d140153-c072-42ce-b19c-c5c9b528dba4",
				handler.QueryParamPageSize:         "5",
			},
			token:        mockoauth.PremiumTierMockToken,
			expectedCode: http.StatusOK,
			expectedTransactionIDs: []string{
				"04d434cf-a2b7-4c13-82db-8d3b93b0ca40",
				"aac3289e-e691-4216-a551-be52858a5a5c",
				"cfe030c5-363c-4d56-b9a6-a270689d3f53",
				"c53caf93-0f45-4b84-8d93-0a141ae9a93f",
				"d1a19260-aa4d-4b75-b639-301be5cae12d",
			},
		},
		{
			description: "filters transactions by category id",
			queryParams: map[string]string{
				handler.QueryParamCategoryIDs: "059efe62-9a56-414b-bc8e-65caf03f12e4,ed80ba2a-1b70-40b1-b14c-ff63797dd58e",
			},
			token:        mockoauth.PremiumTierMockToken,
			expectedCode: http.StatusOK,
			expectedTransactionIDs: []string{
				"332df5e6-2624-451b-b95f-437094731851",
				"fc8bdec5-d809-4009-bb8e-22b0a043cb0d",
				"4dfefca2-fcf4-4e14-86ca-caaed8603f0f",
				"3d61f0b6-abd6-4642-ba72-a6eac603fdae",
				"8274fdc6-5a6e-40d0-9baf-10e44b1eebd1",
			},
		},
		{
			description: "filters transactions by institution id",
			queryParams: map[string]string{
				handler.QueryParamInstitutionIDs: "88f812ab-9bc9-4830-afc6-7ac0ba67b1ec",
				handler.QueryParamPageSize:       "5",
			},
			token:        mockoauth.PremiumTierMockToken,
			expectedCode: http.StatusOK,
			expectedTransactionIDs: []string{
				"cac59381-300c-4f63-be76-7a2f654cd480",
				"e827eec2-1fc4-4976-a3ae-86294a0fc338",
				"5a509fd8-a0bf-43d0-b142-60a803b34141",
				"2b5fe271-4da7-4332-9226-170073e07d2e",
				"c7e63e56-b4e6-4423-af2c-bf5a1c529519",
			},
		},
		{
			description: "filters transactions by is expense",
			queryParams: map[string]string{
				handler.QueryParamIsExpense: "TRUE",
				handler.QueryParamPageSize:  "5",
			},
			token:        mockoauth.PremiumTierMockToken,
			expectedCode: http.StatusOK,
			expectedTransactionIDs: []string{
				"cac59381-300c-4f63-be76-7a2f654cd480",
				"e827eec2-1fc4-4976-a3ae-86294a0fc338",
				"5a509fd8-a0bf-43d0-b142-60a803b34141",
				"2b5fe271-4da7-4332-9226-170073e07d2e",
				"c7e63e56-b4e6-4423-af2c-bf5a1c529519",
			},
		},
		{
			description: "filters transactions by is income",
			queryParams: map[string]string{
				handler.QueryParamIsIncome: "TRUE",
				handler.QueryParamPageSize: "5",
			},
			token:        mockoauth.PremiumTierMockToken,
			expectedCode: http.StatusOK,
			expectedTransactionIDs: []string{
				"d42f49d6-6652-4268-8970-f4293eb63c03",
				"e6cb4beb-b431-42ab-b654-94748247ba93",
				"4cd343af-9822-4018-9636-03b465176485",
				"aac3289e-e691-4216-a551-be52858a5a5c",
				"cfe030c5-363c-4d56-b9a6-a270689d3f53",
			},
		},
		{
			description: "filters transactions by is ignored",
			queryParams: map[string]string{
				handler.QueryParamIsIgnored: "TRUE",
				handler.QueryParamPageSize:  "5",
			},
			token:        mockoauth.PremiumTierMockToken,
			expectedCode: http.StatusOK,
			expectedTransactionIDs: []string{
				"c0c71a23-41f5-4688-9441-544eaf2bdc76",
				"b3047f7e-5636-4862-bca5-1f9dbc40d70f",
				"49508dc4-ab8f-46da-b6c5-8431a1dadae0",
				"0bbd955a-062b-4b7c-82cb-3e86a8c8b714",
				"c0047189-4395-4270-88c4-a5a443bd6350",
			},
		},
		{
			description: "filters transactions by date period",
			queryParams: map[string]string{
				handler.QueryParamStartDate: "2024-11-29T00:00:00.000-03:00",
				handler.QueryParamEndDate:   "2024-11-30T23:59:59.999-03:00",
				handler.QueryParamPageSize:  "5",
			},
			token:        mockoauth.PremiumTierMockToken,
			expectedCode: http.StatusOK,
			expectedTransactionIDs: []string{
				"68c59c47-1359-4f04-89bf-8f97f4b0b8f9",
				"fc8bdec5-d809-4009-bb8e-22b0a043cb0d",
				"e1c73c22-7d52-43e2-80a8-63ce6da99e53",
			},
		},
		{
			description: "search transactions by institution name",
			queryParams: map[string]string{
				handler.QueryParamSearch:   "nubank",
				handler.QueryParamPageSize: "5",
			},
			token:        mockoauth.PremiumTierMockToken,
			expectedCode: http.StatusOK,
			expectedTransactionIDs: []string{
				"cac59381-300c-4f63-be76-7a2f654cd480",
				"e827eec2-1fc4-4976-a3ae-86294a0fc338",
				"5a509fd8-a0bf-43d0-b142-60a803b34141",
				"2b5fe271-4da7-4332-9226-170073e07d2e",
				"c7e63e56-b4e6-4423-af2c-bf5a1c529519",
			},
		},
		{
			description: "search transactions by payment method name",
			queryParams: map[string]string{
				handler.QueryParamSearch:   "pix",
				handler.QueryParamPageSize: "5",
			},
			token:        mockoauth.PremiumTierMockToken,
			expectedCode: http.StatusOK,
			expectedTransactionIDs: []string{
				"04d434cf-a2b7-4c13-82db-8d3b93b0ca40",
				"aac3289e-e691-4216-a551-be52858a5a5c",
				"cfe030c5-363c-4d56-b9a6-a270689d3f53",
				"c53caf93-0f45-4b84-8d93-0a141ae9a93f",
				"d1a19260-aa4d-4b75-b639-301be5cae12d",
			},
		},
		{
			description: "search transactions by category name",
			queryParams: map[string]string{
				handler.QueryParamSearch:   "transprt",
				handler.QueryParamPageSize: "5",
			},
			token:        mockoauth.PremiumTierMockToken,
			expectedCode: http.StatusOK,
			expectedTransactionIDs: []string{
				"33f86bd0-cbb9-40e9-bdb9-c7a5d0ab0e5f",
				"351093fb-d63d-406f-a8e0-ddd8a37b9bcb",
				"dcbaccad-1c50-4f68-91d7-9645671fc66e",
				"7cad86b1-5579-41b0-b927-0d6d6ccc7669",
				"7ed055ed-fbff-492c-9767-11e07d60cc5c",
			},
		},
		{
			description: "search transactions by transaction name",
			queryParams: map[string]string{
				handler.QueryParamSearch:   "daniel mesquita",
				handler.QueryParamPageSize: "5",
			},
			token:        mockoauth.PremiumTierMockToken,
			expectedCode: http.StatusOK,
			expectedTransactionIDs: []string{
				"7f00a7c6-2f07-4e19-9be6-ab6be20e2bdd",
				"155fdf1b-9199-44d9-8fb2-2f7ebbd06440",
				"312cf698-f131-4f14-8fd9-b0ffda75bfa2",
				"ec793a7d-37ba-47e3-b12a-04499c08b762",
				"96463c37-4b69-4dd5-9ad3-136e86aae2bf",
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

			var out dto.ListTransactionsResponse
			var errRes dto.ErrorResponse
			statusCode, rawBody, err := app.MakeRequest(
				http.MethodGet,
				"/api/v1/transactions",
				WithQueryParams(test.queryParams),
				WithBearerToken(signInRes.AccessToken),
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

func TestGetTransaction(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description           string
		token                 string
		transactionID         string
		expectedCode          int
		expectedTransactionID string
	}{
		{
			description:   "fails without token",
			token:         "",
			transactionID: "e1c73c22-7d52-43e2-80a8-63ce6da99e53",
			expectedCode:  http.StatusBadRequest,
		},
		{
			description:           "Get transactions",
			token:                 mockoauth.PremiumTierMockToken,
			transactionID:         "e1c73c22-7d52-43e2-80a8-63ce6da99e53",
			expectedCode:          http.StatusOK,
			expectedTransactionID: "e1c73c22-7d52-43e2-80a8-63ce6da99e53",
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

			var out dto.GetTransactionResponse
			statusCode, rawBody, err := app.MakeRequest(
				http.MethodGet,
				"/api/v1/transactions/"+test.transactionID,
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

func TestUpdateTransaction(t *testing.T) {
	t.Parallel()

	type Test struct {
		description         string
		token               string
		transactionID       string
		body                dto.UpdateTransactionRequest
		expectedCode        int
		expectedTransaction *entity.Transaction
	}

	tests := []Test{
		{
			description:   "fails without token",
			token:         "",
			transactionID: "e1c73c22-7d52-43e2-80a8-63ce6da99e53",
			expectedCode:  http.StatusBadRequest,
		},
		{
			description:   "fails with non-existing transaction",
			token:         mockoauth.PremiumTierMockToken,
			transactionID: "5fde4a75-f4df-415e-86bb-d7e24d488e36",
			expectedCode:  http.StatusNotFound,
			body: dto.UpdateTransactionRequest{
				UpdateTransactionUseCaseInput: transaction.UpdateTransactionUseCaseInput{
					Name: "Lorem ipsum",
				},
			},
		},
		func() Test {
			name := "Foo bar"
			date := time.Now()
			amount := int64(5436)
			paymentMethodID := uuid.MustParse(
				"fc7adfa0-259c-430e-99f5-bef5281add10",
			)
			accountID := uuid.MustParse(
				"ac4d82a0-9eff-4936-8a2e-8d12591c9d00",
			)
			institutionID := uuid.MustParse(
				"df5dbd97-89c7-4776-8b3f-7992bc2bb16b",
			)
			categoryID := uuid.MustParse(
				"2a226707-1f75-4276-9697-3e7aac3c7db6",
			)
			return Test{
				description:   "updates full transaction",
				token:         mockoauth.PremiumTierMockToken,
				transactionID: "e1c73c22-7d52-43e2-80a8-63ce6da99e53",
				expectedCode:  http.StatusNoContent,
				body: dto.UpdateTransactionRequest{
					UpdateTransactionUseCaseInput: transaction.UpdateTransactionUseCaseInput{
						Name:            name,
						Amount:          amount,
						Date:            date,
						PaymentMethodID: paymentMethodID,
						AccountID:       accountID,
						InstitutionID:   institutionID,
						CategoryID:      categoryID,
					},
				},
				expectedTransaction: &entity.Transaction{
					Name:            name,
					Amount:          amount,
					Date:            date,
					PaymentMethodID: paymentMethodID,
					AccountID:       &accountID,
					InstitutionID:   &institutionID,
					CategoryID:      categoryID,
				},
			}
		}(),

		func() Test {
			name := "Foo bar"
			date := time.Now()
			amount := int64(5436)
			paymentMethodID := uuid.MustParse(
				"fc7adfa0-259c-430e-99f5-bef5281add10",
			)
			accountID := uuid.MustParse(
				"e5f31705-cb65-42a5-9072-2b9b59e338a8",
			)
			institutionID := uuid.MustParse(
				"88f812ab-9bc9-4830-afc6-7ac0ba67b1ec",
			)
			categoryID := uuid.MustParse(
				"373b150b-94bd-44b2-abdd-2aab14e74fad",
			)
			return Test{
				description:   "updates partial transaction",
				token:         mockoauth.PremiumTierMockToken,
				transactionID: "e1c73c22-7d52-43e2-80a8-63ce6da99e53",
				expectedCode:  http.StatusNoContent,
				body: dto.UpdateTransactionRequest{
					UpdateTransactionUseCaseInput: transaction.UpdateTransactionUseCaseInput{
						Name:            name,
						Amount:          amount,
						Date:            date,
						PaymentMethodID: paymentMethodID,
					},
				},
				expectedTransaction: &entity.Transaction{
					Name:            name,
					Amount:          amount,
					Date:            date,
					PaymentMethodID: paymentMethodID,
					AccountID:       &accountID,
					InstitutionID:   &institutionID,
					CategoryID:      categoryID,
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
				"/api/v1/transactions/"+test.transactionID,
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

			actualTransaction, err := app.db.GetTransactionByID(
				ctx,
				uuid.MustParse(test.transactionID),
			)
			assert.Nil(t, err)

			assert.Equal(
				t,
				signInRes.User.ID.String(),
				actualTransaction.UserID.String(),
			)
			assert.Equal(
				t,
				test.expectedTransaction.Name,
				actualTransaction.Name,
			)
			assert.Equal(
				t,
				test.expectedTransaction.Amount,
				actualTransaction.Amount,
			)
			assert.Equal(
				t,
				test.expectedTransaction.Date.Format(time.RFC3339),
				actualTransaction.Date.Format(time.RFC3339),
			)
			assert.Equal(
				t,
				test.expectedTransaction.PaymentMethodID.String(),
				actualTransaction.PaymentMethodID.String(),
			)
			assert.Equal(
				t,
				test.expectedTransaction.CategoryID.String(),
				actualTransaction.CategoryID.String(),
			)
			assert.Equal(
				t,
				test.expectedTransaction.AccountID.String(),
				actualTransaction.AccountID.String(),
			)
			assert.Equal(
				t,
				test.expectedTransaction.InstitutionID.String(),
				actualTransaction.InstitutionID.String(),
			)
		})
	}
}

func TestCreateTransaction(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description  string
		token        string
		body         dto.CreateTransactionRequest
		expectedCode int
	}{
		{
			description:  "fails without token",
			token:        "",
			expectedCode: http.StatusBadRequest,
		},
		{
			description:  "creates transaction",
			token:        mockoauth.PremiumTierMockToken,
			expectedCode: http.StatusCreated,
			body: func() dto.CreateTransactionRequest {
				categoryID := uuid.MustParse(
					"373b150b-94bd-44b2-abdd-2aab14e74fad",
				)
				return dto.CreateTransactionRequest{
					CreateTransactionUseCaseInput: transaction.CreateTransactionUseCaseInput{
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
			description:  "creates transaction without category",
			token:        mockoauth.PremiumTierMockToken,
			expectedCode: http.StatusCreated,
			body: dto.CreateTransactionRequest{
				CreateTransactionUseCaseInput: transaction.CreateTransactionUseCaseInput{
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

			signInRes := &dto.SignInResponse{}
			if test.token != "" {
				signInRes = app.SignIn(test.token)
			}

			statusCode, rawBody, err := app.MakeRequest(
				http.MethodPost,
				"/api/v1/transactions",
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

			transaction, err := app.db.GetLatestTransactionByUserID(
				ctx,
				signInRes.User.ID.String(),
			)
			assert.Nil(t, err)

			defaultCategoryID := uuid.MustParse(
				"40086b51-ac58-47c7-9f14-684346af9012",
			)

			expectedCategoryID := test.body.CategoryID
			if expectedCategoryID == nil {
				expectedCategoryID = &defaultCategoryID
			}

			assert.Equal(t, test.body.Name, transaction.Name)
			assert.Equal(t, test.body.Amount, transaction.Amount)
			assert.Equal(
				t,
				test.body.PaymentMethodID.String(),
				transaction.PaymentMethodID.String(),
			)
			assert.Equal(
				t,
				test.body.Date.Format(time.RFC3339),
				transaction.Date.Format(time.RFC3339),
			)
			assert.Equal(
				t,
				expectedCategoryID.String(),
				transaction.CategoryID.String(),
			)
		})
	}
}
