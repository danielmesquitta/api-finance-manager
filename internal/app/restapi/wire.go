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
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/tx"
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

		googleoauth.NewGoogleOAuth,
		mockoauth.NewMockOAuth,

		wire.Bind(new(openfinance.Client), new(*pluggy.Client)),
		pluggy.NewClient,

		db.NewPGXPool,
		db.NewQueries,

		wire.Bind(new(tx.TX), new(*tx.PgxTX)),
		tx.NewPgxTX,

		wire.Bind(new(repo.UserRepo), new(*pgrepo.UserPgRepo)),
		pgrepo.NewUserPgRepo,

		wire.Bind(new(repo.InstitutionRepo), new(*pgrepo.InstitutionPgRepo)),
		pgrepo.NewInstitutionPgRepo,

		wire.Bind(new(repo.CategoryRepo), new(*pgrepo.CategoryPgRepo)),
		pgrepo.NewCategoryPgRepo,

		wire.Bind(new(repo.BudgetRepo), new(*pgrepo.BudgetPgRepo)),
		pgrepo.NewBudgetPgRepo,

		usecase.NewSignInUseCase,
		usecase.NewRefreshTokenUseCase,
		usecase.NewCalculateCompoundInterestUseCase,
		usecase.NewCalculateEmergencyReserveUseCase,
		usecase.NewCalculateRetirementUseCase,
		usecase.NewCalculateSimpleInterestUseCase,
		usecase.NewSyncInstitutionsUseCase,
		usecase.NewSyncCategoriesUseCase,
		usecase.NewListCategoriesUseCase,
		usecase.NewUpsertBudgetUseCase,
		usecase.NewGetBudgetUseCase,
		usecase.NewDeleteBudgetUseCase,

		handler.NewAuthHandler,
		handler.NewHealthHandler,
		handler.NewCalculatorHandler,
		handler.NewInstitutionHandler,
		handler.NewCategoryHandler,
		handler.NewBudgetHandler,

		middleware.NewMiddleware,

		router.NewRouter,

		newApp,
	)

	return &App{}
}
