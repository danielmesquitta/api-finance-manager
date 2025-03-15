// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: account.sql

package sqlc

import (
	"github.com/google/uuid"
)

type CreateAccountsParams struct {
	ID                uuid.UUID `json:"id"`
	ExternalID        string    `json:"external_id"`
	Name              string    `json:"name"`
	Type              string    `json:"type"`
	UserInstitutionID uuid.UUID `json:"user_institution_id"`
}
