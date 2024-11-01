package config

import (
	"log"

	"github.com/danielmesquitta/api-finance-manager/internal/pkg/validator"
	"github.com/spf13/viper"
)

type Environment string

const (
	DevelopmentEnv Environment = "development"
	ProductionEnv  Environment = "production"
)

type Env struct {
	v validator.Validator

	Environment  Environment `mapstructure:"ENVIRONMENT"`
	Port         string      `mapstructure:"PORT"`
	DatabaseURL  string      `mapstructure:"DATABASE_URL"   validate:"required"`
	JWTSecretKey string      `mapstructure:"JWT_SECRET_KEY" validate:"required"`
}

func (e *Env) validate() error {
	if err := e.v.Validate(e); err != nil {
		return err
	}
	if e.Environment == "" {
		e.Environment = DevelopmentEnv
	}
	if e.Port == "" {
		e.Port = "8080"
	}
	return nil
}

func LoadEnv(v validator.Validator) *Env {
	env := &Env{
		v: v,
	}

	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalln(err)
	}

	if err := viper.Unmarshal(&env); err != nil {
		log.Fatalln(err)
	}

	if err := env.validate(); err != nil {
		log.Fatalln(err)
	}

	return env
}
