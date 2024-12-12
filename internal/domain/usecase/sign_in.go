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

type SignInUseCase struct {
	v  validator.Validator
	ur repo.UserRepo
	j  jwtutil.JWTManager
	g  *googleoauth.GoogleOAuth
	m  *mockoauth.MockOAuth
}

func NewSignInUseCase(
	v validator.Validator,
	ur repo.UserRepo,
	j jwtutil.JWTManager,
	g *googleoauth.GoogleOAuth,
	m *mockoauth.MockOAuth,
) *SignInUseCase {
	return &SignInUseCase{
		v:  v,
		ur: ur,
		j:  j,
		g:  g,
		m:  m,
	}
}

type SignInUseCaseInput struct {
	Provider entity.Provider `json:"provider,omitempty" validate:"required,oneof=GOOGLE APPLE REFRESH MOCK"`
	Token    string          `json:"token,omitempty"    validate:"required_without=UserID"`
	UserID   uuid.UUID       `json:"-"                  validate:"required_without=Token"`
}

type SignInUseCaseOutput struct {
	User         entity.User `json:"user,omitempty"`
	AccessToken  string      `json:"access_token,omitempty"`
	RefreshToken string      `json:"refresh_token,omitempty"`
}

func (uc *SignInUseCase) Execute(
	ctx context.Context,
	in SignInUseCaseInput,
) (*SignInUseCaseOutput, error) {
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

func (uc *SignInUseCase) signInWithGoogle(
	ctx context.Context,
	token string,
) (*SignInUseCaseOutput, error) {
	token = strings.TrimPrefix(token, "Bearer ")

	oauthUser, err := uc.g.GetUser(token)
	if err != nil {
		slog.Info("failed to get user from google", "error", err)
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

func (uc *SignInUseCase) signInWithMock(
	ctx context.Context,
	token string,
) (*SignInUseCaseOutput, error) {
	oauthUser, err := uc.m.GetUser(token)
	if err != nil {
		slog.Info("failed to get user from mock", "error", err)
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

func (uc *SignInUseCase) refreshToken(
	ctx context.Context,
	userID uuid.UUID,
) (*SignInUseCaseOutput, error) {
	user, err := uc.ur.GetUserByID(ctx, userID)
	if err != nil {
		return nil, errs.New(err)
	}
	if user == nil {
		return nil, errs.ErrUserNotFound
	}

	return uc.signIn(ctx, user)
}

func (uc *SignInUseCase) signIn(
	ctx context.Context,
	user *entity.User,
) (*SignInUseCaseOutput, error) {
	select {
	case <-ctx.Done():
		return nil, errs.New(ctx.Err())
	default:
	}

	accessToken, err := uc.j.NewToken(jwtutil.UserClaims{
		Tier:                  entity.Tier(user.Tier),
		SubscriptionExpiresAt: user.SubscriptionExpiresAt,
		TokenType:             jwtutil.TokenTypeAccess,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    user.ID.String(),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	})
	if err != nil {
		return nil, errs.New(err)
	}

	refreshToken, err := uc.j.NewToken(jwtutil.UserClaims{
		Tier:                  entity.Tier(user.Tier),
		SubscriptionExpiresAt: user.SubscriptionExpiresAt,
		TokenType:             jwtutil.TokenTypeRefresh,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:   user.ID.String(),
			IssuedAt: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(
				time.Now().Add(time.Hour * 24 * 7),
			),
		},
	})
	if err != nil {
		return nil, errs.New(err)
	}

	return &SignInUseCaseOutput{
		User:         *user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (uc *SignInUseCase) createUser(
	ctx context.Context,
	oauthUser *entity.User,
) (*entity.User, error) {
	params := repo.CreateUserParams{}
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

func (uc *SignInUseCase) updateUser(
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
