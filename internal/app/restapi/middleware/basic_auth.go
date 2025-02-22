package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
)

func (m *Middleware) BasicAuth() fiber.Handler {
	middlewareFunc := basicauth.New(
		basicauth.Config{
			Users: map[string]string{
				m.e.BasicAuthUsername: m.e.BasicAuthPassword,
			},
		},
	)

	return middlewareFunc
}
