package entity

import (
	"time"

	"github.com/google/uuid"
)

type Budget struct {
	ID               uuid.UUID        `json:"id,omitempty"`
	Amount           int64            `json:"amount,omitempty"`
	CreatedAt        time.Time        `json:"created_at,omitempty"`
	UpdatedAt        time.Time        `json:"updated_at,omitempty"`
	UserID           uuid.UUID        `json:"user_id,omitempty"`
	User             *User            `json:"user,omitempty"`
	BudgetCategories []BudgetCategory `json:"budget_categories,omitempty"`
}
