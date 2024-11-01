package config

import (
	"log"
	"log/slog"
	"os"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/validator"
	"github.com/spf13/viper"
)

type Environment string

const (
	EnvDevelopment Environment = "development"
	EnvProduction  Environment = "production"
)

type Env struct {
	v validator.Validator

	Environment        Environment `mapstructure:"ENVIRONMENT"`
	Port               string      `mapstructure:"PORT"`
	DatabaseURL        string      `mapstructure:"DATABASE_URL"         validate:"required"`
	JWTSecretKey       string      `mapstructure:"JWT_SECRET_KEY"       validate:"required"`
	PluggyClientID     string      `mapstructure:"PLUGGY_CLIENT_ID"     validate:"required"`
	PluggyClientSecret string      `mapstructure:"PLUGGY_CLIENT_SECRET" validate:"required"`
}

func LoadEnv(v validator.Validator) *Env {
	env := &Env{
		v: v,
	}

	if err := env.loadEnv(); err != nil {
		log.Fatalf("failed to load environment variables: %v", err)
	}

	env.configLogger()

	return env
}

func (e *Env) loadEnv() error {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return errs.New(err)
	}

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

func (e *Env) configLogger() {
	if e.Environment == EnvProduction {
		slogger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
		slog.SetDefault(slogger)
	}
}
