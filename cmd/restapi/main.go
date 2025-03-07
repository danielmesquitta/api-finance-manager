package main

import (
	"log"
	"log/slog"

	"github.com/danielmesquitta/api-finance-manager/internal/app/restapi"
	"github.com/danielmesquitta/api-finance-manager/internal/config"
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

	slog.Info("Starting server...")

	var app *restapi.App
	switch e.Environment {
	case config.EnvironmentProduction:
		app = restapi.NewProd(v, e, nil)

	case config.EnvironmentTest:
		app = restapi.NewTest(v, e, nil)

	case config.EnvironmentStaging:
		app = restapi.NewStaging(v, e, nil)

	default:
		app = restapi.NewDev(v, e, nil)
	}

	if err := app.Listen(":" + e.Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
