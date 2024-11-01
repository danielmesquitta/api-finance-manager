package entity

type PaymentMethod string

const (
	PaymentMethodPix          PaymentMethod = "PIX"
	PaymentMethodBoleto       PaymentMethod = "BOLETO"
	PaymentMethodCreditCard   PaymentMethod = "CREDIT_CARD"
	PaymentMethodDebitCard    PaymentMethod = "DEBIT_CARD"
	PaymentMethodTransference PaymentMethod = "TRANSFERENCE"
	PaymentMethodUnknown      PaymentMethod = "UNKNOWN"
)
