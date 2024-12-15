package jwtutil

import (
	"encoding/base64"
	"encoding/json"
	"strings"
	"time"

	"github.com/danielmesquitta/api-finance-manager/internal/config"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/golang-jwt/jwt/v5"
)

type TokenType byte

const (
	TokenTypeAccess TokenType = iota
	TokenTypeRefresh
)

type RegisteredClaims struct {
	jwt.RegisteredClaims
}

func (rc *RegisteredClaims) IsExpired() bool {
	nowMinusOneMinute := time.Now().Add(-1 * time.Minute)
	return rc.ExpiresAt.Before(nowMinusOneMinute)
}

type UserClaims struct {
	RegisteredClaims
	Tier                  entity.Tier `json:"tier,omitempty"`
	SubscriptionExpiresAt *time.Time  `json:"subscription_expires_at,omitempty"`
}

// IsExpired checks if the token is expired
// by comparing the expiration time with the current time minus one minute
// to account for requests that may take longer to process

type JWT struct {
	keys map[TokenType][]byte
}

func NewJWT(
	e *config.Env,
) *JWT {
	keys := map[TokenType][]byte{
		TokenTypeAccess:  []byte(e.JWTAccessTokenSecretKey),
		TokenTypeRefresh: []byte(e.JWTRefreshTokenSecretKey),
	}

	return &JWT{
		keys: keys,
	}
}

func (j *JWT) NewToken(claims UserClaims, tokenType TokenType) (string, error) {
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return jwtToken.SignedString(j.keys[tokenType])
}

func (j *JWT) Parse(jwtToken string, tokenType TokenType) (*UserClaims, error) {
	parsedAccessToken, err := jwt.ParseWithClaims(
		jwtToken,
		&UserClaims{},
		func(_ *jwt.Token) (interface{}, error) {
			return j.keys[tokenType], nil
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
		return nil, errs.ErrUnauthorized
	}

	return userClaims, nil
}

// Decode decodes a JWT and extracts the payload.
func (j *JWT) Decode(token string) (*RegisteredClaims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, errs.New("invalid token")
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, errs.New(err)
	}

	var claims RegisteredClaims
	if err := json.Unmarshal(payload, &claims); err != nil {
		return nil, errs.New(err)
	}

	return &claims, nil
}
