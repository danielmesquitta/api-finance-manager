package entity

import "github.com/google/uuid"

type TransactionWithCategoryAndInstitution struct {
	Transaction
	CategoryID      uuid.UUID `json:"category_id,omitempty"`
	CategoryName    string    `json:"category_name,omitempty"`
	InstitutionID   uuid.UUID `json:"institution_id,omitempty"`
	InstitutionName string    `json:"institution_name,omitempty"`
	InstitutionLogo string    `json:"institution_logo,omitempty"`
}
