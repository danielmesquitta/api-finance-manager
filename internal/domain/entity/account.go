package entity

import "time"

type FullAccount struct {
	Account
	OpenFinanceID  *string    `json:"open_finance_id,omitempty"`
	SynchronizedAt *time.Time `json:"synchronized_at,omitempty"`
}
