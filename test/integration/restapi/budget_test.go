package restapi

import (
	"context"
	"net/http"
	"testing"

	"github.com/danielmesquitta/api-finance-manager/internal/app/restapi/dto"
	"github.com/danielmesquitta/api-finance-manager/internal/app/restapi/handler"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/usecase"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/dateutil"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/db/sqlc"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/oauth/mockoauth"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"github.com/stretchr/testify/assert"
)

func TestGetBudget(t *testing.T) {
	t.Parallel()

	type Test struct {
		description      string
		token            string
		queryParams      map[string]string
		expectedCode     int
		expectedResponse *dto.GetBudgetResponse
	}

	tests := []Test{
		{
			description:      "Fail to list budgets without token",
			token:            "",
			expectedCode:     http.StatusBadRequest,
			expectedResponse: nil,
		},
		func() Test {
			dateStr := "2024-11-01T00:00:00-03:00"
			date := dateutil.MustParseISOString(dateStr)
			budgetID := uuid.MustParse("8aa317f8-702c-43b1-897b-e24a4285d2d2")

			startDate := dateutil.MustParseISOString(
				"2024-11-01T00:00:00-03:00",
			)
			endDate := dateutil.MustParseISOString(
				"2024-11-30T23:59:59.999999999-03:00",
			)
			cmpStartDate := dateutil.MustParseISOString(
				"2024-10-01T00:00:00-03:00",
			)
			cmpEndDate := dateutil.MustParseISOString(
				"2024-10-31T23:59:59.999999999-03:00",
			)

			return Test{
				description: "Get budget",
				token:       mockoauth.PremiumTierMockToken,
				queryParams: map[string]string{
					handler.QueryParamDate: dateStr,
				},
				expectedCode: http.StatusOK,
				expectedResponse: &dto.GetBudgetResponse{
					GetBudgetOutput: usecase.GetBudgetOutput{
						Budget: entity.Budget{
							ID:     budgetID,
							Amount: 20_000_00,
							Date:   date,
						},
						Spent:                              18_379_70,
						Available:                          1_620_30,
						AvailablePercentageVariation:       -83_22,
						AvailablePerDay:                    0,
						AvailablePerDayPercentageVariation: 0,
						ComparisonDates: dateutil.ComparisonDates{
							StartDate:           startDate,
							EndDate:             endDate,
							ComparisonStartDate: cmpStartDate,
							ComparisonEndDate:   cmpEndDate,
						},
						BudgetCategories: []usecase.GetBudgetBudgetCategories{
							{
								BudgetCategory: entity.BudgetCategory{
									ID: uuid.MustParse(
										"5932b6b2-52eb-4e5c-ab00-19de3c534578",
									),
								},
								Spent:     666_30,
								Available: 9_333_70,
							},
							{
								BudgetCategory: entity.BudgetCategory{
									ID: uuid.MustParse(
										"39a948f7-0619-4383-a6b3-17fe653651c2",
									),
								},
								Spent:     1_229_18,
								Available: 2_770_82,
							},
							{
								BudgetCategory: entity.BudgetCategory{
									ID: uuid.MustParse(
										"75b30904-600f-4d40-a7cc-f7f7f800679d",
									),
								},
								Spent:     118_26,
								Available: 5_881_74,
							},
						},
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

			accessToken := ""
			if test.token != "" {
				signInRes := app.SignIn(test.token)
				accessToken = signInRes.AccessToken
			}

			var actualResponse dto.GetBudgetResponse
			statusCode, rawBody, err := app.MakeRequest(
				http.MethodGet,
				"/api/v1/budgets",
				WithBearerToken(accessToken),
				WithQueryParams(test.queryParams),
				WithResponse(&actualResponse),
			)
			assert.Nil(t, err)

			assert.Equal(
				t,
				test.expectedCode,
				statusCode,
				rawBody,
			)

			if test.expectedResponse == nil {
				return
			}

			assert.Equal(
				t,
				test.expectedResponse.Budget.ID,
				actualResponse.Budget.ID,
			)
			assert.Equal(
				t,
				test.expectedResponse.Budget.Amount,
				actualResponse.Budget.Amount,
			)
			assert.True(
				t,
				test.expectedResponse.Budget.Date.After(
					actualResponse.Budget.Date,
				),
			)

			assert.Equal(t, test.expectedResponse.Spent, actualResponse.Spent)
			assert.Equal(
				t,
				test.expectedResponse.Available,
				actualResponse.Available,
			)
			assert.Equal(
				t,
				test.expectedResponse.AvailablePercentageVariation,
				actualResponse.AvailablePercentageVariation,
			)
			assert.Equal(
				t,
				test.expectedResponse.AvailablePerDay,
				actualResponse.AvailablePerDay,
			)
			assert.Equal(
				t,
				test.expectedResponse.AvailablePerDayPercentageVariation,
				actualResponse.AvailablePerDayPercentageVariation,
			)

			assert.Equal(
				t,
				test.expectedResponse.ComparisonDates.StartDate,
				actualResponse.ComparisonDates.StartDate,
			)
			assert.Equal(
				t,
				test.expectedResponse.ComparisonDates.EndDate,
				actualResponse.ComparisonDates.EndDate,
			)
			assert.Equal(
				t,
				test.expectedResponse.ComparisonDates.ComparisonStartDate,
				actualResponse.ComparisonDates.ComparisonStartDate,
			)
			assert.Equal(
				t,
				test.expectedResponse.ComparisonDates.ComparisonEndDate,
				actualResponse.ComparisonDates.ComparisonEndDate,
			)

			assert.Equal(
				t,
				len(test.expectedResponse.BudgetCategories),
				len(actualResponse.BudgetCategories),
			)

			actualBudgetCategories := map[uuid.UUID]usecase.GetBudgetBudgetCategories{}
			for _, budgetCategory := range actualResponse.BudgetCategories {
				actualBudgetCategories[budgetCategory.BudgetCategory.ID] = budgetCategory
			}

			for _, expectedCategory := range test.expectedResponse.BudgetCategories {
				actualBudgetCategory, ok := actualBudgetCategories[expectedCategory.BudgetCategory.ID]
				assert.True(t, ok, expectedCategory.BudgetCategory.ID)
				if !ok {
					continue
				}

				assert.Equal(
					t,
					expectedCategory.Spent,
					actualBudgetCategory.Spent,
				)
				assert.Equal(
					t,
					expectedCategory.Available,
					actualBudgetCategory.Available,
				)
			}
		})
	}
}

func TestUpsertBudget(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description  string
		token        string
		body         dto.UpsertBudgetRequest
		expectedCode int
	}{
		{
			description:  "Fail to update budget without token",
			token:        "",
			expectedCode: http.StatusBadRequest,
		},
		{
			description:  "Create budget",
			token:        mockoauth.PremiumTierMockToken,
			expectedCode: http.StatusNoContent,
			body: dto.UpsertBudgetRequest{
				UpsertBudgetInput: usecase.UpsertBudgetInput{
					Amount: 15_000_00,
					Date:   "2025-03-14T00:00:00-03:00",
					Categories: []usecase.UpsertBudgetCategoryInput{
						{
							Amount: 5_000_00,
							CategoryID: uuid.MustParse(
								"059efe62-9a56-414b-bc8e-65caf03f12e4",
							),
						},
						{
							Amount: 1_000_00,
							CategoryID: uuid.MustParse(
								"84b266ed-d64d-49f8-bb86-c6f9cc4cf45a",
							),
						},
						{
							Amount: 9_000_00,
							CategoryID: uuid.MustParse(
								"ed80ba2a-1b70-40b1-b14c-ff63797dd58e",
							),
						},
					},
				},
			},
		},
		{
			description:  "Update budget",
			token:        mockoauth.PremiumTierMockToken,
			expectedCode: http.StatusNoContent,
			body: dto.UpsertBudgetRequest{
				UpsertBudgetInput: usecase.UpsertBudgetInput{
					Amount: 15_000_00,
					Date:   "2024-10-15T00:00:00.000-03:00",
					Categories: []usecase.UpsertBudgetCategoryInput{
						{
							Amount: 5_000_00,
							CategoryID: uuid.MustParse(
								"059efe62-9a56-414b-bc8e-65caf03f12e4",
							),
						},
						{
							Amount: 1_000_00,
							CategoryID: uuid.MustParse(
								"84b266ed-d64d-49f8-bb86-c6f9cc4cf45a",
							),
						},
						{
							Amount: 9_000_00,
							CategoryID: uuid.MustParse(
								"ed80ba2a-1b70-40b1-b14c-ff63797dd58e",
							),
						},
					},
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
				"/api/v1/budgets",
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

			parsedDate := dateutil.MustParseISOString(test.body.Date)

			actualBudget, err := app.db.GetBudget(
				ctx,
				sqlc.GetBudgetParams{
					UserID: user.ID,
					Date:   parsedDate,
				},
			)
			assert.Nil(t, err)
			assert.Equal(t, test.body.Amount, actualBudget.Amount)
			assert.Equal(
				t,
				dateutil.ToMonthStart(parsedDate),
				actualBudget.Date,
			)

			actualBudgetCategories, err := app.db.ListBudgetCategories(
				ctx,
				actualBudget.ID,
			)
			assert.Nil(t, err)

			actualBudgetCategoriesMap := map[uuid.UUID]entity.BudgetCategory{}
			for _, abc := range actualBudgetCategories {
				bc := entity.BudgetCategory{}
				if err := copier.Copy(&bc, abc.BudgetCategory); err != nil {
					t.Fatal(err)
				}
				actualBudgetCategoriesMap[abc.TransactionCategory.ID] = bc
			}

			for _, expectedCategory := range test.body.Categories {
				actualCategory, ok := actualBudgetCategoriesMap[expectedCategory.CategoryID]
				assert.True(t, ok, expectedCategory.CategoryID)
				if !ok {
					continue
				}

				assert.Equal(t, expectedCategory.Amount, actualCategory.Amount)
			}

		})
	}
}
