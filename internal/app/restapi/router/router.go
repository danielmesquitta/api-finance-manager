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
	e   *config.Env
	m   *middleware.Middleware
	hh  *handler.HealthHandler
	ah  *handler.AuthHandler
	ch  *handler.CalculatorHandler
	ih  *handler.InstitutionHandler
	cth *handler.CategoryHandler
}

func NewRouter(
	e *config.Env,
	m *middleware.Middleware,
	hh *handler.HealthHandler,
	ah *handler.AuthHandler,
	ch *handler.CalculatorHandler,
	ih *handler.InstitutionHandler,
	cth *handler.CategoryHandler,
) *Router {
	return &Router{
		e:   e,
		m:   m,
		hh:  hh,
		ah:  ah,
		ch:  ch,
		ih:  ih,
		cth: cth,
	}
}

func (r *Router) Register(
	app *echo.Echo,
) {
	basePath := "/api"

	api := app.Group(basePath)
	api.GET("/docs/*", echoSwagger.WrapHandler)
	api.GET("/health", r.hh.Health)

	apiV1 := app.Group(basePath + "/v1")
	apiV1.POST("/auth/sign-in", r.ah.SignIn)
	apiV1.POST("/auth/refresh", r.ah.RefreshToken, r.m.BearerAuthRefreshToken)

	adminApiV1 := apiV1.Group("/admin", r.m.BasicAuth)
	adminApiV1.POST("/institutions/sync", r.ih.Sync)
	adminApiV1.POST("/categories/sync", r.cth.Sync)

	usersApiV1 := apiV1.Group("", r.m.BearerAuthAccessToken())
	usersApiV1.POST("/calculator/compound-interest", r.ch.CompoundInterest)
	usersApiV1.POST("/calculator/emergency-reserve", r.ch.EmergencyReserve)
	usersApiV1.POST("/calculator/retirement", r.ch.Retirement)
	usersApiV1.POST("/calculator/simple-interest", r.ch.SimpleInterest)

	usersApiV1.GET("/categories", r.cth.List)
}
