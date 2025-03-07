package container

import (
	"context"
	"path/filepath"
	"runtime"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func NewPostgresContainer(
	ctx context.Context,
) (connectionString string, cleanUp func(context.Context) error) {
	_, filename, _, _ := runtime.Caller(0)
	baseDir := filepath.Dir(filename)
	projectRoot := filepath.Join(baseDir, "..", "..")

	migrations, err := filepath.Glob(
		filepath.Join(projectRoot, "sql", "migrations", "*.sql"),
	)
	if err != nil {
		panic(err)
	}

	testdata, err := filepath.Glob(
		filepath.Join(projectRoot, "sql", "testdata", "*.sql"),
	)
	if err != nil {
		panic(err)
	}

	initScripts := append(migrations, testdata...)

	pgCont, err := postgres.Run(ctx,
		"postgres:alpine",
		postgres.WithInitScripts(initScripts...),
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

	cleanUp = func(ctx context.Context) error {
		return pgCont.Terminate(ctx)
	}

	return connStr, cleanUp
}
