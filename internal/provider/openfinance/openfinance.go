package openfinance

import (
	"context"
)

type Institution struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	ImageURL string `json:"image_url"`
}

type Client interface {
	ListInstitutions(
		ctx context.Context,
	) ([]Institution, error)
}
