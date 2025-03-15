package usecase

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/hash"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/jwtutil"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/tx"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/validator"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/oauth"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/oauth/googleoauth"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/oauth/mockoauth"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

type SignIn struct {
	v    *validator.Validator
	tx   tx.TX
	h    *hash.Hasher
	ur   repo.UserRepo
	uapr repo.UserAuthProviderRepo
	j    *jwtutil.JWT
	g    *googleoauth.GoogleOAuth
	m    *mockoauth.MockOAuth
}

func NewSignIn(
	v *validator.Validator,
	tx tx.TX,
	h *hash.Hasher,
	ur repo.UserRepo,
	uapr repo.UserAuthProviderRepo,
	j *jwtutil.JWT,
	g *googleoauth.GoogleOAuth,
	m *mockoauth.MockOAuth,
) *SignIn {
	return &SignIn{
		v:    v,
		tx:   tx,
		h:    h,
		ur:   ur,
		uapr: uapr,
		j:    j,
		g:    g,
		m:    m,
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

	if in.Provider == entity.ProviderRefresh {
		return uc.refreshToken(ctx, in.UserID)
	}

	providerMap := map[entity.Provider]oauth.Provider{
		entity.ProviderGoogle: uc.g,
		entity.ProviderMock:   uc.m,
	}

	oauthProvider, ok := providerMap[in.Provider]
	if !ok {
		return nil, errs.ErrInvalidProvider
	}

	token := strings.TrimPrefix(in.Token, "Bearer ")

	authUser, authProvider, err := oauthProvider.GetUser(ctx, token)
	if err != nil {
		slog.Error(
			"sign-in: failed to get user from auth provider",
			"provider", in.Provider,
			"error", err,
		)
		return nil, errs.ErrUnauthorized
	}

	registeredUser, err := uc.ur.GetUserByEmail(ctx, authUser.Email)
	if err != nil {
		return nil, errs.New(err)
	}

	if registeredUser != nil {
		return uc.updateUserAndSignIn(
			ctx,
			registeredUser,
			authUser,
			authProvider,
		)
	}

	hashedEmail, err := uc.h.Hash(authUser.Email)
	if err != nil {
		return nil, errs.New(err)
	}

	deletedUser, err := uc.ur.GetUserByEmail(ctx, hashedEmail)
	if err != nil {
		return nil, errs.New(err)
	}

	return uc.createUserAndSignIn(ctx, deletedUser, authUser, authProvider)
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
	in7Days := time.Now().Add(time.Hour * 24 * 7)
	tokenClaims := jwtutil.UserClaims{
		Issuer:    user.ID.String(),
		IssuedAt:  time.Now(),
		ExpiresAt: in7Days,
	}

	refreshToken, err := uc.j.NewToken(tokenClaims, jwtutil.TokenTypeRefresh)
	if err != nil {
		return nil, errs.New(err)
	}

	tokenClaims.Tier = entity.Tier(user.Tier)
	tokenClaims.SubscriptionExpiresAt = user.SubscriptionExpiresAt
	in24Hours := time.Now().Add(time.Hour * 24)
	tokenClaims.ExpiresAt = in24Hours

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

func (uc *SignIn) createUserAndSignIn(
	ctx context.Context,
	deletedUser *entity.User,
	authUser *entity.User,
	authProvider *entity.UserAuthProvider,
) (*SignInOutput, error) {
	var out *SignInOutput
	err := uc.tx.Do(ctx, func(ctx context.Context) error {
		if deletedUser != nil {
			err := uc.ur.DestroyUser(ctx, deletedUser.ID)
			if err != nil {
				return errs.New(err)
			}
		}

		user, err := uc.createUser(ctx, authUser)
		if err != nil {
			return errs.New(err)
		}

		if err := uc.createUserAuthProvider(ctx, user.ID, authProvider); err != nil {
			return errs.New(err)
		}

		out, err = uc.signIn(ctx, user)
		if err != nil {
			return errs.New(err)
		}

		return nil
	})
	if err != nil {
		return nil, errs.New(err)
	}

	return out, nil
}

func (uc *SignIn) updateUserAndSignIn(
	ctx context.Context,
	registeredUser *entity.User,
	authUser *entity.User,
	authProvider *entity.UserAuthProvider,
) (*SignInOutput, error) {
	registeredAuthProvider, err := uc.uapr.GetUserAuthProvider(
		ctx,
		repo.GetUserAuthProviderParams{
			UserID:   registeredUser.ID,
			Provider: authProvider.Provider,
		},
	)
	if err != nil {
		return nil, errs.New(err)
	}

	var out *SignInOutput
	err = uc.tx.Do(ctx, func(ctx context.Context) error {
		updatedUser, err := uc.updateUser(ctx, registeredUser, authUser)
		if err != nil {
			return errs.New(err)
		}

		if registeredAuthProvider == nil {
			if err := uc.createUserAuthProvider(ctx, registeredUser.ID, authProvider); err != nil {
				return errs.New(err)
			}
		} else {
			if err := uc.updateUserAuthProvider(ctx, registeredAuthProvider, authProvider); err != nil {
				return errs.New(err)
			}
		}

		out, err = uc.signIn(ctx, updatedUser)
		if err != nil {
			return errs.New(err)
		}

		return nil
	})

	if err != nil {
		return nil, errs.New(err)
	}

	return out, nil
}

func (uc *SignIn) createUser(
	ctx context.Context,
	authUser *entity.User,
) (*entity.User, error) {
	params := repo.CreateUserParams{
		Tier: string(entity.TierFree),
	}
	if err := copier.CopyWithOption(
		&params,
		authUser,
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
	authUser *entity.User,
) (*entity.User, error) {
	params := repo.UpdateUserParams{}
	if err := copier.Copy(&params, user); err != nil {
		return nil, errs.New(err)
	}

	if err := copier.CopyWithOption(
		&params,
		authUser,
		copier.Option{IgnoreEmpty: true},
	); err != nil {
		return nil, errs.New(err)
	}

	// do not update user's
	// email and name if it is already set
	if user.Name != "" {
		params.Name = user.Name
	}

	if user.Email != "" {
		params.Email = user.Email
	}

	updatedUser, err := uc.ur.UpdateUser(ctx, params)
	if err != nil {
		return nil, errs.New(err)
	}

	return updatedUser, nil
}

func (uc *SignIn) createUserAuthProvider(
	ctx context.Context,
	userID uuid.UUID,
	authProvider *entity.UserAuthProvider,
) error {
	params := repo.CreateUserAuthProviderParams{}
	if err := copier.Copy(&params, authProvider); err != nil {
		return errs.New(err)
	}
	params.UserID = userID

	if err := uc.uapr.CreateUserAuthProvider(ctx, params); err != nil {
		return errs.New(err)
	}

	return nil
}

func (uc *SignIn) updateUserAuthProvider(
	ctx context.Context,
	registeredAuthProvider *entity.UserAuthProvider,
	authProvider *entity.UserAuthProvider,
) error {
	params := repo.UpdateUserAuthProviderParams{}
	if err := copier.Copy(&params, authProvider); err != nil {
		return errs.New(err)
	}
	params.ID = registeredAuthProvider.ID

	if err := uc.uapr.UpdateUserAuthProvider(ctx, params); err != nil {
		return errs.New(err)
	}

	return nil
}
