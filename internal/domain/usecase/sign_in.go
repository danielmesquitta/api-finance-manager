package usecase

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/jwtutil"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/validator"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/oauth/googleoauth"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/oauth/mockoauth"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

type SignIn struct {
	v  *validator.Validator
	ur repo.UserRepo
	j  *jwtutil.JWT
	g  *googleoauth.GoogleOAuth
	m  *mockoauth.MockOAuth
}

func NewSignIn(
	v *validator.Validator,
	ur repo.UserRepo,
	j *jwtutil.JWT,
	g *googleoauth.GoogleOAuth,
	m *mockoauth.MockOAuth,
) *SignIn {
	return &SignIn{
		v:  v,
		ur: ur,
		j:  j,
		g:  g,
		m:  m,
	}
}

type SignInInput struct {
	Provider entity.Provider `json:"provider" validate:"required,oneof=GOOGLE APPLE REFRESH MOCK"`
	Token    string          `json:"-"        validate:"required_without=UserID"`
	UserID   uuid.UUID       `json:"-"        validate:"required_without=Token"`
}

type SignInOutput struct {
	User         entity.User `json:"user"`
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
}

func (uc *SignIn) Execute(
	ctx context.Context,
	in SignInInput,
) (*SignInOutput, error) {
	if err := uc.v.Validate(in); err != nil {
		return nil, errs.New(err)
	}

	switch in.Provider {
	case entity.ProviderGoogle:
		return uc.signInWithGoogle(ctx, in.Token)
	case entity.ProviderRefresh:
		return uc.refreshToken(ctx, in.UserID)
	case entity.ProviderMock:
		return uc.signInWithMock(ctx, in.Token)
	default:
		return nil, errs.New("Provedor não implementado")
	}
}

func (uc *SignIn) signInWithGoogle(
	ctx context.Context,
	token string,
) (*SignInOutput, error) {
	token = strings.TrimPrefix(token, "Bearer ")

	oauthUser, err := uc.g.GetUser(ctx, token)
	if err != nil {
		slog.Error("failed to get user from google", "error", err)
		return nil, errs.ErrUnauthorized
	}

	registeredUser, err := uc.ur.GetUserByEmail(ctx, oauthUser.Email)
	if err != nil {
		return nil, errs.New(err)
	}

	if registeredUser != nil {
		updatedUser, err := uc.updateUser(ctx, registeredUser, oauthUser)
		if err != nil {
			return nil, errs.New(err)
		}

		return uc.signIn(ctx, updatedUser)
	}

	user, err := uc.createUser(ctx, oauthUser)
	if err != nil {
		return nil, errs.New(err)
	}

	return uc.signIn(ctx, user)
}

func (uc *SignIn) signInWithMock(
	ctx context.Context,
	token string,
) (*SignInOutput, error) {
	oauthUser, err := uc.m.GetUser(ctx, token)
	if err != nil {
		slog.Error("failed to get user from mock", "error", err)
		return nil, errs.ErrUnauthorized
	}

	registeredUser, err := uc.ur.GetUserByEmail(ctx, oauthUser.Email)
	if err != nil {
		return nil, errs.New(err)
	}

	if registeredUser != nil {
		updatedUser, err := uc.updateUser(ctx, registeredUser, oauthUser)
		if err != nil {
			return nil, errs.New(err)
		}

		return uc.signIn(ctx, updatedUser)
	}

	user, err := uc.createUser(ctx, oauthUser)
	if err != nil {
		return nil, errs.New(err)
	}

	return uc.signIn(ctx, user)
}

func (uc *SignIn) refreshToken(
	ctx context.Context,
	userID uuid.UUID,
) (*SignInOutput, error) {
	user, err := uc.ur.GetUserByID(ctx, userID)
	if err != nil {
		return nil, errs.New(err)
	}
	if user == nil {
		return nil, errs.ErrUserNotFound
	}

	return uc.signIn(ctx, user)
}

func (uc *SignIn) signIn(
	ctx context.Context,
	user *entity.User,
) (*SignInOutput, error) {
	select {
	case <-ctx.Done():
		return nil, errs.New(ctx.Err())
	default:
	}

	tokenClaims := jwtutil.UserClaims{}
	tokenClaims.Issuer = user.ID.String()
	tokenClaims.IssuedAt = jwt.NewNumericDate(time.Now())
	in7Days := time.Now().Add(time.Hour * 24 * 7)
	tokenClaims.ExpiresAt = jwt.NewNumericDate(in7Days)

	refreshToken, err := uc.j.NewToken(tokenClaims, jwtutil.TokenTypeRefresh)
	if err != nil {
		return nil, errs.New(err)
	}

	tokenClaims.Tier = entity.Tier(user.Tier)
	tokenClaims.SubscriptionExpiresAt = user.SubscriptionExpiresAt
	in24Hours := time.Now().Add(time.Hour * 24)
	tokenClaims.ExpiresAt = jwt.NewNumericDate(in24Hours)

	accessToken, err := uc.j.NewToken(
		tokenClaims,
		jwtutil.TokenTypeAccess,
	)
	if err != nil {
		return nil, errs.New(err)
	}

	return &SignInOutput{
		User:         *user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (uc *SignIn) createUser(
	ctx context.Context,
	oauthUser *entity.User,
) (*entity.User, error) {
	params := repo.CreateUserParams{
		Tier: string(entity.TierFree),
	}
	if err := copier.CopyWithOption(
		&params,
		oauthUser,
		copier.Option{IgnoreEmpty: true},
	); err != nil {
		return nil, errs.New(err)
	}

	user, err := uc.ur.CreateUser(ctx, params)
	if err != nil {
		return nil, errs.New(err)
	}

	return user, nil
}

func (uc *SignIn) updateUser(
	ctx context.Context,
	user *entity.User,
	oauthUser *entity.User,
) (*entity.User, error) {
	params := repo.UpdateUserParams{}
	if err := copier.Copy(&params, user); err != nil {
		return nil, errs.New(err)
	}

	if err := copier.CopyWithOption(
		&params,
		oauthUser,
		copier.Option{IgnoreEmpty: true},
	); err != nil {
		return nil, errs.New(err)
	}

	updatedUser, err := uc.ur.UpdateUser(ctx, params)
	if err != nil {
		return nil, errs.New(err)
	}

	return updatedUser, nil
}
