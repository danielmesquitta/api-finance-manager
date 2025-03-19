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
	UserID                    *uuid.UUID `db:"user_id"                      json:"user_id,omitzero"`
	InstitutionID             *uuid.UUID `db:"institution_id"               json:"institution_id,omitzero"`
	UserInstitutionExternalID *string    `db:"user_institution_external_id" json:"user_institution_external_id,omitzero"`
	SynchronizedAt            *time.Time `db:"synchronized_at"              json:"synchronized_at,omitzero"`
}
