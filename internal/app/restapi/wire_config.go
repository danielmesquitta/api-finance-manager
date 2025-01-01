package restapi

import (
	"github.com/google/wire"

	"github.com/danielmesquitta/api-finance-manager/internal/app/restapi/handler"
	"github.com/danielmesquitta/api-finance-manager/internal/app/restapi/middleware"
	"github.com/danielmesquitta/api-finance-manager/internal/app/restapi/router"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/usecase"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/jwtutil"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/tx"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/db"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/db/query"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/oauth/googleoauth"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/oauth/mockoauth"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/openfinance"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/openfinance/mockpluggy"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/openfinance/pluggy"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo/pgrepo"
)

func init() {
	_ = providers
	_ = devProviders
	_ = prodProviders
}

var providers = []any{
	jwtutil.NewJWT,

	googleoauth.NewGoogleOAuth,

	pluggy.NewClient,

	db.NewPGXPool,
	db.NewQueries,
	query.NewQueryBuilder,

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

	wire.Bind(new(repo.AccountRepo), new(*pgrepo.AccountPgRepo)),
	pgrepo.NewAccountPgRepo,

	usecase.NewSignIn,
	usecase.NewRefreshToken,
	usecase.NewCalculateCompoundInterest,
	usecase.NewCalculateEmergencyReserve,
	usecase.NewCalculateRetirement,
	usecase.NewCalculateSimpleInterest,
	usecase.NewSyncInstitutions,
	usecase.NewSyncCategories,
	usecase.NewListCategories,
	usecase.NewUpsertBudget,
	usecase.NewGetBudget,
	usecase.NewDeleteBudget,
	usecase.NewGetUser,
	usecase.NewSyncAccounts,

	handler.NewAuthHandler,
	handler.NewHealthHandler,
	handler.NewCalculatorHandler,
	handler.NewInstitutionHandler,
	handler.NewCategoryHandler,
	handler.NewBudgetHandler,
	handler.NewUserHandler,
	handler.NewAccountHandler,

	middleware.NewMiddleware,

	router.NewRouter,

	newApp,
}

var devProviders = []any{
	mockoauth.NewMockOAuth,

	wire.Bind(new(openfinance.Client), new(*mockpluggy.Client)),
	mockpluggy.NewClient,
}

var prodProviders = []any{
	wire.Value((*mockoauth.MockOAuth)(nil)),

	wire.Bind(new(openfinance.Client), new(*pluggy.Client)),
}
