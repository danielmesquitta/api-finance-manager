package main

import (
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/danielmesquitta/api-finance-manager/internal/config"
	"github.com/danielmesquitta/api-finance-manager/internal/config/env"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/validator"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/db"
	"github.com/jmoiron/sqlx"
)

func main() {
	v := validator.New()
	e := config.LoadConfig(v)
	if e.PostgresTestDatabaseURL == "" {
		log.Fatal("POSTGRES_TEST_DATABASE_URL is not set")
	}

	testDBName, defaultConfig := prepareDatabaseConfigs(e)

	createTestDatabase(testDBName, &defaultConfig)

	pgx := db.NewPGXPool(e)
	sqlx := db.NewSQLX(pgx)

	if err := runSQLFiles(sqlx, "sql/migrations/*.sql"); err != nil {
		log.Fatalf("failed to run migration SQL files: %v", err)
		return
	}

	if err := runSQLFiles(sqlx, "sql/testdata/*.sql"); err != nil {
		log.Fatalf("failed to run test data SQL files: %v", err)
		return
	}
}

// prepareDatabaseConfigs extracts the test database name from the DSN and
// creates a default config to connect to the default database (typically "postgres").
func prepareDatabaseConfigs(e *env.Env) (string, env.Env) {
	parsedDSN, err := url.Parse(e.PostgresTestDatabaseURL)
	if err != nil {
		log.Fatalf("failed to parse DSN: %v", err)
	}
	testDBName := strings.TrimPrefix(parsedDSN.Path, "/")
	if testDBName == "" {
		log.Fatal("test database name could not be parsed from DSN")
	}

	// Change to the default database ("postgres") to perform drop/create.
	parsedDSN.Path = "/postgres"
	defaultDSN := parsedDSN.String()

	// Create a configuration copy for connecting to the default database.
	defaultConfig := *e
	defaultConfig.PostgresDatabaseURL = defaultDSN

	return testDBName, defaultConfig
}

// createTestDatabase drops the test database if it exists and recreates it.
func createTestDatabase(testDBName string, defaultConfig *env.Env) {
	defPGX := db.NewPGXPool(defaultConfig)
	defer defPGX.Close()
	defSQLX := db.NewSQLX(defPGX)

	if _, err := defSQLX.Exec("DROP DATABASE IF EXISTS " + testDBName); err != nil {
		log.Fatalf("failed to drop test database: %v", err)
	}

	if _, err := defSQLX.Exec("CREATE DATABASE " + testDBName); err != nil {
		log.Fatalf("failed to create test database: %v", err)
	}
}

// runSQLFiles reads and executes all SQL statements from files that match the given pattern.
func runSQLFiles(sqlx *sqlx.DB, pattern string) error {
	files, err := filepath.Glob(pattern)
	if err != nil {
		log.Fatalf("error retrieving files with pattern %s: %v", pattern, err)
	}
	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			return err
		}
		if _, err = sqlx.Exec(string(content)); err != nil {
			return err
		}
	}

	return nil
}
