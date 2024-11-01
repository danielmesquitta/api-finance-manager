package usecase

import (
	"context"
	"time"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/jwtutil"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/validator"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/oauth"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jinzhu/copier"
)

type SignInUseCase struct {
	v  validator.Validator
	ur repo.UserRepo
	op oauth.Provider
	j  jwtutil.JWTManager
}

func NewSignInUseCase(
	v validator.Validator,
	ur repo.UserRepo,
	op oauth.Provider,
	j jwtutil.JWTManager,
) *SignInUseCase {
	return &SignInUseCase{
		v:  v,
		ur: ur,
		op: op,
		j:  j,
	}
}

type SignInUseCaseInput struct {
	Token string `json:"token,omitempty" validate:"required,jwt"`
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

	oauthUser, err := uc.op.GetUser(in.Token)
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
			Tier:                  user.Tier,
			SubscriptionExpiresAt: user.SubscriptionExpiresAt,
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
	oauthUser *oauth.User,
) (*entity.User, error) {
	params := repo.CreateUserParams{
		Name:   oauthUser.Name,
		Email:  oauthUser.Email,
		Tier:   entity.TierTRIAL,
		Avatar: &oauthUser.Picture,
	}

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
	oauthUser *oauth.User,
) (*entity.User, error) {
	params := repo.UpdateUserParams{}
	if err := copier.Copy(&params, user); err != nil {
		return nil, errs.New(err)
	}

	params.Name = oauthUser.Name
	params.Email = oauthUser.Email
	params.Avatar = &oauthUser.Picture

	updatedUser, err := uc.ur.UpdateUser(ctx, params)
	if err != nil {
		return nil, errs.New(err)
	}

	return updatedUser, nil
}
