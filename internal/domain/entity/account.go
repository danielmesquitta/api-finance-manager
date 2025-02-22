package entity

import "time"

type FullAccount struct {
	Account
	OpenFinanceID  *string    `json:"open_finance_id,omitzero"`
	SynchronizedAt *time.Time `json:"synchronized_at,omitzero"`
}
