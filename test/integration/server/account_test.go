package server

import (
	"context"
	"net/http"
	"testing"

	"github.com/danielmesquitta/api-finance-manager/internal/app/server/dto"
	"github.com/danielmesquitta/api-finance-manager/internal/app/server/handler"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/usecase/account"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/dateutil"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/oauth/mockoauth"
	"github.com/stretchr/testify/assert"
)

func TestGetAccountsBalance(t *testing.T) {
	t.Parallel()

	type Test struct {
		description      string
		token            string
		queryParams      map[string]string
		expectedCode     int
		expectedResponse *dto.GetAccountsBalanceResponse
	}

	tests := []Test{
		{
			description:      "fails without token",
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
				description: "gets accounts balance",
				token:       mockoauth.PremiumTierMockToken,
				queryParams: map[string]string{
					handler.QueryParamStartDate: startDateStr,
					handler.QueryParamEndDate:   endDateStr,
				},
				expectedCode: http.StatusOK,
				expectedResponse: &dto.GetAccountsBalanceResponse{
					GetAccountsBalanceUseCaseOutput: account.GetAccountsBalanceUseCaseOutput{
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

			signInRes := &dto.SignInResponse{}
			if test.token != "" {
				signInRes = app.SignIn(test.token)
			}

			var actualResponse dto.GetAccountsBalanceResponse
			statusCode, rawBody, err := app.MakeRequest(
				http.MethodGet,
				"/api/v1/accounts/balances",
				WithBearerToken(signInRes.AccessToken),
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
				test.expectedResponse.GetAccountsBalanceUseCaseOutput.ComparisonDates,
				actualResponse.GetAccountsBalanceUseCaseOutput.ComparisonDates,
			)

			assert.Equal(
				t,
				test.expectedResponse.GetAccountsBalanceUseCaseOutput.CurrentBalance,
				actualResponse.GetAccountsBalanceUseCaseOutput.CurrentBalance,
			)

			assert.Equal(
				t,
				test.expectedResponse.GetAccountsBalanceUseCaseOutput.PreviousBalance,
				actualResponse.GetAccountsBalanceUseCaseOutput.PreviousBalance,
			)

			assert.Equal(
				t,
				test.expectedResponse.GetAccountsBalanceUseCaseOutput.BalancePercentageVariation,
				actualResponse.GetAccountsBalanceUseCaseOutput.BalancePercentageVariation,
			)

			assert.Equal(
				t,
				test.expectedResponse.GetAccountsBalanceUseCaseOutput.CurrentIncome,
				actualResponse.GetAccountsBalanceUseCaseOutput.CurrentIncome,
			)

			assert.Equal(
				t,
				test.expectedResponse.GetAccountsBalanceUseCaseOutput.PreviousIncome,
				actualResponse.GetAccountsBalanceUseCaseOutput.PreviousIncome,
			)

			assert.Equal(
				t,
				test.expectedResponse.GetAccountsBalanceUseCaseOutput.IncomePercentageVariation,
				actualResponse.GetAccountsBalanceUseCaseOutput.IncomePercentageVariation,
			)

			assert.Equal(
				t,
				test.expectedResponse.GetAccountsBalanceUseCaseOutput.CurrentExpense,
				actualResponse.GetAccountsBalanceUseCaseOutput.CurrentExpense,
			)

			assert.Equal(
				t,
				test.expectedResponse.GetAccountsBalanceUseCaseOutput.PreviousExpense,
				actualResponse.GetAccountsBalanceUseCaseOutput.PreviousExpense,
			)

			assert.Equal(
				t,
				test.expectedResponse.GetAccountsBalanceUseCaseOutput.ExpensePercentageVariation,
				actualResponse.GetAccountsBalanceUseCaseOutput.ExpensePercentageVariation,
			)
		})
	}
}
