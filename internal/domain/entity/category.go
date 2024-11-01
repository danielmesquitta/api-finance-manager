package entity

import (
	"time"

	"github.com/google/uuid"
)

type Category struct {
	ID         uuid.UUID `json:"id,omitempty"`
	ExternalID string    `json:"external_id,omitempty"`
	Name       string    `json:"name,omitempty"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
	UpdatedAt  time.Time `json:"updated_at,omitempty"`

	Transactions     []Transaction    `json:"transactions,omitempty"`
	BudgetCategories []BudgetCategory `json:"budget_categories,omitempty"`
}
