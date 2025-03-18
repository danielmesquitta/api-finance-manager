package config

import (
	"github.com/danielmesquitta/api-finance-manager/internal/config/env"
	"github.com/danielmesquitta/api-finance-manager/internal/config/log"
	"github.com/danielmesquitta/api-finance-manager/internal/config/time"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/validator"
)

func LoadConfig(v *validator.Validator) *env.Env {
	e := env.NewEnv(v)
	log.SetDefaultLogger(e)
	time.SetServerTimeZone()

	return e
}
