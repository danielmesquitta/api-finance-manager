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
	"github.com/danielmesquitta/api-finance-manager/internal/provider/openfinance"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/openfinance/pluggy"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo/pgrepo"
)

func New() *App {
	wire.Build(
		validator.NewValidator,

		jwtutil.NewJWT,

		config.LoadEnv,

		db.NewPGXConn,
		db.NewQueries,

		googleoauth.NewGoogleOAuth,
		mockoauth.NewMockOAuth,

		wire.Bind(new(openfinance.Client), new(*pluggy.Client)),
		pluggy.NewClient,

		wire.Bind(new(repo.UserRepo), new(*pgrepo.UserPgRepo)),
		pgrepo.NewUserPgRepo,

		wire.Bind(new(repo.InstitutionRepo), new(*pgrepo.InstitutionPgRepo)),
		pgrepo.NewInstitutionPgRepo,

		usecase.NewSignInUseCase,
		usecase.NewRefreshTokenUseCase,
		usecase.NewCalculateCompoundInterestUseCase,
		usecase.NewCalculateEmergencyReserveUseCase,
		usecase.NewCalculateRetirementUseCase,
		usecase.NewCalculateSimpleInterestUseCase,
		usecase.NewSyncInstitutionsUseCase,

		handler.NewAuthHandler,
		handler.NewHealthHandler,
		handler.NewCalculatorHandler,
		handler.NewInstitutionHandler,

		middleware.NewMiddleware,

		router.NewRouter,

		newApp,
	)

	return &App{}
}
