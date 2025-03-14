package restapi

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/danielmesquitta/api-finance-manager/internal/app/restapi/dto"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/usecase"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/oauth/mockoauth"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetBudgetRoute(t *testing.T) {
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
			date, _ := time.Parse(time.RFC3339, dateStr)
			budgetID := uuid.MustParse("8aa317f8-702c-43b1-897b-e24a4285d2d2")

			startDate, _ := time.Parse(
				time.RFC3339,
				"2024-11-01T00:00:00-03:00",
			)
			endDate, _ := time.Parse(
				time.RFC3339,
				"2024-11-30T23:59:59.999999999-03:00",
			)
			cmpStartDate, _ := time.Parse(
				time.RFC3339,
				"2024-10-01T00:00:00-03:00",
			)
			cmpEndDate, _ := time.Parse(
				time.RFC3339,
				"2024-10-31T23:59:59.999999999-03:00",
			)

			return Test{
				description:  "Get budget",
				token:        mockoauth.DefaultMockToken,
				queryParams:  map[string]string{"date": dateStr},
				expectedCode: http.StatusOK,
				expectedResponse: &dto.GetBudgetResponse{
					GetBudgetOutput: usecase.GetBudgetOutput{
						Budget: entity.Budget{
							ID:     budgetID,
							Amount: 2000000,
							Date:   date,
						},
						Spent:                              1837970,
						Available:                          162030,
						AvailablePercentageVariation:       8322,
						AvailablePerDay:                    0,
						AvailablePerDayPercentageVariation: 0,
						ComparisonDates: usecase.ComparisonDates{
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
								Spent:     66630,
								Available: 933370,
							},
							{
								BudgetCategory: entity.BudgetCategory{
									ID: uuid.MustParse(
										"39a948f7-0619-4383-a6b3-17fe653651c2",
									),
								},
								Spent:     122918,
								Available: 277082,
							},
							{
								BudgetCategory: entity.BudgetCategory{
									ID: uuid.MustParse(
										"75b30904-600f-4d40-a7cc-f7f7f800679d",
									),
								},
								Spent:     11826,
								Available: 588174,
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
