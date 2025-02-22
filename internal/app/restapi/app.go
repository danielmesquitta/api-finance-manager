package restapi

import (
	"net/http"
	"time"

	root "github.com/danielmesquitta/api-finance-manager"
	"github.com/danielmesquitta/api-finance-manager/internal/app/restapi/middleware"
	"github.com/danielmesquitta/api-finance-manager/internal/app/restapi/router"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/idempotency"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

type App struct {
	*fiber.App
}

func newApp(
	m *middleware.Middleware,
	r *router.Router,
) *App {
	app := fiber.New(fiber.Config{
		ErrorHandler: m.ErrorHandler,
	})

	app.Use(cors.New())
	app.Use(recover.New())
	app.Use(requestid.New())
	app.Use(limiter.New())
	app.Use(idempotency.New())
	app.Use(helmet.New())
	app.Use(healthcheck.New())
	app.Use(cache.New())
	app.Use(csrf.New())
	app.Use(m.Timeout(60 * time.Second))

	app.Use("/docs", filesystem.New(filesystem.Config{
		Root:       http.FS(root.StaticFiles),
		PathPrefix: "docs",
		Browse:     true,
	}))

	r.Register(app)

	return &App{
		App: app,
	}
}
