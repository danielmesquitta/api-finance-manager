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
	"github.com/danielmesquitta/api-finance-manager/internal/provider/db"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/db/query"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/oauth/googleoauth"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/oauth/mockoauth"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/openfinance/pluggy"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo/pgrepo"
)

// Injectors from wire.go:

func New() *App {
	validatorValidator := validator.NewValidator()
	env := config.LoadEnv(validatorValidator)
	jwt := jwtutil.NewJWT(env)
	middlewareMiddleware := middleware.NewMiddleware(env, jwt)
	healthHandler := handler.NewHealthHandler()
	pool := db.NewPGXPool(env)
	dbDB := db.NewQueries(pool)
	userPgRepo := pgrepo.NewUserPgRepo(dbDB)
	googleOAuth := googleoauth.NewGoogleOAuth()
	mockOAuth := mockoauth.NewMockOAuth(env)
	signInUseCase := usecase.NewSignInUseCase(validatorValidator, userPgRepo, jwt, googleOAuth, mockOAuth)
	refreshTokenUseCase := usecase.NewRefreshTokenUseCase(signInUseCase)
	authHandler := handler.NewAuthHandler(signInUseCase, refreshTokenUseCase)
	calculateCompoundInterestUseCase := usecase.NewCalculateCompoundInterestUseCase(validatorValidator)
	calculateEmergencyReserveUseCase := usecase.NewCalculateEmergencyReserveUseCase(validatorValidator)
	calculateRetirementUseCase := usecase.NewCalculateRetirementUseCase(validatorValidator, calculateCompoundInterestUseCase)
	calculateSimpleInterestUseCase := usecase.NewCalculateSimpleInterestUseCase(validatorValidator)
	calculatorHandler := handler.NewCalculatorHandler(calculateCompoundInterestUseCase, calculateEmergencyReserveUseCase, calculateRetirementUseCase, calculateSimpleInterestUseCase)
	client := pluggy.NewClient(env, jwt)
	institutionPgRepo := pgrepo.NewInstitutionPgRepo(dbDB)
	syncInstitutionsUseCase := usecase.NewSyncInstitutionsUseCase(client, institutionPgRepo)
	institutionHandler := handler.NewInstitutionHandler(syncInstitutionsUseCase)
	queryBuilder := query.NewQueryBuilder(env, dbDB)
	categoryPgRepo := pgrepo.NewCategoryPgRepo(dbDB, queryBuilder)
	syncCategoriesUseCase := usecase.NewSyncCategoriesUseCase(client, categoryPgRepo)
	listCategoriesUseCase := usecase.NewListCategoriesUseCase(categoryPgRepo)
	categoryHandler := handler.NewCategoryHandler(syncCategoriesUseCase, listCategoriesUseCase)
	pgxTX := tx.NewPgxTX(pool)
	budgetPgRepo := pgrepo.NewBudgetPgRepo(dbDB)
	upsertBudgetUseCase := usecase.NewUpsertBudgetUseCase(validatorValidator, pgxTX, budgetPgRepo)
	getBudgetUseCase := usecase.NewGetBudgetUseCase(validatorValidator, budgetPgRepo)
	deleteBudgetUseCase := usecase.NewDeleteBudgetUseCase(pgxTX, budgetPgRepo)
	budgetHandler := handler.NewBudgetHandler(upsertBudgetUseCase, getBudgetUseCase, deleteBudgetUseCase)
	routerRouter := router.NewRouter(env, middlewareMiddleware, healthHandler, authHandler, calculatorHandler, institutionHandler, categoryHandler, budgetHandler)
	app := newApp(env, middlewareMiddleware, routerRouter)
	return app
}
