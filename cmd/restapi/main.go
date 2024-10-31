package main

import (
	"log"

	"github.com/danielmesquitta/api-finance-manager/internal/app/restapi"
)

func main() {
	app := restapi.New()

	if err := app.Start(":" + app.Env.Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
