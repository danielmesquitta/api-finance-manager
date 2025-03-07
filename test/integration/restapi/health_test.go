package restapi

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/danielmesquitta/api-finance-manager/internal/app/restapi/dto"
	"github.com/stretchr/testify/assert"
)

func TestHealthRoutes(t *testing.T) {
	tests := []struct {
		description  string
		route        string
		expectedCode int
		expectedBody *dto.HealthResponse
	}{
		{
			description:  "live route",
			route:        "/api/health",
			expectedCode: 200,
			expectedBody: &dto.HealthResponse{Status: "ok"},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			t.Parallel()

			app, cleanUp := Setup(t)
			defer func() {
				err := cleanUp(context.Background())
				assert.Nil(t, err)
			}()

			req, _ := http.NewRequest(
				http.MethodGet,
				test.route,
				nil,
			)

			res, err := app.Test(req, -1)
			assert.Nil(t, err)

			assert.Equal(
				t,
				test.expectedCode,
				res.StatusCode,
			)

			bytesBody, err := io.ReadAll(res.Body)
			assert.Nil(t, err)

			var response dto.HealthResponse
			err = json.Unmarshal(bytesBody, &response)
			assert.Nil(t, err)

			assert.Equal(
				t,
				test.expectedBody,
				&response,
				string(bytesBody),
			)
		})
	}
}
