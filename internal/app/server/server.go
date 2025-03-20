package server

import (
	"time"

	"github.com/danielmesquitta/api-finance-manager/internal/app/server/middleware"
	"github.com/danielmesquitta/api-finance-manager/internal/app/server/router"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/cache"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/cache/fibercache"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/db"

	"github.com/gofiber/fiber/v2"
	middlewareCache "github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/idempotency"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

type App struct {
	*fiber.App
	DB *db.DB
}

func Build(
	m *middleware.Middleware,
	r *router.Router,
	c cache.Cache,
	DB *db.DB,
) *App {
	app := fiber.New(fiber.Config{
		ErrorHandler: m.ErrorHandler,
	})

	app.Use(cors.New())
	app.Use(m.Recover())
	app.Use(requestid.New(requestid.Config{
		ContextKey: middleware.RequestIDContextKey,
	}))
	app.Use(middlewareCache.New(
		middlewareCache.Config{
			Storage: fibercache.NewFiberCache(c),
		},
	))
	app.Use(limiter.New(limiter.Config{
		Max:        20,
		Expiration: 1 * time.Minute,
		Storage:    fibercache.NewFiberCache(c),
	}))
	app.Use(idempotency.New(
		idempotency.Config{
			Storage: fibercache.NewFiberCache(c),
		},
	))
	app.Use(helmet.New())
	app.Use(m.Timeout(60 * time.Second))

	r.Register(app)

	return &App{
		App: app,
		DB:  DB,
	}
}
