package pluggy

import (
	"context"
	"encoding/json"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
)

type categoriesResponse struct {
	Results []categoriesResult `json:"results"`
}

type categoriesResult struct {
	ID                    string `json:"id"`
	Description           string `json:"name"`
	DescriptionTranslated string `json:"descriptionTranslated"`
}

func (c *Client) ListCategories(
	ctx context.Context,
) ([]entity.Category, error) {
	if err := c.refreshAccessToken(ctx); err != nil {
		return nil, errs.New(err)
	}

	res, err := c.c.R().
		SetContext(ctx).
		Get("/categories")
	if err != nil {
		return nil, errs.New(err)
	}
	body := res.Body()
	if statusCode := res.StatusCode(); statusCode < 200 || statusCode >= 300 {
		return nil, errs.New(body)
	}

	categories := categoriesResponse{}
	if err := json.Unmarshal(body, &categories); err != nil {
		return nil, errs.New(err)
	}

	var institutions []entity.Category
	for _, connector := range categories.Results {
		institutions = append(institutions, entity.Category{
			ExternalID: connector.ID,
			Name:       connector.DescriptionTranslated,
		})
	}

	return institutions, nil
}
