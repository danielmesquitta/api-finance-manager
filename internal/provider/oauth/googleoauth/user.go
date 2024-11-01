package googleoauth

import (
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
	token string,
) (*entity.User, error) {
	res, err := g.c.R().
		SetHeader(authorizationHeaderKey, "Bearer "+token).
		Get("/userinfo/v2/me")

	if err != nil {
		return nil, errs.New(err)
	}

	body := res.Body()
	if statusCode := res.StatusCode(); statusCode < 200 || statusCode >= 300 {
		return nil, errs.New(body)
	}

	userInfo := UserInfo{}
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, errs.New(err)
	}

	user := entity.User{
		Name:   userInfo.Name,
		Email:  userInfo.Email,
		Avatar: &userInfo.Picture,
	}

	return &user, nil
}
