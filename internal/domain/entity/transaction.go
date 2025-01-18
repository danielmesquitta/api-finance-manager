package entity

type TransactionWithCategoryAndInstitution struct {
	Transaction
	CategoryName      string `json:"category_name,omitempty"`
	PaymentMethodName string `json:"payment_method_name,omitempty"`
	InstitutionName   string `json:"institution_name,omitempty"`
	InstitutionLogo   string `json:"institution_logo,omitempty"`
}
