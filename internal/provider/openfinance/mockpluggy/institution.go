package mockpluggy

import (
	"context"
	"encoding/json"
	"slices"
	"strconv"
	"strings"

	root "github.com/danielmesquitta/api-finance-manager"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/openfinance"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/openfinance/pluggy"
)

func (c *Client) ListInstitutions(
	ctx context.Context,
	options ...openfinance.InstitutionOption,
) ([]entity.Institution, error) {
	opts := openfinance.InstitutionOptions{}
	for _, opt := range options {
		opt(&opts)
	}

	data, err := root.TestData.ReadFile("testdata/pluggy/connectors.json")
	if err != nil {
		return nil, errs.New(err)
	}

	connectors := pluggy.ConnectorsResponse{}
	if err := json.Unmarshal(data, &connectors); err != nil {
		return nil, errs.New(err)
	}

	var institutions []entity.Institution
	for _, r := range connectors.Results {
		if len(opts.Types) > 0 && !slices.Contains(opts.Types, r.Type) {
			continue
		}
		if opts.Search != "" && !strings.Contains(r.Name, opts.Search) {
			continue
		}

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
