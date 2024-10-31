package restapi

import (
	"github.com/danielmesquitta/api-finance-manager/internal/app/restapi/middleware"
	"github.com/danielmesquitta/api-finance-manager/internal/app/restapi/router"
	"github.com/danielmesquitta/api-finance-manager/internal/config"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"
)

type App struct {
	*echo.Echo
	Env *config.Env
}

type T interface {
	string
	bool
}

func newApp(
	e *config.Env,
	m *middleware.Middleware,
	r *router.Router,
) *App {
	app := echo.New()

	defaultErrorHandler := app.HTTPErrorHandler
	customErrorHandler := m.ErrorHandler(defaultErrorHandler)
	app.HTTPErrorHandler = customErrorHandler

	app.Use(echoMiddleware.CORS())
	app.Use(echoMiddleware.Recover())
	app.Use(echoMiddleware.Timeout())
	app.Use(
		echoMiddleware.RateLimiter(
			echoMiddleware.NewRateLimiterMemoryStore(
				rate.Limit(20),
			),
		),
	)

	r.Register(app)

	return &App{
		Echo: app,
		Env:  e,
	}
}
