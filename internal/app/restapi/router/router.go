package router

import (
	"github.com/danielmesquitta/api-finance-manager/internal/app/restapi/handler"
	"github.com/danielmesquitta/api-finance-manager/internal/config"
	"github.com/labstack/echo/v4"
)

type Router struct {
	e  *config.Env
	hh *handler.HealthHandler
	uh *handler.UserHandler
}

func NewRouter(
	e *config.Env,
	hh *handler.HealthHandler,
	uh *handler.UserHandler,
) *Router {
	return &Router{
		e:  e,
		hh: hh,
		uh: uh,
	}
}

func (r *Router) Register(
	app *echo.Echo,
) {
	basePath := "/api/v1"
	apiV1 := app.Group(basePath)

	apiV1.GET("/health", r.hh.Health)

	apiV1.POST("/users", r.uh.Create)
}
