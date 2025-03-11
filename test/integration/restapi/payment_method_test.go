package restapi

import (
	"context"
	"net/http"
	"testing"

	"github.com/danielmesquitta/api-finance-manager/internal/app/restapi/dto"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/oauth/mockoauth"
	"github.com/stretchr/testify/assert"
)

func TestListPaymentMethodsRoute(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description              string
		queryParams              map[string]string
		token                    string
		expectedCode             int
		expectedPaymentMethodIDs []string
	}{
		{
			description:              "Fail to list payment methods without token",
			queryParams:              map[string]string{},
			token:                    "",
			expectedCode:             http.StatusBadRequest,
			expectedPaymentMethodIDs: []string{},
		},
		{
			description:  "List all payment methods",
			queryParams:  map[string]string{},
			token:        mockoauth.DefaultMockToken,
			expectedCode: http.StatusOK,
			expectedPaymentMethodIDs: []string{
				"abbedc1f-0812-4ed1-9ec9-f51ca13e1069",
				"2158b0b6-844f-44b6-b487-282d0c1b045c",
				"66a40300-3fa2-415e-9480-ada220d07afb",
				"262f50e1-a751-4184-9427-90a23f485482",
				"1897980c-7f92-4eeb-8e6c-95690cea4ece",
			},
		},
		{
			description: "Search payment methods",
			queryParams: map[string]string{
				"search": "Cartao",
			},
			token:        mockoauth.DefaultMockToken,
			expectedCode: http.StatusOK,
			expectedPaymentMethodIDs: []string{
				"abbedc1f-0812-4ed1-9ec9-f51ca13e1069",
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

			var out dto.ListPaymentMethodsResponse
			statusCode, rawBody, err := app.MakeRequest(
				http.MethodGet,
				"/api/v1/payment-methods",
				WithQueryParams(test.queryParams),
				WithBearerToken(accessToken),
				WithResponse(&out),
			)
			assert.Nil(t, err)

			assert.Equal(
				t,
				test.expectedCode,
				statusCode,
			)

			if len(test.expectedPaymentMethodIDs) == 0 {
				assert.Empty(t, out.Items, rawBody)
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
