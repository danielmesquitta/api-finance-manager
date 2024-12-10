package main

import (
	"log"

	"github.com/danielmesquitta/api-finance-manager/internal/app/restapi"
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
	app := restapi.New()
	if err := app.Start(":" + app.Env.Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
