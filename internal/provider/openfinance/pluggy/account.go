package pluggy

import (
	"context"
	"encoding/json"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/google/uuid"
)

type accountsResponse struct {
	Results []accountsResult `json:"results"`
}

type accountsResult struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	Name string `json:"name"`
}

func (c *Client) ListAccounts(
	ctx context.Context,
	connectionID uuid.UUID,
) ([]entity.Account, error) {
	if err := c.refreshAccessToken(ctx); err != nil {
		return nil, errs.New(err)
	}

	queryParams := map[string]string{
		"itemId": connectionID.String(),
	}

	res, err := c.c.R().
		SetContext(ctx).
		SetQueryParams(queryParams).
		Get("/accounts")
	if err != nil {
		return nil, errs.New(err)
	}
	if err := res.Error(); err != nil {
		return nil, errs.New(err)
	}

	accountsRes := accountsResponse{}
	if err := json.Unmarshal(res.Body(), &accountsRes); err != nil {
		return nil, errs.New(err)
	}

	var accounts []entity.Account
	for _, a := range accountsRes.Results {
		accounts = append(accounts, entity.Account{
			ExternalID: a.ID,
			Type:       a.Type,
			Name:       a.Name,
		})
	}

	return accounts, nil
}