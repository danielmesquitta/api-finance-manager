package pluggy

import (
	"github.com/danielmesquitta/api-finance-manager/internal/config"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/jwtutil"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/openfinance"
	"github.com/go-resty/resty/v2"
)

type Client struct {
	c *resty.Client
	e *config.Env
	j *jwtutil.JWT
}

func NewClient(
	e *config.Env,
	j *jwtutil.JWT,
) *Client {
	client := resty.New().
		SetBaseURL("https://api.pluggy.ai")

	return &Client{
		c: client,
		e: e,
		j: j,
	}
}

var _ openfinance.Client = (*Client)(nil)
