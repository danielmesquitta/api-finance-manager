package pluggy

import (
	"context"
	"encoding/json"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
)

type authResponse struct {
	APIKey string `json:"apiKey"`
}

type authRequest struct {
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
}

const authHeaderKey = "x-api-key"

func (c *Client) refreshAccessToken(ctx context.Context) error {
	accessToken := c.c.Header.Get(authHeaderKey)

	if accessToken == "" {
		if err := c.authenticate(ctx); err != nil {
			return errs.New(err)
		}
		return nil
	}

	claims, err := c.j.Decode(accessToken)
	if err != nil {
		return errs.New(err)
	}

	if claims.IsExpired() {
		if err := c.authenticate(ctx); err != nil {
			return errs.New(err)
		}
		return nil
	}

	return nil
}

func (c *Client) authenticate(ctx context.Context) error {
	authRequest := authRequest{
		ClientID:     c.e.PluggyClientID,
		ClientSecret: c.e.PluggyClientSecret,
	}

	res, err := c.c.R().
		SetContext(ctx).
		SetBody(authRequest).
		Post("/auth")
	if err != nil {
		return errs.New(err)
	}
	body := res.Body()
	if statusCode := res.StatusCode(); statusCode < 200 || statusCode >= 300 {
		return errs.New(body)
	}

	data := authResponse{}
	if err := json.Unmarshal(res.Body(), &data); err != nil {
		return errs.New(err)
	}

	if data.APIKey == "" {
		return errs.New("api key is empty")
	}

	c.c.SetHeader(authHeaderKey, data.APIKey)

	return nil
}
