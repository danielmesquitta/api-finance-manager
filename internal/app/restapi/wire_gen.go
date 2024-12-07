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
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/validator"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/db"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/oauth/googleoauth"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo/pgrepo"
)

// Injectors from wire.go:

func New() *App {
	validate := validator.NewValidate()
	env := config.LoadEnv(validate)
	middlewareMiddleware := middleware.NewMiddleware()
	healthHandler := handler.NewHealthHandler()
	conn := db.NewPGXConn(env)
	queries := db.NewQueries(conn)
	userPgRepo := pgrepo.NewUserPgRepo(queries)
	jwt := jwtutil.NewJWT(env)
	googleOAuth := googleoauth.NewGoogleOAuth()
	signInUseCase := usecase.NewSignInUseCase(validate, userPgRepo, jwt, googleOAuth)
	authHandler := handler.NewAuthHandler(signInUseCase)
	routerRouter := router.NewRouter(env, healthHandler, authHandler)
	app := newApp(env, middlewareMiddleware, routerRouter)
	return app
}
