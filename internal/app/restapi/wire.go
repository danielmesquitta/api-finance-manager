//go:build wireinject
// +build wireinject

package restapi

import (
	"github.com/google/wire"

	"github.com/danielmesquitta/api-finance-manager/internal/app/restapi/handler"
	"github.com/danielmesquitta/api-finance-manager/internal/app/restapi/middleware"
	"github.com/danielmesquitta/api-finance-manager/internal/app/restapi/router"
	"github.com/danielmesquitta/api-finance-manager/internal/config"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/usecase"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/validator"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/db"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo/pgrepo"
)

func New() *App {
	wire.Build(
		wire.Bind(new(validator.Validator), new(*validator.Validate)),
		validator.NewValidate,

		config.LoadEnv,

		db.NewPGXConn,
		db.NewQueries,

		wire.Bind(new(repo.UserRepo), new(*pgrepo.UserPgRepo)),
		pgrepo.NewUserPgRepo,

		usecase.NewCreateUserUseCase,

		handler.NewUserHandler,
		handler.NewHealthHandler,

		middleware.NewMiddleware,

		router.NewRouter,

		newApp,
	)

	return &App{}
}
