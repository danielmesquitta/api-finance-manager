package testlog

import (
	"log/slog"
	"testing"

	"github.com/danielmesquitta/api-finance-manager/internal/config"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/log/jsonlog"
	"github.com/neilotoole/slogt"
)

func NewLogger(
	e *config.Env,
	t *testing.T,
) *slog.Logger {
	if t == nil {
		return jsonlog.NewLogger(e)
	}

	logger := slogt.New(t, slogt.JSON())
	return logger
}
