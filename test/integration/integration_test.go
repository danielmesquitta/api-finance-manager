package integration

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/danielmesquitta/api-finance-manager/internal/app/restapi"
	"github.com/danielmesquitta/api-finance-manager/internal/config"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/validator"
)

type PostgresContainer struct {
	*postgres.PostgresContainer
	ConnectionString string
}

func createContainer(ctx context.Context) *PostgresContainer {
	migrations, err := filepath.Glob(
		filepath.Join("sql", "migrations", "**", "*.sql"),
	)
	if err != nil {
		panic(err)
	}

	pgContainer, err := postgres.Run(ctx,
		"postgres:alpine",
		postgres.WithInitScripts(migrations...),
		postgres.WithDatabase("test-db"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		panic(err)
	}
	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		panic(err)
	}

	return &PostgresContainer{
		PostgresContainer: pgContainer,
		ConnectionString:  connStr,
	}
}

func setup(
	ctx context.Context,
) (app *restapi.App, cleanUp func()) {
	pgContainer := createContainer(ctx)
	cleanUp = func() {
		_ = pgContainer.Terminate(context.Background())
	}

	v := validator.New()
	e := config.LoadEnv(v)

	e.Environment = config.EnvironmentTest
	e.DatabaseURL = pgContainer.ConnectionString

	fmt.Printf("%+v\n", e)

	app = restapi.NewDev(v, e)
	go func() {
		if err := app.Start(":" + e.Port); err != nil {
			log.Fatalf("failed to start server: %v", err)
		}
	}()

	return app, cleanUp
}
