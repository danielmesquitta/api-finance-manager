package restapi

import (
	"context"
	"log/slog"
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
		// {
		// 	description:      "Fail to list budgets without token",
		// 	token:            "",
		// 	expectedCode:     http.StatusBadRequest,
		// 	expectedResponse: nil,
		// },
		func() Test {
			dateStr := "2024-11-14T00:00:00.000-03:00"
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
			comparisonStartDate, _ := time.Parse(
				time.RFC3339,
				"2024-10-01T00:00:00-03:00",
			)
			comparisonEndDate, _ := time.Parse(
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
							Amount: 500000,
							Date:   date,
						},
						Spent:                              8807768,
						Available:                          -8307768,
						AvailablePercentageVariation:       -6572,
						AvailablePerDay:                    -276926,
						AvailablePerDayPercentageVariation: 0,
						ComparisonDates: usecase.ComparisonDates{
							StartDate:           startDate,
							EndDate:             endDate,
							ComparisonStartDate: comparisonStartDate,
							ComparisonEndDate:   comparisonEndDate,
						},
						BudgetCategories: []usecase.GetBudgetBudgetCategories{
							{
								BudgetCategory: entity.BudgetCategory{
									ID: uuid.MustParse(
										"5932b6b2-52eb-4e5c-ab00-19de3c534578",
									),
								},
								Spent:     64754,
								Available: 185246,
							},
							{
								BudgetCategory: entity.BudgetCategory{
									ID: uuid.MustParse(
										"39a948f7-0619-4383-a6b3-17fe653651c2",
									),
								},
								Spent:     7319,
								Available: 92681,
							},
							{
								BudgetCategory: entity.BudgetCategory{
									ID: uuid.MustParse(
										"75b30904-600f-4d40-a7cc-f7f7f800679d",
									),
								},
								Spent:     0,
								Available: 150000,
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

			slog.Debug("actualResponse", "actualResponse", actualResponse)

			assert.Equal(
				t,
				test.expectedCode,
				statusCode,
				rawBody,
			)

			if test.expectedResponse == nil {
				return
			}

			// Assert budget properties
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
			assert.Equal(
				t,
				test.expectedResponse.Budget.Date.Format(time.RFC3339),
				actualResponse.Budget.Date.Format(time.RFC3339),
			)

			// Assert budget metrics
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

			// Assert comparison dates
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

			// Assert budget categories
			assert.Equal(
				t,
				len(test.expectedResponse.BudgetCategories),
				len(actualResponse.BudgetCategories),
			)
			for i, expectedCategory := range test.expectedResponse.BudgetCategories {
				if i < len(actualResponse.BudgetCategories) {
					assert.Equal(
						t,
						expectedCategory.BudgetCategory,
						actualResponse.BudgetCategories[i].BudgetCategory,
					)
					assert.Equal(
						t,
						expectedCategory.Spent,
						actualResponse.BudgetCategories[i].Spent,
					)
					assert.Equal(
						t,
						expectedCategory.Available,
						actualResponse.BudgetCategories[i].Available,
					)
					assert.Equal(
						t,
						expectedCategory.Category,
						actualResponse.BudgetCategories[i].Category,
					)
				}
			}
		})
	}
}
