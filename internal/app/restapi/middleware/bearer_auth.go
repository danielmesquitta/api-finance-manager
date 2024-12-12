package middleware

import (
	"slices"
	"strings"
	"time"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/jwtutil"
	"github.com/labstack/echo/v4"
)

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

	bearerAuthOptions := []bearerAuthOption{
		withTokenType(jwtutil.TokenTypeAccess),
	}
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
			withTokenType(jwtutil.TokenTypeRefresh),
		)
	}
}

type bearerAuthOptions struct {
	TokenType jwtutil.TokenType
	Tiers     []entity.Tier
}

type bearerAuthOption func(*bearerAuthOptions)

func withTokenType(tokenType jwtutil.TokenType) bearerAuthOption {
	return func(o *bearerAuthOptions) {
		o.TokenType = tokenType
	}
}

func withTiers(tiers []entity.Tier) bearerAuthOption {
	return func(o *bearerAuthOptions) {
		o.Tiers = tiers
	}
}

func (m *Middleware) bearerAuth(
	c echo.Context,
	next echo.HandlerFunc,
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

	// Split the header to get the token part
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return errs.ErrUnauthorized
	}

	accessToken := parts[1]

	// Parse and validate the token
	claims, err := m.j.Parse(accessToken)
	if err != nil {
		return errs.ErrUnauthorized
	}

	if defaultOptions.TokenType != "" &&
		claims.TokenType != defaultOptions.TokenType {
		return errs.ErrUnauthorized
	}

	now := time.Now()
	if claims.SubscriptionExpiresAt != nil &&
		now.After(*claims.SubscriptionExpiresAt) {
		return errs.ErrSubscriptionExpired
	}

	if len(defaultOptions.Tiers) > 0 {
		if !slices.Contains(defaultOptions.Tiers, claims.Tier) {
			return errs.ErrUnauthorized
		}
	}

	// Set the claims in the context
	c.Set("claims", claims)

	// Token is valid, proceed with the request
	return next(c)
}
