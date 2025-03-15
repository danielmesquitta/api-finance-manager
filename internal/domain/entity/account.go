package entity

import (
	"time"

	"github.com/google/uuid"
)

type AccountType = string

const (
	AccountTypeBank   AccountType = "BANK"
	AccountTypeCredit AccountType = "CREDIT"
)

type FullAccount struct {
	Account
	UserID                    *uuid.UUID `json:"user_id,omitzero"`
	InstitutionID             *uuid.UUID `json:"institution_id,omitzero"`
	UserInstitutionExternalID *string    `json:"user_institution_external_id,omitzero"`
	SynchronizedAt            *time.Time `json:"synchronized_at,omitzero"`
}
