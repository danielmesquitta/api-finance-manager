package pluggy

import (
	"context"
	"encoding/json"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/money"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/openfinance"
)

type accountsResponse struct {
	Results []accountsResult `json:"results"`
}

type accountsResult struct {
	ID      string  `json:"id"`
	Type    string  `json:"type"`
	Name    string  `json:"name"`
	Balance float64 `json:"balance"`
}

func (c *Client) ListAccounts(
	ctx context.Context,
	connectionID string,
) ([]openfinance.Account, error) {
	if err := c.refreshAccessToken(ctx); err != nil {
		return nil, errs.New(err)
	}

	queryParams := map[string]string{
		"itemId": connectionID,
	}

	res, err := c.c.R().
		SetContext(ctx).
		SetQueryParams(queryParams).
		Get("/accounts")
	if err != nil {
		return nil, errs.New(err)
	}
	body := res.Body()
	if res.IsError() {
		return nil, errs.New(body)
	}

	accountsRes := accountsResponse{}
	if err := json.Unmarshal(body, &accountsRes); err != nil {
		return nil, errs.New(err)
	}

	var accounts []openfinance.Account
	for _, a := range accountsRes.Results {
		accounts = append(accounts, openfinance.Account{
			Account: entity.Account{
				ExternalID: a.ID,
				Type:       a.Type,
				Name:       a.Name,
			},
			Balance: money.ToCents(a.Balance),
		})
	}

	return accounts, nil
}
