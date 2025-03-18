package main

import (
	"log"

	"github.com/danielmesquitta/api-finance-manager/internal/app/restapi"
	"github.com/danielmesquitta/api-finance-manager/internal/config"
	"github.com/danielmesquitta/api-finance-manager/internal/config/env"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/validator"
)

// @title API Finance Manager
// @version 1.0
// @description API Finance Manager
// @contact.name Daniel Mesquita
// @contact.email danielmesquitta123@gmail.com
// @BasePath /api
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
// @securityDefinitions.basic BasicAuth
func main() {
	v := validator.New()
	e := config.LoadConfig(v)

	var app *restapi.App
	switch e.Environment {
	case env.EnvironmentProduction:
		app = restapi.NewProd(v, e)

	case env.EnvironmentTest:
		app = restapi.NewTest(v, e)

	case env.EnvironmentStaging:
		app = restapi.NewStaging(v, e)

	default:
		app = restapi.NewDev(v, e)
	}

	if err := app.Listen(":" + e.Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
