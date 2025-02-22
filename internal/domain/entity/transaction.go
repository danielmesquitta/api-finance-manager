package entity

type FullTransaction struct {
	Transaction
	CategoryName      string `json:"category_name,omitzero"`
	PaymentMethodName string `json:"payment_method_name,omitzero"`
	InstitutionName   string `json:"institution_name,omitzero"`
	InstitutionLogo   string `json:"institution_logo,omitzero"`
}
