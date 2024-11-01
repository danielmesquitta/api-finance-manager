package entity

import (
	"time"

	"github.com/google/uuid"
)

type Category struct {
	ID         uuid.UUID `json:"id"`
	ExternalID string    `json:"external_id"`
	Name       string    `json:"name"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`

	Transactions     []Transaction    `json:"transactions"`
	BudgetCategories []BudgetCategory `json:"budget_categories"`
}
