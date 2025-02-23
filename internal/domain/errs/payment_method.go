package errs

var (
	ErrPaymentMethodNotFound = New(
		"Método de pagamento não encontrado",
		ErrCodeNotFound,
	)
)
