package router

import (
	"github.com/gofiber/fiber/v2"

	_ "github.com/danielmesquitta/api-finance-manager/docs" // swagger docs
	"github.com/danielmesquitta/api-finance-manager/internal/app/restapi/handler"
	"github.com/danielmesquitta/api-finance-manager/internal/app/restapi/middleware"
	"github.com/danielmesquitta/api-finance-manager/internal/config"
)

type Router struct {
	e    *config.Env
	m    *middleware.Middleware
	hh   *handler.HealthHandler
	dh   *handler.DocHandler
	ah   *handler.AuthHandler
	ch   *handler.CalculatorHandler
	ih   *handler.InstitutionHandler
	cth  *handler.CategoryHandler
	bh   *handler.BudgetHandler
	uh   *handler.UserHandler
	ach  *handler.AccountHandler
	th   *handler.TransactionHandler
	bah  *handler.BalanceHandler
	fh   *handler.FeedbackHandler
	pmh  *handler.PaymentMethodHandler
	aih  *handler.AIChatHandler
	acmh *handler.AIChatMessageHandler
}

func NewRouter(
	e *config.Env,
	m *middleware.Middleware,
	hh *handler.HealthHandler,
	dh *handler.DocHandler,
	ah *handler.AuthHandler,
	ch *handler.CalculatorHandler,
	ih *handler.InstitutionHandler,
	cth *handler.CategoryHandler,
	bh *handler.BudgetHandler,
	uh *handler.UserHandler,
	ach *handler.AccountHandler,
	th *handler.TransactionHandler,
	bah *handler.BalanceHandler,
	fh *handler.FeedbackHandler,
	pmh *handler.PaymentMethodHandler,
	aih *handler.AIChatHandler,
	acmh *handler.AIChatMessageHandler,
) *Router {
	return &Router{
		e:    e,
		m:    m,
		hh:   hh,
		dh:   dh,
		ah:   ah,
		ch:   ch,
		ih:   ih,
		cth:  cth,
		bh:   bh,
		uh:   uh,
		ach:  ach,
		th:   th,
		bah:  bah,
		fh:   fh,
		pmh:  pmh,
		aih:  aih,
		acmh: acmh,
	}
}

func (r *Router) Register(
	app *fiber.App,
) {
	basePath := "/api"

	api := app.Group(basePath)

	api.Use("/docs", r.dh.Get)
	api.Get("/health", r.hh.Health)

	apiV1 := app.Group(basePath + "/v1")

	apiV1.Post("/auth/sign-in", r.ah.SignIn)
	apiV1.Post("/auth/refresh", r.m.BearerAuthRefreshToken(), r.ah.RefreshToken)

	adminApiV1 := apiV1.Group("/admin", r.m.BasicAuth())
	adminApiV1.Post("/institutions/sync", r.ih.Sync)
	adminApiV1.Post("/accounts", r.ach.Create)
	adminApiV1.Post("/transactions/categories/sync", r.cth.Sync)
	adminApiV1.Post("/transactions/sync", r.th.Sync)
	adminApiV1.Post("/balances/sync", r.bah.Sync)

	usersApiV1 := apiV1.Group("", r.m.BearerAuthAccessToken())

	usersApiV1.Get("/users/profile", r.uh.GetProfile)
	usersApiV1.Put("/users/profile", r.uh.UpdateProfile)
	usersApiV1.Delete("/users/profile", r.uh.DeleteProfile)

	usersApiV1.Post("/calculator/compound-interest", r.ch.CompoundInterest)
	usersApiV1.Post("/calculator/emergency-reserve", r.ch.EmergencyReserve)
	usersApiV1.Post("/calculator/retirement", r.ch.Retirement)
	usersApiV1.Post("/calculator/simple-interest", r.ch.SimpleInterest)
	usersApiV1.Post("/calculator/cash-vs-installments", r.ch.CashVsInstallments)

	usersApiV1.Get("/transactions/categories", r.cth.List)

	usersApiV1.Get("/institutions", r.ih.List)
	usersApiV1.Get("/users/institutions", r.ih.ListUserInstitutions)

	usersApiV1.Post("/budgets", r.bh.Upsert)
	usersApiV1.Get("/budgets", r.bh.Get)
	usersApiV1.Get(
		"/budgets/categories/:category_id",
		r.bh.GetTransactionCategory,
	)
	usersApiV1.Get(
		"/budgets/categories/:category_id/transactions",
		r.bh.ListCategoryTransactions,
	)
	usersApiV1.Delete("/budgets", r.bh.Delete)

	usersApiV1.Get("/balances", r.bah.Get)

	usersApiV1.Post("/transactions", r.th.Create)
	usersApiV1.Get("/transactions", r.th.List)
	usersApiV1.Get("/transactions/:transaction_id", r.th.Get)
	usersApiV1.Put("/transactions/:transaction_id", r.th.Update)

	usersApiV1.Post("/feedbacks", r.fh.Create)

	usersApiV1.Get("/payment-methods", r.pmh.List)

	usersApiV1.Post("/ai-chats", r.aih.Create)
	usersApiV1.Delete("/ai-chats/:ai_chat_id", r.aih.Delete)
	usersApiV1.Put("/ai-chats/:ai_chat_id", r.aih.Update)
	usersApiV1.Get("/ai-chats", r.aih.List)

	usersApiV1.Get("/ai-chats/:ai_chat_id/messages", r.acmh.List)
}
