package container

import (
	"context"
	"path/filepath"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func NewPostgresContainer(
	ctx context.Context,
) (connectionString string, cleanUp func(context.Context) error) {
	migrations, err := filepath.Glob(
		filepath.Join("sql", "migrations", "**", "*.sql"),
	)
	if err != nil {
		panic(err)
	}

	pgCont, err := postgres.Run(ctx,
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
	connStr, err := pgCont.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		panic(err)
	}

	return connStr, pgCont.Terminate
}
