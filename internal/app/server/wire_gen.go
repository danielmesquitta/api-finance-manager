// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package server

import (
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
	"github.com/danielmesquitta/api-finance-manager/internal/provider/cache/rediscache"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/db"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/db/query"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/gpt/openai"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/oauth/googleoauth"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/oauth/mockoauth"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/openfinance/mockpluggy"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/openfinance/pluggy"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo/pgrepo"
)

// Injectors from wire.go:

// NewTest wires up the application in test mode.
func NewTest(v *validator.Validator, e *env.Env) *App {
	jwt := jwtutil.NewJWT(e)
	middlewareMiddleware := middleware.NewMiddleware(e, jwt)
	healthHandler := handler.NewHealthHandler()
	docHandler := handler.NewDocHandler()
	pool := db.NewPGXPool(e)
	pgxTX := tx.NewPgxTX(pool)
	hasher := hash.NewHasher(e)
	sqlxDB := db.NewSQLX(pool)
	queryBuilder := query.NewQueryBuilder(e, sqlxDB)
	dbDB := db.NewDB(pool, queryBuilder)
	userRepo := pgrepo.NewUserRepo(dbDB)
	userAuthProviderRepo := pgrepo.NewUserAuthProviderRepo(dbDB)
	googleOAuth := googleoauth.NewGoogleOAuth()
	mockOAuth := mockoauth.NewMockOAuth(e)
	signInUseCase := auth.NewSignInUseCase(v, pgxTX, hasher, userRepo, userAuthProviderRepo, jwt, googleOAuth, mockOAuth)
	refreshTokenUseCase := auth.NewRefreshTokenUseCase(signInUseCase)
	authHandler := handler.NewAuthHandler(signInUseCase, refreshTokenUseCase)
	calculateCompoundInterestUseCase := calc.NewCalculateCompoundInterestUseCase(v)
	calculateEmergencyReserveUseCase := calc.NewCalculateEmergencyReserveUseCase(v)
	calculateRetirementUseCase := calc.NewCalculateRetirementUseCase(v, calculateCompoundInterestUseCase)
	calculateSimpleInterestUseCase := calc.NewCalculateSimpleInterestUseCase(v)
	calculateCashVsInstallmentsUseCase := calc.NewCalculateCashVsInstallmentsUseCase(v, calculateCompoundInterestUseCase)
	calculatorHandler := handler.NewCalculatorHandler(calculateCompoundInterestUseCase, calculateEmergencyReserveUseCase, calculateRetirementUseCase, calculateSimpleInterestUseCase, calculateCashVsInstallmentsUseCase)
	client := pluggy.NewClient(e, jwt)
	mockpluggyClient := mockpluggy.NewClient(client)
	institutionRepo := pgrepo.NewInstitutionRepo(dbDB)
	syncInstitutionsUseCase := institution.NewSyncInstitutionsUseCase(mockpluggyClient, institutionRepo)
	listInstitutionsUseCase := institution.NewListInstitutionsUseCase(institutionRepo)
	institutionHandler := handler.NewInstitutionHandler(syncInstitutionsUseCase, listInstitutionsUseCase)
	transactionCategoryRepo := pgrepo.NewCategoryRepo(dbDB)
	syncTransactionCategoriesUseCase := transactioncategory.NewSyncTransactionCategoriesUseCase(mockpluggyClient, transactionCategoryRepo)
	listTransactionCategoriesUseCase := transactioncategory.NewListTransactionCategoriesUseCase(transactionCategoryRepo)
	transactionCategoryHandler := handler.NewTransactionCategoryHandler(syncTransactionCategoriesUseCase, listTransactionCategoriesUseCase)
	budgetRepo := pgrepo.NewBudgetRepo(dbDB)
	upsertBudgetUseCase := budget.NewUpsertBudgetUseCase(v, pgxTX, budgetRepo, transactionCategoryRepo)
	transactionRepo := pgrepo.NewTransactionRepo(dbDB)
	getBudgetUseCase := budget.NewGetBudgetUseCase(v, budgetRepo, transactionRepo)
	getBudgetCategoryUseCase := budget.NewGetBudgetCategoryUseCase(v, budgetRepo, transactionRepo, transactionCategoryRepo)
	deleteBudgetUseCase := budget.NewDeleteBudgetUseCase(pgxTX, budgetRepo)
	listTransactionsUseCase := transaction.NewListTransactionsUseCase(v, transactionRepo)
	listBudgetCategoryTransactionsUseCase := budget.NewListBudgetCategoryTransactionsUseCase(v, listTransactionsUseCase)
	budgetHandler := handler.NewBudgetHandler(upsertBudgetUseCase, getBudgetUseCase, getBudgetCategoryUseCase, deleteBudgetUseCase, listBudgetCategoryTransactionsUseCase)
	getUserUseCase := user.NewGetUserUseCase(userRepo)
	updateUserUseCase := user.NewUpdateUserUseCase(v, userRepo)
	deleteUserUseCase := user.NewDeleteUserUseCase(hasher, userRepo)
	userHandler := handler.NewUserHandler(getUserUseCase, updateUserUseCase, deleteUserUseCase)
	accountRepo := pgrepo.NewAccountRepo(dbDB, queryBuilder)
	accountBalanceRepo := pgrepo.NewAccountBalanceRepo(dbDB)
	userInstitutionRepo := pgrepo.NewUserInstitutionRepo(dbDB)
	createAccountsUseCase := account.NewCreateAccountsUseCase(v, mockpluggyClient, pgxTX, accountRepo, accountBalanceRepo, institutionRepo, userInstitutionRepo)
	getAccountsBalanceUseCase := account.NewGetAccountsBalanceUseCase(v, transactionRepo, accountBalanceRepo)
	redisCache := rediscache.NewRedisCache(e)
	syncAccountsBalancesUseCase := account.NewSyncAccountsBalancesUseCase(e, pgxTX, mockpluggyClient, redisCache, accountRepo, accountBalanceRepo)
	accountHandler := handler.NewAccountHandler(createAccountsUseCase, getAccountsBalanceUseCase, syncAccountsBalancesUseCase)
	paymentMethodRepo := pgrepo.NewPaymentMethodRepo(dbDB)
	syncTransactionsUseCase := transaction.NewSyncTransactionsUseCase(e, mockpluggyClient, redisCache, pgxTX, accountRepo, userRepo, transactionRepo, transactionCategoryRepo, paymentMethodRepo)
	getTransactionUseCase := transaction.NewGetTransactionUseCase(transactionRepo)
	updateTransactionUseCase := transaction.NewUpdateTransactionUseCase(v, transactionRepo)
	createTransactionUseCase := transaction.NewCreateTransactionUseCase(v, transactionRepo, userRepo, transactionCategoryRepo, paymentMethodRepo)
	transactionHandler := handler.NewTransactionHandler(syncTransactionsUseCase, listTransactionsUseCase, getTransactionUseCase, updateTransactionUseCase, createTransactionUseCase)
	feedbackRepo := pgrepo.NewFeedbackRepo(dbDB)
	createFeedbackUseCase := feedback.NewCreateFeedbackUseCase(v, feedbackRepo)
	feedbackHandler := handler.NewFeedbackHandler(createFeedbackUseCase)
	listPaymentMethodsUseCase := paymentmethod.NewListPaymentMethodsUseCase(paymentMethodRepo)
	paymentMethodHandler := handler.NewPaymentMethodHandler(listPaymentMethodsUseCase)
	aiChatRepo := pgrepo.NewAIChatRepo(dbDB)
	createAIChatUseCase := aichat.NewCreateAIChatUseCase(v, aiChatRepo)
	deleteAIChatUseCase := aichat.NewDeleteAIChatUseCase(aiChatRepo)
	updateAIChatUseCase := aichat.NewUpdateAIChatUseCase(v, aiChatRepo)
	listAIChatsUseCase := aichat.NewListAIChatsUseCase(aiChatRepo)
	listAIChatMessagesAndAnswersUseCase := aichat.NewListAIChatMessagesAndAnswersUseCase(aiChatRepo)
	openAI := openai.NewOpenAI(e)
	aiChatMessageRepo := pgrepo.NewAIChatMessageRepo(dbDB)
	aiChatAnswerRepo := pgrepo.NewAIChatAnswerRepo(dbDB)
	generateAIChatMessageUseCase := aichat.NewGenerateAIChatMessageUseCase(v, pgxTX, openAI, aiChatRepo, aiChatMessageRepo, aiChatAnswerRepo, paymentMethodRepo, transactionCategoryRepo, institutionRepo, transactionRepo, getBudgetUseCase, getAccountsBalanceUseCase)
	aiChatHandler := handler.NewAIChatHandler(createAIChatUseCase, deleteAIChatUseCase, updateAIChatUseCase, listAIChatsUseCase, listAIChatMessagesAndAnswersUseCase, generateAIChatMessageUseCase)
	routerRouter := router.NewRouter(e, middlewareMiddleware, healthHandler, docHandler, authHandler, calculatorHandler, institutionHandler, transactionCategoryHandler, budgetHandler, userHandler, accountHandler, transactionHandler, feedbackHandler, paymentMethodHandler, aiChatHandler)
	app := Build(middlewareMiddleware, routerRouter, redisCache, dbDB)
	return app
}

// NewProd wires up the application in prod mode.
func NewProd(v *validator.Validator, e *env.Env) *App {
	jwt := jwtutil.NewJWT(e)
	middlewareMiddleware := middleware.NewMiddleware(e, jwt)
	healthHandler := handler.NewHealthHandler()
	docHandler := handler.NewDocHandler()
	pool := db.NewPGXPool(e)
	pgxTX := tx.NewPgxTX(pool)
	hasher := hash.NewHasher(e)
	sqlxDB := db.NewSQLX(pool)
	queryBuilder := query.NewQueryBuilder(e, sqlxDB)
	dbDB := db.NewDB(pool, queryBuilder)
	userRepo := pgrepo.NewUserRepo(dbDB)
	userAuthProviderRepo := pgrepo.NewUserAuthProviderRepo(dbDB)
	googleOAuth := googleoauth.NewGoogleOAuth()
	mockOAuth := _wireMockOAuthValue
	signInUseCase := auth.NewSignInUseCase(v, pgxTX, hasher, userRepo, userAuthProviderRepo, jwt, googleOAuth, mockOAuth)
	refreshTokenUseCase := auth.NewRefreshTokenUseCase(signInUseCase)
	authHandler := handler.NewAuthHandler(signInUseCase, refreshTokenUseCase)
	calculateCompoundInterestUseCase := calc.NewCalculateCompoundInterestUseCase(v)
	calculateEmergencyReserveUseCase := calc.NewCalculateEmergencyReserveUseCase(v)
	calculateRetirementUseCase := calc.NewCalculateRetirementUseCase(v, calculateCompoundInterestUseCase)
	calculateSimpleInterestUseCase := calc.NewCalculateSimpleInterestUseCase(v)
	calculateCashVsInstallmentsUseCase := calc.NewCalculateCashVsInstallmentsUseCase(v, calculateCompoundInterestUseCase)
	calculatorHandler := handler.NewCalculatorHandler(calculateCompoundInterestUseCase, calculateEmergencyReserveUseCase, calculateRetirementUseCase, calculateSimpleInterestUseCase, calculateCashVsInstallmentsUseCase)
	client := pluggy.NewClient(e, jwt)
	institutionRepo := pgrepo.NewInstitutionRepo(dbDB)
	syncInstitutionsUseCase := institution.NewSyncInstitutionsUseCase(client, institutionRepo)
	listInstitutionsUseCase := institution.NewListInstitutionsUseCase(institutionRepo)
	institutionHandler := handler.NewInstitutionHandler(syncInstitutionsUseCase, listInstitutionsUseCase)
	transactionCategoryRepo := pgrepo.NewCategoryRepo(dbDB)
	syncTransactionCategoriesUseCase := transactioncategory.NewSyncTransactionCategoriesUseCase(client, transactionCategoryRepo)
	listTransactionCategoriesUseCase := transactioncategory.NewListTransactionCategoriesUseCase(transactionCategoryRepo)
	transactionCategoryHandler := handler.NewTransactionCategoryHandler(syncTransactionCategoriesUseCase, listTransactionCategoriesUseCase)
	budgetRepo := pgrepo.NewBudgetRepo(dbDB)
	upsertBudgetUseCase := budget.NewUpsertBudgetUseCase(v, pgxTX, budgetRepo, transactionCategoryRepo)
	transactionRepo := pgrepo.NewTransactionRepo(dbDB)
	getBudgetUseCase := budget.NewGetBudgetUseCase(v, budgetRepo, transactionRepo)
	getBudgetCategoryUseCase := budget.NewGetBudgetCategoryUseCase(v, budgetRepo, transactionRepo, transactionCategoryRepo)
	deleteBudgetUseCase := budget.NewDeleteBudgetUseCase(pgxTX, budgetRepo)
	listTransactionsUseCase := transaction.NewListTransactionsUseCase(v, transactionRepo)
	listBudgetCategoryTransactionsUseCase := budget.NewListBudgetCategoryTransactionsUseCase(v, listTransactionsUseCase)
	budgetHandler := handler.NewBudgetHandler(upsertBudgetUseCase, getBudgetUseCase, getBudgetCategoryUseCase, deleteBudgetUseCase, listBudgetCategoryTransactionsUseCase)
	getUserUseCase := user.NewGetUserUseCase(userRepo)
	updateUserUseCase := user.NewUpdateUserUseCase(v, userRepo)
	deleteUserUseCase := user.NewDeleteUserUseCase(hasher, userRepo)
	userHandler := handler.NewUserHandler(getUserUseCase, updateUserUseCase, deleteUserUseCase)
	accountRepo := pgrepo.NewAccountRepo(dbDB, queryBuilder)
	accountBalanceRepo := pgrepo.NewAccountBalanceRepo(dbDB)
	userInstitutionRepo := pgrepo.NewUserInstitutionRepo(dbDB)
	createAccountsUseCase := account.NewCreateAccountsUseCase(v, client, pgxTX, accountRepo, accountBalanceRepo, institutionRepo, userInstitutionRepo)
	getAccountsBalanceUseCase := account.NewGetAccountsBalanceUseCase(v, transactionRepo, accountBalanceRepo)
	redisCache := rediscache.NewRedisCache(e)
	syncAccountsBalancesUseCase := account.NewSyncAccountsBalancesUseCase(e, pgxTX, client, redisCache, accountRepo, accountBalanceRepo)
	accountHandler := handler.NewAccountHandler(createAccountsUseCase, getAccountsBalanceUseCase, syncAccountsBalancesUseCase)
	paymentMethodRepo := pgrepo.NewPaymentMethodRepo(dbDB)
	syncTransactionsUseCase := transaction.NewSyncTransactionsUseCase(e, client, redisCache, pgxTX, accountRepo, userRepo, transactionRepo, transactionCategoryRepo, paymentMethodRepo)
	getTransactionUseCase := transaction.NewGetTransactionUseCase(transactionRepo)
	updateTransactionUseCase := transaction.NewUpdateTransactionUseCase(v, transactionRepo)
	createTransactionUseCase := transaction.NewCreateTransactionUseCase(v, transactionRepo, userRepo, transactionCategoryRepo, paymentMethodRepo)
	transactionHandler := handler.NewTransactionHandler(syncTransactionsUseCase, listTransactionsUseCase, getTransactionUseCase, updateTransactionUseCase, createTransactionUseCase)
	feedbackRepo := pgrepo.NewFeedbackRepo(dbDB)
	createFeedbackUseCase := feedback.NewCreateFeedbackUseCase(v, feedbackRepo)
	feedbackHandler := handler.NewFeedbackHandler(createFeedbackUseCase)
	listPaymentMethodsUseCase := paymentmethod.NewListPaymentMethodsUseCase(paymentMethodRepo)
	paymentMethodHandler := handler.NewPaymentMethodHandler(listPaymentMethodsUseCase)
	aiChatRepo := pgrepo.NewAIChatRepo(dbDB)
	createAIChatUseCase := aichat.NewCreateAIChatUseCase(v, aiChatRepo)
	deleteAIChatUseCase := aichat.NewDeleteAIChatUseCase(aiChatRepo)
	updateAIChatUseCase := aichat.NewUpdateAIChatUseCase(v, aiChatRepo)
	listAIChatsUseCase := aichat.NewListAIChatsUseCase(aiChatRepo)
	listAIChatMessagesAndAnswersUseCase := aichat.NewListAIChatMessagesAndAnswersUseCase(aiChatRepo)
	openAI := openai.NewOpenAI(e)
	aiChatMessageRepo := pgrepo.NewAIChatMessageRepo(dbDB)
	aiChatAnswerRepo := pgrepo.NewAIChatAnswerRepo(dbDB)
	generateAIChatMessageUseCase := aichat.NewGenerateAIChatMessageUseCase(v, pgxTX, openAI, aiChatRepo, aiChatMessageRepo, aiChatAnswerRepo, paymentMethodRepo, transactionCategoryRepo, institutionRepo, transactionRepo, getBudgetUseCase, getAccountsBalanceUseCase)
	aiChatHandler := handler.NewAIChatHandler(createAIChatUseCase, deleteAIChatUseCase, updateAIChatUseCase, listAIChatsUseCase, listAIChatMessagesAndAnswersUseCase, generateAIChatMessageUseCase)
	routerRouter := router.NewRouter(e, middlewareMiddleware, healthHandler, docHandler, authHandler, calculatorHandler, institutionHandler, transactionCategoryHandler, budgetHandler, userHandler, accountHandler, transactionHandler, feedbackHandler, paymentMethodHandler, aiChatHandler)
	app := Build(middlewareMiddleware, routerRouter, redisCache, dbDB)
	return app
}

var (
	_wireMockOAuthValue = (*mockoauth.MockOAuth)(nil)
)

// NewDev wires up the application in dev mode.
func NewDev(v *validator.Validator, e *env.Env) *App {
	jwt := jwtutil.NewJWT(e)
	middlewareMiddleware := middleware.NewMiddleware(e, jwt)
	healthHandler := handler.NewHealthHandler()
	docHandler := handler.NewDocHandler()
	pool := db.NewPGXPool(e)
	pgxTX := tx.NewPgxTX(pool)
	hasher := hash.NewHasher(e)
	sqlxDB := db.NewSQLX(pool)
	queryBuilder := query.NewQueryBuilder(e, sqlxDB)
	dbDB := db.NewDB(pool, queryBuilder)
	userRepo := pgrepo.NewUserRepo(dbDB)
	userAuthProviderRepo := pgrepo.NewUserAuthProviderRepo(dbDB)
	googleOAuth := googleoauth.NewGoogleOAuth()
	mockOAuth := mockoauth.NewMockOAuth(e)
	signInUseCase := auth.NewSignInUseCase(v, pgxTX, hasher, userRepo, userAuthProviderRepo, jwt, googleOAuth, mockOAuth)
	refreshTokenUseCase := auth.NewRefreshTokenUseCase(signInUseCase)
	authHandler := handler.NewAuthHandler(signInUseCase, refreshTokenUseCase)
	calculateCompoundInterestUseCase := calc.NewCalculateCompoundInterestUseCase(v)
	calculateEmergencyReserveUseCase := calc.NewCalculateEmergencyReserveUseCase(v)
	calculateRetirementUseCase := calc.NewCalculateRetirementUseCase(v, calculateCompoundInterestUseCase)
	calculateSimpleInterestUseCase := calc.NewCalculateSimpleInterestUseCase(v)
	calculateCashVsInstallmentsUseCase := calc.NewCalculateCashVsInstallmentsUseCase(v, calculateCompoundInterestUseCase)
	calculatorHandler := handler.NewCalculatorHandler(calculateCompoundInterestUseCase, calculateEmergencyReserveUseCase, calculateRetirementUseCase, calculateSimpleInterestUseCase, calculateCashVsInstallmentsUseCase)
	client := pluggy.NewClient(e, jwt)
	mockpluggyClient := mockpluggy.NewClient(client)
	institutionRepo := pgrepo.NewInstitutionRepo(dbDB)
	syncInstitutionsUseCase := institution.NewSyncInstitutionsUseCase(mockpluggyClient, institutionRepo)
	listInstitutionsUseCase := institution.NewListInstitutionsUseCase(institutionRepo)
	institutionHandler := handler.NewInstitutionHandler(syncInstitutionsUseCase, listInstitutionsUseCase)
	transactionCategoryRepo := pgrepo.NewCategoryRepo(dbDB)
	syncTransactionCategoriesUseCase := transactioncategory.NewSyncTransactionCategoriesUseCase(mockpluggyClient, transactionCategoryRepo)
	listTransactionCategoriesUseCase := transactioncategory.NewListTransactionCategoriesUseCase(transactionCategoryRepo)
	transactionCategoryHandler := handler.NewTransactionCategoryHandler(syncTransactionCategoriesUseCase, listTransactionCategoriesUseCase)
	budgetRepo := pgrepo.NewBudgetRepo(dbDB)
	upsertBudgetUseCase := budget.NewUpsertBudgetUseCase(v, pgxTX, budgetRepo, transactionCategoryRepo)
	transactionRepo := pgrepo.NewTransactionRepo(dbDB)
	getBudgetUseCase := budget.NewGetBudgetUseCase(v, budgetRepo, transactionRepo)
	getBudgetCategoryUseCase := budget.NewGetBudgetCategoryUseCase(v, budgetRepo, transactionRepo, transactionCategoryRepo)
	deleteBudgetUseCase := budget.NewDeleteBudgetUseCase(pgxTX, budgetRepo)
	listTransactionsUseCase := transaction.NewListTransactionsUseCase(v, transactionRepo)
	listBudgetCategoryTransactionsUseCase := budget.NewListBudgetCategoryTransactionsUseCase(v, listTransactionsUseCase)
	budgetHandler := handler.NewBudgetHandler(upsertBudgetUseCase, getBudgetUseCase, getBudgetCategoryUseCase, deleteBudgetUseCase, listBudgetCategoryTransactionsUseCase)
	getUserUseCase := user.NewGetUserUseCase(userRepo)
	updateUserUseCase := user.NewUpdateUserUseCase(v, userRepo)
	deleteUserUseCase := user.NewDeleteUserUseCase(hasher, userRepo)
	userHandler := handler.NewUserHandler(getUserUseCase, updateUserUseCase, deleteUserUseCase)
	accountRepo := pgrepo.NewAccountRepo(dbDB, queryBuilder)
	accountBalanceRepo := pgrepo.NewAccountBalanceRepo(dbDB)
	userInstitutionRepo := pgrepo.NewUserInstitutionRepo(dbDB)
	createAccountsUseCase := account.NewCreateAccountsUseCase(v, mockpluggyClient, pgxTX, accountRepo, accountBalanceRepo, institutionRepo, userInstitutionRepo)
	getAccountsBalanceUseCase := account.NewGetAccountsBalanceUseCase(v, transactionRepo, accountBalanceRepo)
	redisCache := rediscache.NewRedisCache(e)
	syncAccountsBalancesUseCase := account.NewSyncAccountsBalancesUseCase(e, pgxTX, mockpluggyClient, redisCache, accountRepo, accountBalanceRepo)
	accountHandler := handler.NewAccountHandler(createAccountsUseCase, getAccountsBalanceUseCase, syncAccountsBalancesUseCase)
	paymentMethodRepo := pgrepo.NewPaymentMethodRepo(dbDB)
	syncTransactionsUseCase := transaction.NewSyncTransactionsUseCase(e, mockpluggyClient, redisCache, pgxTX, accountRepo, userRepo, transactionRepo, transactionCategoryRepo, paymentMethodRepo)
	getTransactionUseCase := transaction.NewGetTransactionUseCase(transactionRepo)
	updateTransactionUseCase := transaction.NewUpdateTransactionUseCase(v, transactionRepo)
	createTransactionUseCase := transaction.NewCreateTransactionUseCase(v, transactionRepo, userRepo, transactionCategoryRepo, paymentMethodRepo)
	transactionHandler := handler.NewTransactionHandler(syncTransactionsUseCase, listTransactionsUseCase, getTransactionUseCase, updateTransactionUseCase, createTransactionUseCase)
	feedbackRepo := pgrepo.NewFeedbackRepo(dbDB)
	createFeedbackUseCase := feedback.NewCreateFeedbackUseCase(v, feedbackRepo)
	feedbackHandler := handler.NewFeedbackHandler(createFeedbackUseCase)
	listPaymentMethodsUseCase := paymentmethod.NewListPaymentMethodsUseCase(paymentMethodRepo)
	paymentMethodHandler := handler.NewPaymentMethodHandler(listPaymentMethodsUseCase)
	aiChatRepo := pgrepo.NewAIChatRepo(dbDB)
	createAIChatUseCase := aichat.NewCreateAIChatUseCase(v, aiChatRepo)
	deleteAIChatUseCase := aichat.NewDeleteAIChatUseCase(aiChatRepo)
	updateAIChatUseCase := aichat.NewUpdateAIChatUseCase(v, aiChatRepo)
	listAIChatsUseCase := aichat.NewListAIChatsUseCase(aiChatRepo)
	listAIChatMessagesAndAnswersUseCase := aichat.NewListAIChatMessagesAndAnswersUseCase(aiChatRepo)
	openAI := openai.NewOpenAI(e)
	aiChatMessageRepo := pgrepo.NewAIChatMessageRepo(dbDB)
	aiChatAnswerRepo := pgrepo.NewAIChatAnswerRepo(dbDB)
	generateAIChatMessageUseCase := aichat.NewGenerateAIChatMessageUseCase(v, pgxTX, openAI, aiChatRepo, aiChatMessageRepo, aiChatAnswerRepo, paymentMethodRepo, transactionCategoryRepo, institutionRepo, transactionRepo, getBudgetUseCase, getAccountsBalanceUseCase)
	aiChatHandler := handler.NewAIChatHandler(createAIChatUseCase, deleteAIChatUseCase, updateAIChatUseCase, listAIChatsUseCase, listAIChatMessagesAndAnswersUseCase, generateAIChatMessageUseCase)
	routerRouter := router.NewRouter(e, middlewareMiddleware, healthHandler, docHandler, authHandler, calculatorHandler, institutionHandler, transactionCategoryHandler, budgetHandler, userHandler, accountHandler, transactionHandler, feedbackHandler, paymentMethodHandler, aiChatHandler)
	app := Build(middlewareMiddleware, routerRouter, redisCache, dbDB)
	return app
}

// NewStaging wires up the application in staging mode.
func NewStaging(v *validator.Validator, e *env.Env) *App {
	jwt := jwtutil.NewJWT(e)
	middlewareMiddleware := middleware.NewMiddleware(e, jwt)
	healthHandler := handler.NewHealthHandler()
	docHandler := handler.NewDocHandler()
	pool := db.NewPGXPool(e)
	pgxTX := tx.NewPgxTX(pool)
	hasher := hash.NewHasher(e)
	sqlxDB := db.NewSQLX(pool)
	queryBuilder := query.NewQueryBuilder(e, sqlxDB)
	dbDB := db.NewDB(pool, queryBuilder)
	userRepo := pgrepo.NewUserRepo(dbDB)
	userAuthProviderRepo := pgrepo.NewUserAuthProviderRepo(dbDB)
	googleOAuth := googleoauth.NewGoogleOAuth()
	mockOAuth := mockoauth.NewMockOAuth(e)
	signInUseCase := auth.NewSignInUseCase(v, pgxTX, hasher, userRepo, userAuthProviderRepo, jwt, googleOAuth, mockOAuth)
	refreshTokenUseCase := auth.NewRefreshTokenUseCase(signInUseCase)
	authHandler := handler.NewAuthHandler(signInUseCase, refreshTokenUseCase)
	calculateCompoundInterestUseCase := calc.NewCalculateCompoundInterestUseCase(v)
	calculateEmergencyReserveUseCase := calc.NewCalculateEmergencyReserveUseCase(v)
	calculateRetirementUseCase := calc.NewCalculateRetirementUseCase(v, calculateCompoundInterestUseCase)
	calculateSimpleInterestUseCase := calc.NewCalculateSimpleInterestUseCase(v)
	calculateCashVsInstallmentsUseCase := calc.NewCalculateCashVsInstallmentsUseCase(v, calculateCompoundInterestUseCase)
	calculatorHandler := handler.NewCalculatorHandler(calculateCompoundInterestUseCase, calculateEmergencyReserveUseCase, calculateRetirementUseCase, calculateSimpleInterestUseCase, calculateCashVsInstallmentsUseCase)
	client := pluggy.NewClient(e, jwt)
	mockpluggyClient := mockpluggy.NewClient(client)
	institutionRepo := pgrepo.NewInstitutionRepo(dbDB)
	syncInstitutionsUseCase := institution.NewSyncInstitutionsUseCase(mockpluggyClient, institutionRepo)
	listInstitutionsUseCase := institution.NewListInstitutionsUseCase(institutionRepo)
	institutionHandler := handler.NewInstitutionHandler(syncInstitutionsUseCase, listInstitutionsUseCase)
	transactionCategoryRepo := pgrepo.NewCategoryRepo(dbDB)
	syncTransactionCategoriesUseCase := transactioncategory.NewSyncTransactionCategoriesUseCase(mockpluggyClient, transactionCategoryRepo)
	listTransactionCategoriesUseCase := transactioncategory.NewListTransactionCategoriesUseCase(transactionCategoryRepo)
	transactionCategoryHandler := handler.NewTransactionCategoryHandler(syncTransactionCategoriesUseCase, listTransactionCategoriesUseCase)
	budgetRepo := pgrepo.NewBudgetRepo(dbDB)
	upsertBudgetUseCase := budget.NewUpsertBudgetUseCase(v, pgxTX, budgetRepo, transactionCategoryRepo)
	transactionRepo := pgrepo.NewTransactionRepo(dbDB)
	getBudgetUseCase := budget.NewGetBudgetUseCase(v, budgetRepo, transactionRepo)
	getBudgetCategoryUseCase := budget.NewGetBudgetCategoryUseCase(v, budgetRepo, transactionRepo, transactionCategoryRepo)
	deleteBudgetUseCase := budget.NewDeleteBudgetUseCase(pgxTX, budgetRepo)
	listTransactionsUseCase := transaction.NewListTransactionsUseCase(v, transactionRepo)
	listBudgetCategoryTransactionsUseCase := budget.NewListBudgetCategoryTransactionsUseCase(v, listTransactionsUseCase)
	budgetHandler := handler.NewBudgetHandler(upsertBudgetUseCase, getBudgetUseCase, getBudgetCategoryUseCase, deleteBudgetUseCase, listBudgetCategoryTransactionsUseCase)
	getUserUseCase := user.NewGetUserUseCase(userRepo)
	updateUserUseCase := user.NewUpdateUserUseCase(v, userRepo)
	deleteUserUseCase := user.NewDeleteUserUseCase(hasher, userRepo)
	userHandler := handler.NewUserHandler(getUserUseCase, updateUserUseCase, deleteUserUseCase)
	accountRepo := pgrepo.NewAccountRepo(dbDB, queryBuilder)
	accountBalanceRepo := pgrepo.NewAccountBalanceRepo(dbDB)
	userInstitutionRepo := pgrepo.NewUserInstitutionRepo(dbDB)
	createAccountsUseCase := account.NewCreateAccountsUseCase(v, mockpluggyClient, pgxTX, accountRepo, accountBalanceRepo, institutionRepo, userInstitutionRepo)
	getAccountsBalanceUseCase := account.NewGetAccountsBalanceUseCase(v, transactionRepo, accountBalanceRepo)
	redisCache := rediscache.NewRedisCache(e)
	syncAccountsBalancesUseCase := account.NewSyncAccountsBalancesUseCase(e, pgxTX, mockpluggyClient, redisCache, accountRepo, accountBalanceRepo)
	accountHandler := handler.NewAccountHandler(createAccountsUseCase, getAccountsBalanceUseCase, syncAccountsBalancesUseCase)
	paymentMethodRepo := pgrepo.NewPaymentMethodRepo(dbDB)
	syncTransactionsUseCase := transaction.NewSyncTransactionsUseCase(e, mockpluggyClient, redisCache, pgxTX, accountRepo, userRepo, transactionRepo, transactionCategoryRepo, paymentMethodRepo)
	getTransactionUseCase := transaction.NewGetTransactionUseCase(transactionRepo)
	updateTransactionUseCase := transaction.NewUpdateTransactionUseCase(v, transactionRepo)
	createTransactionUseCase := transaction.NewCreateTransactionUseCase(v, transactionRepo, userRepo, transactionCategoryRepo, paymentMethodRepo)
	transactionHandler := handler.NewTransactionHandler(syncTransactionsUseCase, listTransactionsUseCase, getTransactionUseCase, updateTransactionUseCase, createTransactionUseCase)
	feedbackRepo := pgrepo.NewFeedbackRepo(dbDB)
	createFeedbackUseCase := feedback.NewCreateFeedbackUseCase(v, feedbackRepo)
	feedbackHandler := handler.NewFeedbackHandler(createFeedbackUseCase)
	listPaymentMethodsUseCase := paymentmethod.NewListPaymentMethodsUseCase(paymentMethodRepo)
	paymentMethodHandler := handler.NewPaymentMethodHandler(listPaymentMethodsUseCase)
	aiChatRepo := pgrepo.NewAIChatRepo(dbDB)
	createAIChatUseCase := aichat.NewCreateAIChatUseCase(v, aiChatRepo)
	deleteAIChatUseCase := aichat.NewDeleteAIChatUseCase(aiChatRepo)
	updateAIChatUseCase := aichat.NewUpdateAIChatUseCase(v, aiChatRepo)
	listAIChatsUseCase := aichat.NewListAIChatsUseCase(aiChatRepo)
	listAIChatMessagesAndAnswersUseCase := aichat.NewListAIChatMessagesAndAnswersUseCase(aiChatRepo)
	openAI := openai.NewOpenAI(e)
	aiChatMessageRepo := pgrepo.NewAIChatMessageRepo(dbDB)
	aiChatAnswerRepo := pgrepo.NewAIChatAnswerRepo(dbDB)
	generateAIChatMessageUseCase := aichat.NewGenerateAIChatMessageUseCase(v, pgxTX, openAI, aiChatRepo, aiChatMessageRepo, aiChatAnswerRepo, paymentMethodRepo, transactionCategoryRepo, institutionRepo, transactionRepo, getBudgetUseCase, getAccountsBalanceUseCase)
	aiChatHandler := handler.NewAIChatHandler(createAIChatUseCase, deleteAIChatUseCase, updateAIChatUseCase, listAIChatsUseCase, listAIChatMessagesAndAnswersUseCase, generateAIChatMessageUseCase)
	routerRouter := router.NewRouter(e, middlewareMiddleware, healthHandler, docHandler, authHandler, calculatorHandler, institutionHandler, transactionCategoryHandler, budgetHandler, userHandler, accountHandler, transactionHandler, feedbackHandler, paymentMethodHandler, aiChatHandler)
	app := Build(middlewareMiddleware, routerRouter, redisCache, dbDB)
	return app
}
