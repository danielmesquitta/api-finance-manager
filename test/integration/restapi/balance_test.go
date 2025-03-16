package restapi

import (
	"context"
	"net/http"
	"testing"

	"github.com/danielmesquitta/api-finance-manager/internal/app/restapi/dto"
	"github.com/danielmesquitta/api-finance-manager/internal/app/restapi/handler"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/usecase"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/dateutil"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/oauth/mockoauth"
	"github.com/stretchr/testify/assert"
)

func TestGetBalance(t *testing.T) {
	t.Parallel()

	type Test struct {
		description      string
		token            string
		queryParams      map[string]string
		expectedCode     int
		expectedResponse *dto.GetBalanceResponse
	}

	tests := []Test{
		{
			description:      "Fail to get balance without token",
			token:            "",
			expectedCode:     http.StatusBadRequest,
			expectedResponse: nil,
		},
		func() Test {
			startDateStr := "2024-11-01T00:00:00-03:00"
			endDateStr := "2024-11-30T23:59:59.999999999-03:00"

			startDate := dateutil.MustParseISOString(startDateStr)
			endDate := dateutil.MustParseISOString(endDateStr)
			cmpStartDate := dateutil.MustParseISOString(
				"2024-10-01T00:00:00-03:00",
			)
			cmpEndDate := dateutil.MustParseISOString(
				"2024-10-31T23:59:59.999999999-03:00",
			)

			return Test{
				description: "Get balance",
				token:       mockoauth.DefaultMockToken,
				queryParams: map[string]string{
					handler.QueryParamStartDate: startDateStr,
					handler.QueryParamEndDate:   endDateStr,
				},
				expectedCode: http.StatusOK,
				expectedResponse: &dto.GetBalanceResponse{
					GetBalanceOutput: usecase.GetBalanceOutput{
						ComparisonDates: dateutil.ComparisonDates{
							StartDate:           startDate,
							EndDate:             endDate,
							ComparisonStartDate: cmpStartDate,
							ComparisonEndDate:   cmpEndDate,
						},
						CurrentBalance:             12_505_08,
						PreviousBalance:            10_405_08,
						BalancePercentageVariation: 20_18,
						CurrentIncome:              160_758_24,
						PreviousIncome:             55_143_52,
						IncomePercentageVariation:  191_52,
						CurrentExpense:             -123_649_55,
						PreviousExpense:            -55_775_91,
						ExpensePercentageVariation: 121_68,
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

			var actualResponse dto.GetBalanceResponse
			statusCode, rawBody, err := app.MakeRequest(
				http.MethodGet,
				"/api/v1/balances",
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
				test.expectedResponse.GetBalanceOutput.ComparisonDates,
				actualResponse.GetBalanceOutput.ComparisonDates,
			)

			assert.Equal(
				t,
				test.expectedResponse.GetBalanceOutput.CurrentBalance,
				actualResponse.GetBalanceOutput.CurrentBalance,
			)

			assert.Equal(
				t,
				test.expectedResponse.GetBalanceOutput.PreviousBalance,
				actualResponse.GetBalanceOutput.PreviousBalance,
			)

			assert.Equal(
				t,
				test.expectedResponse.GetBalanceOutput.BalancePercentageVariation,
				actualResponse.GetBalanceOutput.BalancePercentageVariation,
			)

			assert.Equal(
				t,
				test.expectedResponse.GetBalanceOutput.CurrentIncome,
				actualResponse.GetBalanceOutput.CurrentIncome,
			)

			assert.Equal(
				t,
				test.expectedResponse.GetBalanceOutput.PreviousIncome,
				actualResponse.GetBalanceOutput.PreviousIncome,
			)

			assert.Equal(
				t,
				test.expectedResponse.GetBalanceOutput.IncomePercentageVariation,
				actualResponse.GetBalanceOutput.IncomePercentageVariation,
			)

			assert.Equal(
				t,
				test.expectedResponse.GetBalanceOutput.CurrentExpense,
				actualResponse.GetBalanceOutput.CurrentExpense,
			)

			assert.Equal(
				t,
				test.expectedResponse.GetBalanceOutput.PreviousExpense,
				actualResponse.GetBalanceOutput.PreviousExpense,
			)

			assert.Equal(
				t,
				test.expectedResponse.GetBalanceOutput.ExpensePercentageVariation,
				actualResponse.GetBalanceOutput.ExpensePercentageVariation,
			)
		})
	}
}
