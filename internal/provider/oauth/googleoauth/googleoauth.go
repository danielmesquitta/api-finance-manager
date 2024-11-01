package googleoauth

import (
	"github.com/danielmesquitta/api-finance-manager/internal/provider/oauth"
	"github.com/go-resty/resty/v2"
)

type GoogleOAuth struct {
	c *resty.Client
}

func NewGoogleOAuth() *GoogleOAuth {
	client := resty.New().
		SetBaseURL("https://googleapis.com")

	return &GoogleOAuth{
		c: client,
	}
}

var _ oauth.Provider = &GoogleOAuth{}
