package entity

import (
	"time"

	"github.com/google/uuid"
)

type BudgetCategory struct {
	ID        uuid.UUID `json:"id"`
	Amount    int64     `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	BudgetID   string    `json:"budget_id"`
	Budget     *Budget   `json:"budget"`
	CategoryID string    `json:"category_id"`
	Category   *Category `json:"category"`
}
