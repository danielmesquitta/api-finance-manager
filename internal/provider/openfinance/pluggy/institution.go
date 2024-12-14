package pluggy

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
)

type connectorsResponse struct {
	Results []connectorsResult `json:"results"`
}

type connectorsResult struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	ImageURL string `json:"imageUrl"`
}

func (c *Client) ListInstitutions(
	ctx context.Context,
) ([]entity.Institution, error) {
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

	var institutions []entity.Institution
	for _, r := range connectors.Results {
		var logo *string
		if r.ImageURL != "" {
			logo = &r.ImageURL
		}
		institutions = append(institutions, entity.Institution{
			ExternalID: strconv.Itoa(r.ID),
			Name:       r.Name,
			Logo:       logo,
		})
	}

	return institutions, nil
}
