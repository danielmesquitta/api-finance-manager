package server

import (
	"context"
	"net/http"
	"testing"

	"github.com/danielmesquitta/api-finance-manager/internal/app/server/dto"
	"github.com/stretchr/testify/assert"
)

func TestHealth(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description  string
		expectedCode int
		expectedBody *dto.HealthResponse
	}{
		{
			description:  "health check",
			expectedCode: 200,
			expectedBody: &dto.HealthResponse{Status: "ok"},
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

			var out dto.HealthResponse
			statusCode, rawBody, err := app.MakeRequest(
				http.MethodGet,
				"/api/health",
				WithResponse(&out),
			)
			assert.Nil(t, err)

			assert.Equal(
				t,
				test.expectedCode,
				statusCode,
				rawBody,
			)

			assert.Equal(
				t,
				test.expectedBody,
				&out,
			)
		})
	}
}
