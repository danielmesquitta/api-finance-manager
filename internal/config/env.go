package config

import (
	"bytes"
	_ "embed"
	"log"

	root "github.com/danielmesquitta/api-finance-manager"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/validator"
	"github.com/spf13/viper"
)

type Environment string

const (
	EnvironmentDevelopment Environment = "development"
	EnvironmentProduction  Environment = "production"
	EnvironmentStaging     Environment = "staging"
	EnvironmentTest        Environment = "test"
)

type Env struct {
	v *validator.Validator

	Environment                      Environment `mapstructure:"ENVIRONMENT"                         validate:"required,oneof=development production staging test"`
	Port                             string      `mapstructure:"PORT"`
	DatabaseURL                      string      `mapstructure:"DATABASE_URL"                        validate:"required"`
	JWTAccessTokenSecretKey          string      `mapstructure:"JWT_ACCESS_TOKEN_SECRET_KEY"         validate:"required"`
	JWTRefreshTokenSecretKey         string      `mapstructure:"JWT_REFRESH_TOKEN_SECRET_KEY"        validate:"required"`
	PluggyClientID                   string      `mapstructure:"PLUGGY_CLIENT_ID"                    validate:"required"`
	PluggyClientSecret               string      `mapstructure:"PLUGGY_CLIENT_SECRET"                validate:"required"`
	BasicAuthUsername                string      `mapstructure:"BASIC_AUTH_USERNAME"                 validate:"required"`
	BasicAuthPassword                string      `mapstructure:"BASIC_AUTH_PASSWORD"                 validate:"required"`
	MaxLevenshteinDistancePercentage float64     `mapstructure:"MAX_LEVENSHTEIN_DISTANCE_PERCENTAGE" validate:"required,min=0,max=1"`
}

func LoadEnv(v *validator.Validator) *Env {
	e := &Env{
		v: v,
	}

	if err := e.loadEnv(); err != nil {
		log.Fatalf("failed to load environment variables: %v", err)
	}

	return e
}

func (e *Env) loadEnv() error {
	viper.SetConfigType("env")
	err := viper.ReadConfig(bytes.NewBuffer(root.Env))
	if err != nil {
		return errs.New(err)
	}

	viper.AutomaticEnv()

	if err := viper.Unmarshal(&e); err != nil {
		return errs.New(err)
	}

	if err := e.validate(); err != nil {
		return errs.New(err)
	}

	return nil
}

func (e *Env) validate() error {
	if err := e.v.Validate(e); err != nil {
		return err
	}
	if e.Environment == "" {
		e.Environment = EnvironmentDevelopment
	}
	if e.Port == "" {
		e.Port = "8080"
	}
	return nil
}
