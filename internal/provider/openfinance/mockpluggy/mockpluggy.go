package mockpluggy

import (
	"github.com/danielmesquitta/api-finance-manager/internal/provider/openfinance"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/openfinance/pluggy"
)

type Client struct {
	*pluggy.Client
}

func NewClient(
	c *pluggy.Client,
) *Client {
	return &Client{
		Client: c,
	}
}

var _ openfinance.Client = (*Client)(nil)
