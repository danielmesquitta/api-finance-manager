package entity

type FullTransaction struct {
	Transaction
	CategoryName      string  `db:"category_name"       json:"category_name,omitzero"`
	PaymentMethodName string  `db:"payment_method_name" json:"payment_method_name,omitzero"`
	InstitutionName   *string `db:"institution_name"    json:"institution_name,omitzero"`
	InstitutionLogo   *string `db:"institution_logo"    json:"institution_logo,omitzero"`
}
