package googleoauth

import (
	"net/url"

	"github.com/danielmesquitta/api-finance-manager/internal/provider/oauth"
)

type GoogleOAuth struct {
	BaseURL url.URL
}

func NewGoogleOAuth() *GoogleOAuth {
	baseURL := url.URL{
		Scheme: "https",
		Host:   "googleapis.com",
	}

	return &GoogleOAuth{
		BaseURL: baseURL,
	}
}

var _ oauth.Provider = &GoogleOAuth{}
