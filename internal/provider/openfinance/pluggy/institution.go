package pluggy

import (
	"context"
	"encoding/json"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/openfinance"
)

type consentsResponse struct {
	Results []result `json:"results"`
}

type result struct {
	ItemID string `json:"itemId"`
}

func (c *Client) GetInstitution(
	ctx context.Context,
	accountItemID string,
) (*openfinance.Institution, error) {
	if err := c.refreshAccessToken(ctx); err != nil {
		return nil, errs.New(err)
	}

	res, err := c.c.R().
		SetContext(ctx).
		SetQueryParam("item_id", accountItemID).
		Get("/consents")
	if err != nil {
		return nil, errs.New(err)
	}
	body := res.Body()
	if statusCode := res.StatusCode(); statusCode < 200 || statusCode >= 300 {
		return nil, errs.New(body)
	}

	consents := consentsResponse{}
	if err := json.Unmarshal(body, &consents); err != nil {
		return nil, errs.New(err)
	}

	if len(consents.Results) == 0 {
		return nil, errs.New("no consents found")
	}

	institutionID := consents.Results[0].ItemID

	res, err = c.c.R().
		SetContext(ctx).
		Get("/items/" + institutionID)
	if err != nil {
		return nil, errs.New(err)
	}
	body = res.Body()
	if statusCode := res.StatusCode(); statusCode < 200 || statusCode >= 300 {
		return nil, errs.New(body)
	}

	institution := openfinance.Institution{}
	if err := json.Unmarshal(body, &institution); err != nil {
		return nil, errs.New(err)
	}

	return &institution, nil
}
