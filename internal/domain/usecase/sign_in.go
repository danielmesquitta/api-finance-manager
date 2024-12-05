package usecase

import (
	"context"
	"time"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/jwtutil"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/validator"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/oauth/googleoauth"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jinzhu/copier"
)

type SignInUseCase struct {
	v   validator.Validator
	ur  repo.UserRepo
	j   jwtutil.JWTManager
	goa *googleoauth.GoogleOAuth
}

func NewSignInUseCase(
	v validator.Validator,
	ur repo.UserRepo,
	j jwtutil.JWTManager,
	goa *googleoauth.GoogleOAuth,
) *SignInUseCase {
	return &SignInUseCase{
		v:   v,
		ur:  ur,
		j:   j,
		goa: goa,
	}
}

type SignInUseCaseInput struct {
	Provider string `json:"provider,omitempty" validate:"required,oneof=google apple"`
	Token    string `json:"token,omitempty"    validate:"required"`
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

	oauthUser, err := uc.goa.GetUser(in.Token)
	if err != nil {
		return nil, errs.New(err)
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

func (uc *SignInUseCase) signIn(
	ctx context.Context,
	user *entity.User,
) (*SignInUseCaseOutput, error) {
	select {
	case <-ctx.Done():
		return nil, errs.New(ctx.Err())
	default:
		accessToken, err := uc.j.NewToken(jwtutil.UserClaims{
			Tier:                  user.Tier,
			SubscriptionExpiresAt: user.SubscriptionExpiresAt,
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

	params.Tier = entity.TierTrial

	twoWeeksFromNow := time.Now().AddDate(0, 0, 14)
	params.SubscriptionExpiresAt = twoWeeksFromNow

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
