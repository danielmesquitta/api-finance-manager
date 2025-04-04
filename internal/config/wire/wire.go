package wire

import (
	"github.com/google/wire"

	"github.com/danielmesquitta/api-finance-manager/internal/app/server"
	"github.com/danielmesquitta/api-finance-manager/internal/app/server/handler"
	"github.com/danielmesquitta/api-finance-manager/internal/app/server/middleware"
	"github.com/danielmesquitta/api-finance-manager/internal/app/server/router"
	"github.com/danielmesquitta/api-finance-manager/internal/config/env"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/usecase/account"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/usecase/aichat"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/usecase/auth"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/usecase/budget"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/usecase/calc"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/usecase/feedback"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/usecase/institution"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/usecase/paymentmethod"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/usecase/transaction"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/usecase/transactioncategory"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/usecase/user"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/hash"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/jwtutil"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/tx"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/validator"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/cache"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/cache/rediscache"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/db"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/db/query"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/gpt"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/gpt/openai"
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
	_ = params
}

func params(
	v *validator.Validator,
	e *env.Env,
) {
}

var providers = []any{
	jwtutil.NewJWT,
	hash.NewHasher,

	googleoauth.NewGoogleOAuth,

	pluggy.NewClient,

	db.NewPGXPool,
	db.NewSQLX,
	query.NewQueryBuilder,
	db.NewDB,

	wire.Bind(new(gpt.GPT), new(*openai.OpenAI)),
	openai.NewOpenAI,

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

	wire.Bind(new(repo.UserInstitutionRepo), new(*pgrepo.UserInstitutionRepo)),
	pgrepo.NewUserInstitutionRepo,

	wire.Bind(new(repo.AIChatRepo), new(*pgrepo.AIChatRepo)),
	pgrepo.NewAIChatRepo,

	wire.Bind(new(repo.AIChatMessageRepo), new(*pgrepo.AIChatMessageRepo)),
	pgrepo.NewAIChatMessageRepo,

	wire.Bind(new(repo.AIChatAnswerRepo), new(*pgrepo.AIChatAnswerRepo)),
	pgrepo.NewAIChatAnswerRepo,

	wire.Bind(
		new(repo.UserAuthProviderRepo),
		new(*pgrepo.UserAuthProviderRepo),
	),
	pgrepo.NewUserAuthProviderRepo,

	account.NewCreateAccountsUseCase,
	account.NewGetAccountsBalanceUseCase,
	account.NewSyncAccountsBalancesUseCase,

	aichat.NewListAIChatsUseCase,
	aichat.NewCreateAIChatUseCase,
	aichat.NewDeleteAIChatUseCase,
	aichat.NewUpdateAIChatUseCase,
	aichat.NewListAIChatMessagesAndAnswersUseCase,
	aichat.NewGenerateAIChatMessageUseCase,

	auth.NewSignInUseCase,
	auth.NewRefreshTokenUseCase,

	budget.NewUpsertBudgetUseCase,
	budget.NewGetBudgetUseCase,
	budget.NewDeleteBudgetUseCase,
	budget.NewGetBudgetCategoryUseCase,
	budget.NewListBudgetCategoryTransactionsUseCase,

	calc.NewCalculateCompoundInterestUseCase,
	calc.NewCalculateEmergencyReserveUseCase,
	calc.NewCalculateRetirementUseCase,
	calc.NewCalculateSimpleInterestUseCase,
	calc.NewCalculateCashVsInstallmentsUseCase,

	feedback.NewCreateFeedbackUseCase,

	institution.NewSyncInstitutionsUseCase,
	institution.NewListInstitutionsUseCase,

	paymentmethod.NewListPaymentMethodsUseCase,

	transaction.NewSyncTransactionsUseCase,
	transaction.NewListTransactionsUseCase,
	transaction.NewGetTransactionUseCase,
	transaction.NewUpdateTransactionUseCase,
	transaction.NewCreateTransactionUseCase,

	transactioncategory.NewSyncTransactionCategoriesUseCase,
	transactioncategory.NewListTransactionCategoriesUseCase,

	user.NewGetUserUseCase,
	user.NewUpdateUserUseCase,
	user.NewDeleteUserUseCase,

	handler.NewAuthHandler,
	handler.NewCalculatorHandler,
	handler.NewInstitutionHandler,
	handler.NewTransactionCategoryHandler,
	handler.NewBudgetHandler,
	handler.NewUserHandler,
	handler.NewAccountHandler,
	handler.NewTransactionHandler,
	handler.NewDocHandler,
	handler.NewFeedbackHandler,
	handler.NewPaymentMethodHandler,
	handler.NewAIChatHandler,
	handler.NewHealthHandler,

	middleware.NewMiddleware,

	router.NewRouter,

	server.Build,
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
