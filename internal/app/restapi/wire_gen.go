// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package restapi

import (
	"github.com/danielmesquitta/api-finance-manager/internal/app/restapi/handler"
	"github.com/danielmesquitta/api-finance-manager/internal/app/restapi/middleware"
	"github.com/danielmesquitta/api-finance-manager/internal/app/restapi/router"
	"github.com/danielmesquitta/api-finance-manager/internal/config"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/usecase"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/jwtutil"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/tx"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/validator"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/cache/fibercache"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/cache/rediscache"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/db"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/db/query"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/oauth/googleoauth"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/oauth/mockoauth"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/openfinance/mockpluggy"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/openfinance/pluggy"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo/pgrepo"
)

// Injectors from wire_config_gen.go:

// NewDev wires up the application in development mode.
func NewDev(v *validator.Validator, e *config.Env) *App {
	jwt := jwtutil.NewJWT(e)
	middlewareMiddleware := middleware.NewMiddleware(e, jwt)
	pool := db.NewPGXPool(e)
	dbDB := db.NewDB(pool)
	userPgRepo := pgrepo.NewUserPgRepo(dbDB)
	googleOAuth := googleoauth.NewGoogleOAuth()
	mockOAuth := mockoauth.NewMockOAuth(e)
	signIn := usecase.NewSignIn(v, userPgRepo, jwt, googleOAuth, mockOAuth)
	refreshToken := usecase.NewRefreshToken(signIn)
	authHandler := handler.NewAuthHandler(signIn, refreshToken)
	calculateCompoundInterest := usecase.NewCalculateCompoundInterest(v)
	calculateEmergencyReserve := usecase.NewCalculateEmergencyReserve(v)
	calculateRetirement := usecase.NewCalculateRetirement(v, calculateCompoundInterest)
	calculateSimpleInterest := usecase.NewCalculateSimpleInterest(v)
	calculateCashVsInstallments := usecase.NewCalculateCashVsInstallments(v, calculateCompoundInterest)
	calculatorHandler := handler.NewCalculatorHandler(calculateCompoundInterest, calculateEmergencyReserve, calculateRetirement, calculateSimpleInterest, calculateCashVsInstallments)
	client := pluggy.NewClient(e, jwt)
	mockpluggyClient := mockpluggy.NewClient(client)
	queryBuilder := query.NewQueryBuilder(e, dbDB)
	institutionPgRepo := pgrepo.NewInstitutionPgRepo(dbDB, queryBuilder)
	syncInstitutions := usecase.NewSyncInstitutions(mockpluggyClient, institutionPgRepo)
	listInstitutions := usecase.NewListInstitutions(institutionPgRepo)
	institutionHandler := handler.NewInstitutionHandler(syncInstitutions, listInstitutions)
	categoryPgRepo := pgrepo.NewCategoryPgRepo(dbDB, queryBuilder)
	syncCategories := usecase.NewSyncCategories(mockpluggyClient, categoryPgRepo)
	listTransactionCategories := usecase.NewListTransactionCategories(categoryPgRepo)
	categoryHandler := handler.NewCategoryHandler(syncCategories, listTransactionCategories)
	pgxTX := tx.NewPgxTX(pool)
	budgetPgRepo := pgrepo.NewBudgetPgRepo(dbDB)
	upsertBudget := usecase.NewUpsertBudget(v, pgxTX, budgetPgRepo, categoryPgRepo)
	transactionPgRepo := pgrepo.NewTransactionPgRepo(dbDB, queryBuilder)
	getBudget := usecase.NewGetBudget(v, budgetPgRepo, transactionPgRepo)
	getBudgetCategory := usecase.NewGetBudgetCategory(v, budgetPgRepo, transactionPgRepo, categoryPgRepo)
	deleteBudget := usecase.NewDeleteBudget(pgxTX, budgetPgRepo)
	listTransactions := usecase.NewListTransactions(v, transactionPgRepo)
	listBudgetCategoryTransactions := usecase.NewListBudgetCategoryTransactions(v, listTransactions)
	budgetHandler := handler.NewBudgetHandler(upsertBudget, getBudget, getBudgetCategory, deleteBudget, listBudgetCategoryTransactions)
	getUser := usecase.NewGetUser(userPgRepo)
	userHandler := handler.NewUserHandler(getUser)
	accountPgRepo := pgrepo.NewAccountPgRepo(dbDB, queryBuilder)
	accountBalancePgRepo := pgrepo.NewAccountBalancePgRepo(dbDB)
	redisCache := rediscache.NewRedisCache(e)
	paymentMethodPgRepo := pgrepo.NewPaymentMethodPgRepo(dbDB, queryBuilder)
	syncTransactions := usecase.NewSyncTransactions(e, mockpluggyClient, redisCache, pgxTX, accountPgRepo, userPgRepo, transactionPgRepo, categoryPgRepo, paymentMethodPgRepo)
	createAccounts := usecase.NewCreateAccounts(v, mockpluggyClient, pgxTX, userPgRepo, accountPgRepo, accountBalancePgRepo, institutionPgRepo, syncTransactions)
	accountHandler := handler.NewAccountHandler(createAccounts)
	getTransaction := usecase.NewGetTransaction(transactionPgRepo)
	updateTransaction := usecase.NewUpdateTransaction(v, transactionPgRepo)
	transactionHandler := handler.NewTransactionHandler(syncTransactions, listTransactions, getTransaction, updateTransaction)
	getBalance := usecase.NewGetBalance(v, transactionPgRepo, accountBalancePgRepo)
	syncBalances := usecase.NewSyncBalances(e, mockpluggyClient, redisCache, accountPgRepo, accountBalancePgRepo)
	balanceHandler := handler.NewBalanceHandler(getBalance, syncBalances)
	routerRouter := router.NewRouter(e, middlewareMiddleware, authHandler, calculatorHandler, institutionHandler, categoryHandler, budgetHandler, userHandler, accountHandler, transactionHandler, balanceHandler)
	fiberCache := fibercache.NewFiberCache(redisCache)
	app := newApp(middlewareMiddleware, routerRouter, fiberCache)
	return app
}

// NewProd wires up the application in production mode.
func NewProd(v *validator.Validator, e *config.Env) *App {
	jwt := jwtutil.NewJWT(e)
	middlewareMiddleware := middleware.NewMiddleware(e, jwt)
	pool := db.NewPGXPool(e)
	dbDB := db.NewDB(pool)
	userPgRepo := pgrepo.NewUserPgRepo(dbDB)
	googleOAuth := googleoauth.NewGoogleOAuth()
	mockOAuth := _wireMockOAuthValue
	signIn := usecase.NewSignIn(v, userPgRepo, jwt, googleOAuth, mockOAuth)
	refreshToken := usecase.NewRefreshToken(signIn)
	authHandler := handler.NewAuthHandler(signIn, refreshToken)
	calculateCompoundInterest := usecase.NewCalculateCompoundInterest(v)
	calculateEmergencyReserve := usecase.NewCalculateEmergencyReserve(v)
	calculateRetirement := usecase.NewCalculateRetirement(v, calculateCompoundInterest)
	calculateSimpleInterest := usecase.NewCalculateSimpleInterest(v)
	calculateCashVsInstallments := usecase.NewCalculateCashVsInstallments(v, calculateCompoundInterest)
	calculatorHandler := handler.NewCalculatorHandler(calculateCompoundInterest, calculateEmergencyReserve, calculateRetirement, calculateSimpleInterest, calculateCashVsInstallments)
	client := pluggy.NewClient(e, jwt)
	queryBuilder := query.NewQueryBuilder(e, dbDB)
	institutionPgRepo := pgrepo.NewInstitutionPgRepo(dbDB, queryBuilder)
	syncInstitutions := usecase.NewSyncInstitutions(client, institutionPgRepo)
	listInstitutions := usecase.NewListInstitutions(institutionPgRepo)
	institutionHandler := handler.NewInstitutionHandler(syncInstitutions, listInstitutions)
	categoryPgRepo := pgrepo.NewCategoryPgRepo(dbDB, queryBuilder)
	syncCategories := usecase.NewSyncCategories(client, categoryPgRepo)
	listTransactionCategories := usecase.NewListTransactionCategories(categoryPgRepo)
	categoryHandler := handler.NewCategoryHandler(syncCategories, listTransactionCategories)
	pgxTX := tx.NewPgxTX(pool)
	budgetPgRepo := pgrepo.NewBudgetPgRepo(dbDB)
	upsertBudget := usecase.NewUpsertBudget(v, pgxTX, budgetPgRepo, categoryPgRepo)
	transactionPgRepo := pgrepo.NewTransactionPgRepo(dbDB, queryBuilder)
	getBudget := usecase.NewGetBudget(v, budgetPgRepo, transactionPgRepo)
	getBudgetCategory := usecase.NewGetBudgetCategory(v, budgetPgRepo, transactionPgRepo, categoryPgRepo)
	deleteBudget := usecase.NewDeleteBudget(pgxTX, budgetPgRepo)
	listTransactions := usecase.NewListTransactions(v, transactionPgRepo)
	listBudgetCategoryTransactions := usecase.NewListBudgetCategoryTransactions(v, listTransactions)
	budgetHandler := handler.NewBudgetHandler(upsertBudget, getBudget, getBudgetCategory, deleteBudget, listBudgetCategoryTransactions)
	getUser := usecase.NewGetUser(userPgRepo)
	userHandler := handler.NewUserHandler(getUser)
	accountPgRepo := pgrepo.NewAccountPgRepo(dbDB, queryBuilder)
	accountBalancePgRepo := pgrepo.NewAccountBalancePgRepo(dbDB)
	redisCache := rediscache.NewRedisCache(e)
	paymentMethodPgRepo := pgrepo.NewPaymentMethodPgRepo(dbDB, queryBuilder)
	syncTransactions := usecase.NewSyncTransactions(e, client, redisCache, pgxTX, accountPgRepo, userPgRepo, transactionPgRepo, categoryPgRepo, paymentMethodPgRepo)
	createAccounts := usecase.NewCreateAccounts(v, client, pgxTX, userPgRepo, accountPgRepo, accountBalancePgRepo, institutionPgRepo, syncTransactions)
	accountHandler := handler.NewAccountHandler(createAccounts)
	getTransaction := usecase.NewGetTransaction(transactionPgRepo)
	updateTransaction := usecase.NewUpdateTransaction(v, transactionPgRepo)
	transactionHandler := handler.NewTransactionHandler(syncTransactions, listTransactions, getTransaction, updateTransaction)
	getBalance := usecase.NewGetBalance(v, transactionPgRepo, accountBalancePgRepo)
	syncBalances := usecase.NewSyncBalances(e, client, redisCache, accountPgRepo, accountBalancePgRepo)
	balanceHandler := handler.NewBalanceHandler(getBalance, syncBalances)
	routerRouter := router.NewRouter(e, middlewareMiddleware, authHandler, calculatorHandler, institutionHandler, categoryHandler, budgetHandler, userHandler, accountHandler, transactionHandler, balanceHandler)
	fiberCache := fibercache.NewFiberCache(redisCache)
	app := newApp(middlewareMiddleware, routerRouter, fiberCache)
	return app
}

var (
	_wireMockOAuthValue = (*mockoauth.MockOAuth)(nil)
)
