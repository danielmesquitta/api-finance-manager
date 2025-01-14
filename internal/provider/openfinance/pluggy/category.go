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
	ID                    string  `json:"id"`
	ParentID              *string `json:"parentId"`
	Description           string  `json:"name"`
	DescriptionTranslated string  `json:"descriptionTranslated"`
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
	if res.IsError() {
		return nil, errs.New(body)
	}

	categories := categoriesResponse{}
	if err := json.Unmarshal(body, &categories); err != nil {
		return nil, errs.New(err)
	}

	var institutions []entity.Category
	for _, connector := range categories.Results {
		if connector.ParentID != nil {
			continue
		}
		institutions = append(institutions, entity.Category{
			ExternalID: connector.ID,
			Name:       connector.DescriptionTranslated,
		})
	}

	return institutions, nil
}

func (c *Client) GetParentCategoryExternalID(
	externalCategoryID string,
	categoriesByExternalID map[string]entity.Category,
) (string, bool) {
	defaultCategory := "99999999"
	if externalCategoryID == "" {
		return defaultCategory, false
	}

	if _, ok := categoriesByExternalID[externalCategoryID]; ok {
		return externalCategoryID, true
	}

	parentExternalID := externalCategoryID[:2] + "000000"
	if _, ok := categoriesByExternalID[parentExternalID]; ok {
		return parentExternalID, true
	}

	return defaultCategory, false
}
