package restapi

import (
	"github.com/google/wire"

	"github.com/danielmesquitta/api-finance-manager/internal/app/restapi/handler"
	"github.com/danielmesquitta/api-finance-manager/internal/app/restapi/middleware"
	"github.com/danielmesquitta/api-finance-manager/internal/app/restapi/router"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/usecase"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/hash"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/jwtutil"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/tx"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/cache"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/cache/rediscache"
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
	_ = testProviders
	_ = stagingProviders
	_ = prodProviders
}

var providers = []any{
	jwtutil.NewJWT,
	hash.NewHasher,

	googleoauth.NewGoogleOAuth,

	pluggy.NewClient,

	db.NewPGXPool,
	query.NewQueryBuilder,
	db.NewDB,

	wire.Bind(new(tx.TX), new(*tx.PgxTX)),
	tx.NewPgxTX,

	wire.Bind(new(cache.Cache), new(*rediscache.RedisCache)),
	rediscache.NewRedisCache,

	wire.Bind(new(repo.UserRepo), new(*pgrepo.UserRepo)),
	pgrepo.NewUserRepo,

	wire.Bind(new(repo.InstitutionRepo), new(*pgrepo.InstitutionRepo)),
	pgrepo.NewInstitutionRepo,

	wire.Bind(
		new(repo.TransactionCategoryRepo),
		new(*pgrepo.TransactionCategoryRepo),
	),
	pgrepo.NewCategoryRepo,

	wire.Bind(new(repo.BudgetRepo), new(*pgrepo.BudgetRepo)),
	pgrepo.NewBudgetRepo,

	wire.Bind(new(repo.AccountRepo), new(*pgrepo.AccountRepo)),
	pgrepo.NewAccountRepo,

	wire.Bind(new(repo.TransactionRepo), new(*pgrepo.TransactionRepo)),
	pgrepo.NewTransactionRepo,

	wire.Bind(new(repo.PaymentMethodRepo), new(*pgrepo.PaymentMethodRepo)),
	pgrepo.NewPaymentMethodRepo,

	wire.Bind(new(repo.AccountBalanceRepo), new(*pgrepo.AccountBalanceRepo)),
	pgrepo.NewAccountBalanceRepo,

	wire.Bind(new(repo.FeedbackRepo), new(*pgrepo.FeedbackRepo)),
	pgrepo.NewFeedbackRepo,

	wire.Bind(new(repo.AIChatRepo), new(*pgrepo.AIChatRepo)),
	pgrepo.NewAIChatRepo,

	wire.Bind(new(repo.AIChatMessageRepo), new(*pgrepo.AIChatMessageRepo)),
	pgrepo.NewAIChatMessageRepo,

	usecase.NewSignIn,
	usecase.NewRefreshToken,
	usecase.NewCalculateCompoundInterest,
	usecase.NewCalculateEmergencyReserve,
	usecase.NewCalculateRetirement,
	usecase.NewCalculateSimpleInterest,
	usecase.NewSyncInstitutions,
	usecase.NewSyncCategories,
	usecase.NewListTransactionCategories,
	usecase.NewUpsertBudget,
	usecase.NewGetBudget,
	usecase.NewDeleteBudget,
	usecase.NewGetUser,
	usecase.NewCreateAccounts,
	usecase.NewSyncTransactions,
	usecase.NewCalculateCashVsInstallments,
	usecase.NewListTransactions,
	usecase.NewListInstitutions,
	usecase.NewGetBudgetCategory,
	usecase.NewListBudgetCategoryTransactions,
	usecase.NewGetTransaction,
	usecase.NewUpdateTransaction,
	usecase.NewGetBalance,
	usecase.NewSyncBalances,
	usecase.NewCreateTransaction,
	usecase.NewUpdateUser,
	usecase.NewDeleteUser,
	usecase.NewCreateFeedback,
	usecase.NewListPaymentMethods,
	usecase.NewListAIChats,
	usecase.NewListAIChatMessages,
	usecase.NewCreateAIChat,
	usecase.NewDeleteAIChat,
	usecase.NewUpdateAIChat,

	handler.NewAuthHandler,
	handler.NewCalculatorHandler,
	handler.NewInstitutionHandler,
	handler.NewCategoryHandler,
	handler.NewBudgetHandler,
	handler.NewUserHandler,
	handler.NewAccountHandler,
	handler.NewTransactionHandler,
	handler.NewBalanceHandler,
	handler.NewDocHandler,
	handler.NewFeedbackHandler,
	handler.NewPaymentMethodHandler,
	handler.NewAIChatHandler,
	handler.NewAIChatMessageHandler,
	handler.NewHealthHandler,

	middleware.NewMiddleware,

	router.NewRouter,

	newApp,
}

var devProviders = []any{
	mockoauth.NewMockOAuth,

	wire.Bind(new(openfinance.Client), new(*mockpluggy.Client)),
	mockpluggy.NewClient,
}

var testProviders = []any{
	mockoauth.NewMockOAuth,

	wire.Bind(new(openfinance.Client), new(*mockpluggy.Client)),
	mockpluggy.NewClient,
}

var stagingProviders = []any{
	mockoauth.NewMockOAuth,

	wire.Bind(new(openfinance.Client), new(*mockpluggy.Client)),
	mockpluggy.NewClient,
}

var prodProviders = []any{
	wire.Value((*mockoauth.MockOAuth)(nil)),

	wire.Bind(new(openfinance.Client), new(*pluggy.Client)),
}
