package pluggy

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/openfinance"
)

type ConnectorsResponse struct {
	Results []ConnectorsResult `json:"results"`
}

type ConnectorsResult struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	ImageURL string `json:"imageUrl"`
	Type     string `json:"type"`
}

func (c *Client) ListInstitutions(
	ctx context.Context,
	params openfinance.ListInstitutionsParams,
) ([]entity.Institution, error) {
	if err := c.refreshAccessToken(ctx); err != nil {
		return nil, errs.New(err)
	}

	req := c.c.R().SetContext(ctx)
	if len(params.Types) > 0 {
		req.SetQueryParam("types", strings.Join(params.Types, ","))
	}
	if params.Search != "" {
		req.SetQueryParam("name", params.Search)
	}

	res, err := req.
		Get("/connectors")
	if err != nil {
		return nil, errs.New(err)
	}
	body := res.Body()
	if statusCode := res.StatusCode(); statusCode < 200 || statusCode >= 300 {
		return nil, errs.New(body)
	}

	connectors := ConnectorsResponse{}
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
