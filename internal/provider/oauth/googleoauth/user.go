package googleoauth

import (
	"context"
	"encoding/json"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
)

const authorizationHeaderKey = "authorization"

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
) (*entity.User, error) {
	res, err := g.c.R().
		SetContext(ctx).
		SetHeader(authorizationHeaderKey, "Bearer "+token).
		SetQueryParam("alt", "json").
		Get("/oauth2/v1/userinfo")

	if err != nil {
		return nil, errs.New(err)
	}

	body := res.Body()
	if res.IsError() {
		return nil, errs.New(body)
	}

	userInfo := UserInfo{}
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, errs.New(err)
	}

	var avatar *string
	if userInfo.Picture != "" {
		avatar = &userInfo.Picture
	}

	user := entity.User{
		AuthID:        userInfo.ID,
		Name:          userInfo.Name,
		Email:         userInfo.Email,
		VerifiedEmail: userInfo.VerifiedEmail,
		Avatar:        avatar,
		Provider:      string(entity.ProviderGoogle),
	}

	return &user, nil
}
