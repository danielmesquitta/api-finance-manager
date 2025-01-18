package middleware

import (
	"regexp"
	"slices"
	"time"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/jwtutil"
	"github.com/labstack/echo/v4"
)

const ClaimsKey = "claims"

type BearerAuthOptions struct {
	Tiers []entity.Tier
}

type BearerAuthOption func(*BearerAuthOptions)

func WithTiers(tiers []entity.Tier) BearerAuthOption {
	return func(o *BearerAuthOptions) {
		o.Tiers = tiers
	}
}

func (m *Middleware) BearerAuthAccessToken(
	options ...BearerAuthOption,
) echo.MiddlewareFunc {
	defaultOptions := BearerAuthOptions{}
	for _, opt := range options {
		opt(&defaultOptions)
	}

	bearerAuthOptions := []bearerAuthOption{}
	if len(defaultOptions.Tiers) > 0 {
		bearerAuthOptions = append(
			bearerAuthOptions,
			withTiers(defaultOptions.Tiers),
		)
	}

	return func(
		next echo.HandlerFunc,
	) echo.HandlerFunc {
		return func(c echo.Context) error {
			return m.bearerAuth(
				c,
				next,
				jwtutil.TokenTypeAccess,
				bearerAuthOptions...,
			)
		}
	}
}

func (m *Middleware) BearerAuthRefreshToken(
	next echo.HandlerFunc,
) echo.HandlerFunc {
	return func(c echo.Context) error {
		return m.bearerAuth(
			c,
			next,
			jwtutil.TokenTypeRefresh,
		)
	}
}

type bearerAuthOptions struct {
	TokenType jwtutil.TokenType
	Tiers     []entity.Tier
}

type bearerAuthOption func(*bearerAuthOptions)

func withTiers(tiers []entity.Tier) bearerAuthOption {
	return func(o *bearerAuthOptions) {
		o.Tiers = tiers
	}
}

func (m *Middleware) bearerAuth(
	c echo.Context,
	next echo.HandlerFunc,
	tokenType jwtutil.TokenType,
	options ...bearerAuthOption,
) error {
	defaultOptions := bearerAuthOptions{}
	for _, opt := range options {
		opt(&defaultOptions)
	}

	// Get the Authorization header
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return errs.ErrUnauthorized
	}

	// Extract the token from the header
	re := regexp.MustCompile(`(?i)^Bearer\s+`)
	accessToken := re.ReplaceAllString(authHeader, "")

	// Parse and validate the token
	claims, err := m.j.Parse(accessToken, tokenType)
	if err != nil {
		return errs.ErrUnauthorized
	}

	// Block if subscription is expired
	now := time.Now()
	if claims.SubscriptionExpiresAt != nil &&
		now.After(*claims.SubscriptionExpiresAt) {
		return errs.ErrSubscriptionExpired
	}

	// Block if not allowed tier
	if len(defaultOptions.Tiers) > 0 {
		if !slices.Contains(defaultOptions.Tiers, claims.Tier) {
			return errs.ErrUnauthorized
		}
	}

	// Set the claims in the context
	c.Set(ClaimsKey, claims)

	// Token is valid, proceed with the request
	return next(c)
}
