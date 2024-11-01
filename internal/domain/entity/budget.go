package entity

import (
	"time"

	"github.com/google/uuid"
)

type Budget struct {
	ID        uuid.UUID `json:"id"`
	Amount    int64     `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	UserID           string           `json:"user_id"`
	User             *User            `json:"user"`
	BudgetCategories []BudgetCategory `json:"budget_categories,omitempty"`
}
