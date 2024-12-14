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
	EnvDevelopment Environment = "development"
	EnvProduction  Environment = "production"
	EnvStaging     Environment = "staging"
	EnvTest        Environment = "test"
)

type Env struct {
	v validator.Validator

	Environment        Environment `mapstructure:"ENVIRONMENT"`
	Port               string      `mapstructure:"PORT"`
	DatabaseURL        string      `mapstructure:"DATABASE_URL"         validate:"required"`
	JWTSecretKey       string      `mapstructure:"JWT_SECRET_KEY"       validate:"required"`
	PluggyClientID     string      `mapstructure:"PLUGGY_CLIENT_ID"     validate:"required"`
	PluggyClientSecret string      `mapstructure:"PLUGGY_CLIENT_SECRET" validate:"required"`
	BasicAuthUsername  string      `mapstructure:"BASIC_AUTH_USERNAME"  validate:"required"`
	BasicAuthPassword  string      `mapstructure:"BASIC_AUTH_PASSWORD"  validate:"required"`
}

func LoadEnv(v validator.Validator) *Env {
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
		e.Environment = EnvDevelopment
	}
	if e.Port == "" {
		e.Port = "8080"
	}
	return nil
}
