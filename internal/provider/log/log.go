package log

import (
	"log/slog"
	"os"
	"testing"

	"github.com/danielmesquitta/api-finance-manager/internal/config"
	"github.com/neilotoole/slogt"
)

func NewLogger(
	e *config.Env,
	t *testing.T,
) *slog.Logger {
	logger := slog.Default()

	switch e.Environment {
	case config.EnvironmentProduction, config.EnvironmentStaging:
		opts := &slog.HandlerOptions{}
		logger := slog.New(slog.NewJSONHandler(os.Stdout, opts))
		return logger

	case config.EnvironmentTest:
		logger := slogt.New(t)
		return logger

	default:
		return logger
	}
}
