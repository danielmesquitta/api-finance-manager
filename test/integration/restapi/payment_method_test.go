package restapi

import (
	"context"
	"net/http"
	"testing"

	"github.com/danielmesquitta/api-finance-manager/internal/app/restapi/dto"
	"github.com/danielmesquitta/api-finance-manager/internal/app/restapi/handler"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/oauth/mockoauth"
	"github.com/stretchr/testify/assert"
)

func TestListPaymentMethods(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description              string
		queryParams              map[string]string
		token                    string
		expectedCode             int
		expectedPaymentMethodIDs []string
	}{
		{
			description:              "fails without token",
			queryParams:              map[string]string{},
			token:                    "",
			expectedCode:             http.StatusBadRequest,
			expectedPaymentMethodIDs: []string{},
		},
		{
			description:  "lists payment methods",
			queryParams:  map[string]string{},
			token:        mockoauth.PremiumTierMockToken,
			expectedCode: http.StatusOK,
			expectedPaymentMethodIDs: []string{
				"5d140153-c072-42ce-b19c-c5c9b528dba4",
				"61c45664-e3ea-44c4-8d9f-892088db9a8a",
				"b9098717-97ff-4051-a474-b0d703680176",
				"fc7adfa0-259c-430e-99f5-bef5281add10",
				"d3140f28-076f-4371-8c4b-b9e65e5367ef",
			},
		},
		{
			description: "searches payment methods",
			queryParams: map[string]string{
				handler.QueryParamSearch: "Cartao",
			},
			token:        mockoauth.PremiumTierMockToken,
			expectedCode: http.StatusOK,
			expectedPaymentMethodIDs: []string{
				"b9098717-97ff-4051-a474-b0d703680176",
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

			var out dto.ListPaymentMethodsResponse
			statusCode, rawBody, err := app.MakeRequest(
				http.MethodGet,
				"/api/v1/payment-methods",
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

			if len(test.expectedPaymentMethodIDs) == 0 {
				assert.Empty(t, out.Items)
				return
			}

			assert.Len(
				t,
				out.Items,
				len(test.expectedPaymentMethodIDs),
			)

			paymentMethodIDs := make([]string, len(out.Items))
			for i, paymentMethod := range out.Items {
				paymentMethodIDs[i] = paymentMethod.ID.String()
			}

			assert.ElementsMatch(
				t,
				test.expectedPaymentMethodIDs,
				paymentMethodIDs,
			)
		})
	}
}
