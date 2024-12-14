package middleware

import (
	"github.com/danielmesquitta/api-finance-manager/internal/config"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/jwtutil"
)

type Middleware struct {
	e *config.Env
	j *jwtutil.JWT
}

func NewMiddleware(
	e *config.Env,
	j *jwtutil.JWT,
) *Middleware {
	return &Middleware{
		e: e,
		j: j,
	}
}
