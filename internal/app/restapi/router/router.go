package router

import (
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"

	_ "github.com/danielmesquitta/api-finance-manager/docs" // swagger docs
	"github.com/danielmesquitta/api-finance-manager/internal/app/restapi/handler"
	"github.com/danielmesquitta/api-finance-manager/internal/app/restapi/middleware"
	"github.com/danielmesquitta/api-finance-manager/internal/config"
)

type Router struct {
	e  *config.Env
	m  *middleware.Middleware
	hh *handler.HealthHandler
	ah *handler.AuthHandler
	ch *handler.CalculatorHandler
}

func NewRouter(
	e *config.Env,
	m *middleware.Middleware,
	hh *handler.HealthHandler,
	ah *handler.AuthHandler,
	ch *handler.CalculatorHandler,
) *Router {
	return &Router{
		e:  e,
		m:  m,
		hh: hh,
		ah: ah,
		ch: ch,
	}
}

func (r *Router) Register(
	app *echo.Echo,
) {
	basePath := "/api"
	api := app.Group(basePath)

	api.GET("/docs/*", echoSwagger.WrapHandler)

	apiV1 := app.Group(basePath + "/v1")

	apiV1.GET("/health", r.hh.Health)
	apiV1.POST("/auth/sign-in", r.ah.SignIn)

	privateApiV1 := apiV1.Group("/", r.m.BearerAuth)

	privateApiV1.POST("/calculator/compound-interest", r.ch.CompoundInterest)
	privateApiV1.POST("/calculator/emergency-reserve", r.ch.EmergencyReserve)
	privateApiV1.POST("/calculator/retirement", r.ch.Retirement)
	privateApiV1.POST("/calculator/simple-interest", r.ch.SimpleInterest)
}
