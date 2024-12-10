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
	"github.com/danielmesquitta/api-finance-manager/internal/provider/oauth/googleoauth"
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

		// PRODUCTION
		// wire.Bind(new(googleoauth.Provider), new(*googleoauth.GoogleOAuth)),
		// googleoauth.NewGoogleOAuth,

		// DEVELOPMENT
		wire.Bind(new(googleoauth.Provider), new(*mockoauth.MockOAuth)),
		mockoauth.NewMockOAuth,

		wire.Bind(new(repo.UserRepo), new(*pgrepo.UserPgRepo)),
		pgrepo.NewUserPgRepo,

		usecase.NewSignInUseCase,
		usecase.NewCalculateCompoundInterestUseCase,
		usecase.NewCalculateEmergencyReserveUseCase,
		usecase.NewCalculateRetirementUseCase,
		usecase.NewCalculateSimpleInterestUseCase,

		handler.NewAuthHandler,
		handler.NewHealthHandler,
		handler.NewCalculatorHandler,

		middleware.NewMiddleware,

		router.NewRouter,

		newApp,
	)

	return &App{}
}
