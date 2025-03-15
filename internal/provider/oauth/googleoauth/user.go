package googleoauth

import (
	"context"
	"encoding/json"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/gofiber/fiber/v2"
)

type UserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

func (g *GoogleOAuth) GetUser(
	ctx context.Context,
	token string,
) (*entity.User, *entity.UserAuthProvider, error) {
	res, err := g.c.R().
		SetContext(ctx).
		SetHeader(fiber.HeaderAuthorization, "Bearer "+token).
		SetQueryParam("alt", "json").
		Get("/oauth2/v1/userinfo")

	if err != nil {
		return nil, nil, errs.New(err)
	}

	body := res.Body()
	if res.IsError() {
		return nil, nil, errs.New(body)
	}

	userInfo := UserInfo{}
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, nil, errs.New(err)
	}

	var avatar *string
	if userInfo.Picture != "" {
		avatar = &userInfo.Picture
	}

	user := entity.User{
		Name:   userInfo.Name,
		Email:  userInfo.Email,
		Avatar: avatar,
	}

	authProvider := entity.UserAuthProvider{
		ExternalID:    userInfo.ID,
		VerifiedEmail: userInfo.VerifiedEmail,
		Provider:      entity.ProviderGoogle,
	}

	return &user, &authProvider, nil
}
