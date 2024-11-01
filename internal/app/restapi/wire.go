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
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/jwtutil"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/validator"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/db"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/oauth"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/oauth/mockoauth"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo/pgrepo"
)

func New() *App {
	wire.Build(
		wire.Bind(new(validator.Validator), new(*validator.Validate)),
		validator.NewValidate,

		wire.Bind(new(jwtutil.JWTManager), new(*jwtutil.JWT)),
		jwtutil.NewJWT,

		config.LoadEnv,

		db.NewPGXConn,
		db.NewQueries,

		// wire.Bind(new(oauth.Provider), new(*googleoauth.GoogleOAuth)),
		// googleoauth.NewGoogleOAuth,

		wire.Bind(new(oauth.Provider), new(*mockoauth.MockOAuth)),
		mockoauth.NewMockOAuth,

		wire.Bind(new(repo.UserRepo), new(*pgrepo.UserPgRepo)),
		pgrepo.NewUserPgRepo,

		usecase.NewSignInUseCase,

		handler.NewAuthHandler,
		handler.NewHealthHandler,

		middleware.NewMiddleware,

		router.NewRouter,

		newApp,
	)

	return &App{}
}
