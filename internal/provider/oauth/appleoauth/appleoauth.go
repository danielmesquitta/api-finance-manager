package appleoauth

import (
	"github.com/danielmesquitta/api-finance-manager/internal/provider/oauth"
	"github.com/go-resty/resty/v2"
)

type AppleOAuth struct {
	c *resty.Client
}

func NewAppleOAuth() *AppleOAuth {
	client := resty.New().
		SetBaseURL("https://googleapis.com")

	return &AppleOAuth{
		c: client,
	}
}

var _ oauth.Provider = &AppleOAuth{}
