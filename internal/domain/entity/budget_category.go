package entity

import (
	"time"

	"github.com/google/uuid"
)

type BudgetCategory struct {
	ID         uuid.UUID `json:"id,omitempty"`
	Amount     int64     `json:"amount,omitempty"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
	UpdatedAt  time.Time `json:"updated_at,omitempty"`
	BudgetID   uuid.UUID `json:"budget_id,omitempty"`
	Budget     *Budget   `json:"budget,omitempty"`
	CategoryID uuid.UUID `json:"category_id,omitempty"`
	Category   *Category `json:"category,omitempty"`
}
