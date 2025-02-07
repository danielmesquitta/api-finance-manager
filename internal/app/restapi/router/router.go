package router

import (
	"io/fs"
	"net/http"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"

	root "github.com/danielmesquitta/api-finance-manager"
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
	bh  *handler.BudgetHandler
	uh  *handler.UserHandler
	ach *handler.AccountHandler
	th  *handler.TransactionHandler
	bah *handler.BalanceHandler
}

func NewRouter(
	e *config.Env,
	m *middleware.Middleware,
	hh *handler.HealthHandler,
	ah *handler.AuthHandler,
	ch *handler.CalculatorHandler,
	ih *handler.InstitutionHandler,
	cth *handler.CategoryHandler,
	bh *handler.BudgetHandler,
	uh *handler.UserHandler,
	ach *handler.AccountHandler,
	th *handler.TransactionHandler,
	bah *handler.BalanceHandler,
) *Router {
	return &Router{
		e:   e,
		m:   m,
		hh:  hh,
		ah:  ah,
		ch:  ch,
		ih:  ih,
		cth: cth,
		bh:  bh,
		uh:  uh,
		ach: ach,
		th:  th,
		bah: bah,
	}
}

func (r *Router) Register(
	app *echo.Echo,
) {
	basePath := "/api"

	api := app.Group(basePath)

	api.GET("/health", r.hh.Health)

	staticFiles, _ := fs.Sub(root.StaticFiles, "docs")
	staticFilesHandler := http.FileServer(http.FS(staticFiles))
	for _, ext := range []string{".json", ".yaml"} {
		api.GET(
			"/docs/openapi"+ext,
			echo.WrapHandler(
				http.StripPrefix("/api/docs/", staticFilesHandler),
			),
		)
	}

	api.GET("/docs/*", echoSwagger.WrapHandler)

	apiV1 := app.Group(basePath + "/v1")
	apiV1.POST("/auth/sign-in", r.ah.SignIn)
	apiV1.POST("/auth/refresh", r.ah.RefreshToken, r.m.BearerAuthRefreshToken)

	adminApiV1 := apiV1.Group("/admin", r.m.BasicAuth)
	adminApiV1.POST("/institutions/sync", r.ih.Sync)
	adminApiV1.POST("/accounts/sync", r.ach.Sync)
	adminApiV1.POST("/transactions/categories/sync", r.cth.Sync)
	adminApiV1.POST("/transactions/sync", r.th.Sync)

	usersApiV1 := apiV1.Group("", r.m.BearerAuthAccessToken())

	usersApiV1.GET("/users/profile", r.uh.Profile)

	usersApiV1.POST("/calculator/compound-interest", r.ch.CompoundInterest)
	usersApiV1.POST("/calculator/emergency-reserve", r.ch.EmergencyReserve)
	usersApiV1.POST("/calculator/retirement", r.ch.Retirement)
	usersApiV1.POST("/calculator/simple-interest", r.ch.SimpleInterest)
	usersApiV1.POST("/calculator/cash-vs-installments", r.ch.CashVsInstallments)

	usersApiV1.GET("/transactions/categories", r.cth.List)

	usersApiV1.GET("/institutions", r.ih.List)
	usersApiV1.GET("/users/institutions", r.ih.ListUserInstitutions)

	usersApiV1.POST("/budgets", r.bh.Upsert)
	usersApiV1.GET("/budgets", r.bh.Get)
	usersApiV1.GET(
		"/budgets/categories/:category_id",
		r.bh.GetTransactionCategory,
	)
	usersApiV1.GET(
		"/budgets/categories/:category_id/transactions",
		r.bh.ListCategoryTransactions,
	)
	usersApiV1.DELETE("/budgets", r.bh.Delete)

	usersApiV1.GET("/balances", r.bah.Get)

	usersApiV1.GET("/transactions", r.th.List)
	usersApiV1.GET("/transactions/:transaction_id", r.th.Get)
	usersApiV1.PUT("/transactions/:transaction_id", r.th.Update)
}
