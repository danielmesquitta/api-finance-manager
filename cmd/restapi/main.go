package main

import (
	"log"

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
	e := config.LoadEnv(v)

	app := restapi.New(v, e)
	if err := app.Start(":" + e.Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
