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

	Environment                      Environment `mapstructure:"ENVIRONMENT"                         validate:"required,oneof=development production staging test"`
	Host                             string      `mapstructure:"HOST"`
	Port                             string      `mapstructure:"PORT"`
	PostgresDatabaseURL              string      `mapstructure:"POSTGRES_DATABASE_URL"               validate:"required"`
	RedisDatabaseURL                 string      `mapstructure:"REDIS_DATABASE_URL"                  validate:"required"`
	JWTAccessTokenSecretKey          string      `mapstructure:"JWT_ACCESS_TOKEN_SECRET_KEY"         validate:"required"`
	JWTRefreshTokenSecretKey         string      `mapstructure:"JWT_REFRESH_TOKEN_SECRET_KEY"        validate:"required"`
	HashSecretKey                    string      `mapstructure:"HASH_SECRET_KEY"                     validate:"required"`
	PluggyClientID                   string      `mapstructure:"PLUGGY_CLIENT_ID"                    validate:"required"`
	PluggyClientSecret               string      `mapstructure:"PLUGGY_CLIENT_SECRET"                validate:"required"`
	BasicAuthUsername                string      `mapstructure:"BASIC_AUTH_USERNAME"                 validate:"required"`
	BasicAuthPassword                string      `mapstructure:"BASIC_AUTH_PASSWORD"                 validate:"required"`
	MaxLevenshteinDistancePercentage float64     `mapstructure:"MAX_LEVENSHTEIN_DISTANCE_PERCENTAGE" validate:"required,min=0,max=1"`
	SyncBalancesMaxAccounts          int         `mapstructure:"SYNC_BALANCES_MAX_ACCOUNTS"          validate:"required,min=1"`
	SyncTransactionsMaxAccounts      int         `mapstructure:"SYNC_TRANSACTIONS_MAX_ACCOUNTS"      validate:"required,min=1"`
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
