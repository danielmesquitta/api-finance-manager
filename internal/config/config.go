package config

import (
	"log"

	"github.com/danielmesquitta/api-finance-manager/internal/pkg/validator"
)

func LoadConfig(v *validator.Validator) *Env {
	e := &Env{
		v: v,
	}

	if err := e.loadEnv(); err != nil {
		log.Fatalf("failed to load environment variables: %v", err)
	}

	setServerTimeZone()

	return e
}
