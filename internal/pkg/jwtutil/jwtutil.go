package jwtutil

import (
	"time"

	"github.com/danielmesquitta/api-finance-manager/internal/config"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/golang-jwt/jwt/v5"
)

type UserClaims struct {
	jwt.RegisteredClaims
	Tier                  entity.Tier `json:"tier,omitempty"`
	SubscriptionExpiresAt time.Time   `json:"subscription_expires_at,omitempty"`
}

// IsExpired checks if the token is expired
// by comparing the expiration time with the current time minus one minute
// to account for requests that may take longer to process
func (u *UserClaims) IsExpired() bool {
	nowMinusOneMinute := time.Now().Add(-1 * time.Minute)
	return u.RegisteredClaims.ExpiresAt.Before(nowMinusOneMinute)
}

type JWTManager interface {
	NewToken(claims UserClaims) (jwtToken string, err error)
	Parse(jwtToken string) (*UserClaims, error)
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
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return jwtToken.SignedString(j.secretKey)
}

func (j *JWT) Parse(jwtToken string) (*UserClaims, error) {
	parsedAccessToken, err := jwt.ParseWithClaims(
		jwtToken,
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

	if userClaims.IsExpired() {
		return nil, errs.New("token is expired")
	}

	return userClaims, nil
}

var _ JWTManager = &JWT{}
