package jwtutil

import (
	"time"

	"github.com/danielmesquitta/api-finance-manager/internal/config"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/golang-jwt/jwt/v5"
)

type UserClaims struct {
	Tier                  entity.Tier `json:"tier,omitempty"`
	SubscriptionExpiresAt time.Time   `json:"subscription_expires_at,omitempty"`
	jwt.RegisteredClaims
}

type JWTManager interface {
	NewToken(claims UserClaims) (accessToken string, err error)
	ValidateToken(accessToken string) (*UserClaims, error)
}

type JWT struct {
	secretKey []byte
}

func NewJWT(
	env *config.Env,
) *JWT {
	return &JWT{
		secretKey: []byte(env.JWTSecretKey),
	}
}

func (j *JWT) NewToken(claims UserClaims) (string, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return accessToken.SignedString(j.secretKey)
}

func (j *JWT) ValidateToken(accessToken string) (*UserClaims, error) {
	parsedAccessToken, err := jwt.ParseWithClaims(
		accessToken,
		&UserClaims{},
		func(_ *jwt.Token) (interface{}, error) {
			return j.secretKey, nil
		},
	)
	if err != nil {
		return nil, errs.New(err)
	}

	userClaims, ok := parsedAccessToken.Claims.(*UserClaims)
	if !ok {
		return nil, errs.New("invalid claims")
	}

	if j.isExpired(&userClaims.RegisteredClaims) {
		return nil, errs.New("token is expired")
	}

	return userClaims, nil
}

func (j *JWT) isExpired(claims *jwt.RegisteredClaims) bool {
	return claims.ExpiresAt.Before(time.Now())
}

var _ JWTManager = &JWT{}
