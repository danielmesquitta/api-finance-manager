package middleware

import (
	"log/slog"

	"github.com/danielmesquitta/api-finance-manager/internal/config"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/jwtutil"
)

type Middleware struct {
	e *config.Env
	j *jwtutil.JWT
	l *slog.Logger
}

func NewMiddleware(
	e *config.Env,
	j *jwtutil.JWT,
	l *slog.Logger,
) *Middleware {
	return &Middleware{
		e: e,
		j: j,
		l: l,
	}
}
