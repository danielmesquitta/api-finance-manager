package middleware

import (
	"strings"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/labstack/echo/v4"
)

func (m *Middleware) BearerAuth(
	next echo.HandlerFunc,
) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Get the Authorization header
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return errs.ErrUnauthorized
		}

		// Split the header to get the token part
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			return errs.ErrUnauthorized
		}

		accessToken := parts[1]

		// Parse and validate the token
		claims, err := m.j.Parse(accessToken)
		if err != nil {
			return errs.ErrUnauthorized
		}

		// Set the claims in the context
		c.Set("claims", claims)

		// Token is valid, proceed with the request
		return next(c)
	}
}
