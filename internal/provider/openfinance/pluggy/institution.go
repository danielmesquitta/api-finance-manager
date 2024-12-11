package pluggy

import (
	"context"
	"encoding/json"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/openfinance"
)

type connectorsResponse struct {
	Results []result `json:"results"`
}

type result struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	ImageURL string `json:"image_url"`
}

func (c *Client) ListInstitutions(
	ctx context.Context,
) ([]openfinance.Institution, error) {
	if err := c.refreshAccessToken(ctx); err != nil {
		return nil, errs.New(err)
	}

	res, err := c.c.R().
		SetContext(ctx).
		Get("/connectors")
	if err != nil {
		return nil, errs.New(err)
	}
	body := res.Body()
	if statusCode := res.StatusCode(); statusCode < 200 || statusCode >= 300 {
		return nil, errs.New(body)
	}

	connectors := connectorsResponse{}
	if err := json.Unmarshal(body, &connectors); err != nil {
		return nil, errs.New(err)
	}

	var institutions []openfinance.Institution
	for _, connector := range connectors.Results {
		institutions = append(institutions, openfinance.Institution{
			ID:       connector.ID,
			Name:     connector.Name,
			ImageURL: connector.ImageURL,
		})
	}

	return institutions, nil
}
