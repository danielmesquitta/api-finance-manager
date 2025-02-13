package config

import (
	"bytes"
	_ "embed"

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

	Environment                      Environment `validate:"required,oneof=development production staging test" mapstructure:"ENVIRONMENT"`
	Host                             string      `                                                              mapstructure:"HOST"`
	Port                             string      `                                                              mapstructure:"PORT"`
	PostgresDatabaseURL              string      `validate:"required"                                           mapstructure:"POSTGRES_DATABASE_URL"`
	RedisDatabaseURL                 string      `validate:"required"                                           mapstructure:"REDIS_DATABASE_URL"`
	JWTAccessTokenSecretKey          string      `validate:"required"                                           mapstructure:"JWT_ACCESS_TOKEN_SECRET_KEY"`
	JWTRefreshTokenSecretKey         string      `validate:"required"                                           mapstructure:"JWT_REFRESH_TOKEN_SECRET_KEY"`
	PluggyClientID                   string      `validate:"required"                                           mapstructure:"PLUGGY_CLIENT_ID"`
	PluggyClientSecret               string      `validate:"required"                                           mapstructure:"PLUGGY_CLIENT_SECRET"`
	BasicAuthUsername                string      `validate:"required"                                           mapstructure:"BASIC_AUTH_USERNAME"`
	BasicAuthPassword                string      `validate:"required"                                           mapstructure:"BASIC_AUTH_PASSWORD"`
	MaxLevenshteinDistancePercentage float64     `validate:"required,min=0,max=1"                               mapstructure:"MAX_LEVENSHTEIN_DISTANCE_PERCENTAGE"`
	SyncBalancesMaxAccounts          int         `validate:"required,min=1"                                     mapstructure:"SYNC_BALANCES_MAX_ACCOUNTS"`
	SyncTransactionsMaxAccounts      int         `validate:"required,min=1"                                     mapstructure:"SYNC_TRANSACTIONS_MAX_ACCOUNTS"`
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
	if e.Host == "" {
		e.Host = "http://localhost"
	}
	return nil
}
