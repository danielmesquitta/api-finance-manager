package entity

type PaymentMethod string

const (
	PaymentMethodOther      PaymentMethod = "OTHER"
	PaymentMethodCreditCard PaymentMethod = "CREDIT_CARD"
	PaymentMethodBOLETO     PaymentMethod = "BOLETO"
	PaymentMethodDEBIT      PaymentMethod = "DEBIT"
	PaymentMethodPIX        PaymentMethod = "PIX"
	PaymentMethodTED        PaymentMethod = "TED"
	PaymentMethodTEF        PaymentMethod = "TEF"
)
