package entity

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	Name          string        `json:"name"`
	Description   string        `json:"description"`
	Amount        int64         `json:"amount"`
	PaymentMethod PaymentMethod `json:"payment_method"`
	Date          time.Time     `json:"date"`

	CategoryID uuid.UUID `json:"category_id"`
}
