package pluggy

import (
	"context"
	"encoding/json"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/openfinance"
)

func (c *Client) GetAccount(
	ctx context.Context,
	accountID string,
) (*openfinance.Account, error) {
	if err := c.refreshAccessToken(ctx); err != nil {
		return nil, errs.New(err)
	}

	res, err := c.c.R().
		SetContext(ctx).
		Post("/accounts/" + accountID)
	if err != nil {
		return nil, errs.New(err)
	}
	body := res.Body()
	if statusCode := res.StatusCode(); statusCode < 200 || statusCode >= 300 {
		return nil, errs.New(body)
	}

	account := openfinance.Account{}
	if err := json.Unmarshal(body, &account); err != nil {
		return nil, errs.New(err)
	}

	return &account, nil
}
