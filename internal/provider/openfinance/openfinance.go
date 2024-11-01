package openfinance

import (
	"context"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
)

type Account struct {
	ID         string             `json:"id,omitempty"`
	Type       entity.AccountType `json:"type,omitempty"`
	Name       string             `json:"name,omitempty"`
	Balance    float64            `json:"balance,omitempty"`
	ItemID     string             `json:"itemId,omitempty"`
	CreditData *CreditData        `json:"creditData,omitempty"`
}

type CreditData struct {
	Level                string  `json:"level,omitempty"` // BLACK, GOLD, BRONZE ...
	Brand                string  `json:"brand,omitempty"` // VISA, MASTERCARD ...
	AvailableCreditLimit float64 `json:"availableCreditLimit,omitempty"`
	CreditLimit          float64 `json:"creditLimit,omitempty"`
}

type Institution struct {
	Connector Connector `json:"connector"`
}

type Connector struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	ImageURL string `json:"imageUrl"`
}

type Client interface {
	GetAccount(ctx context.Context, accountID string) (*Account, error)
	GetInstitution(
		ctx context.Context,
		accountItemID string,
	) (*Institution, error)
}
