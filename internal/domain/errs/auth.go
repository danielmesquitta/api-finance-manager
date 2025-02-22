package errs

var (
	ErrSubscriptionExpired = New(
		"Assinatura expirada",
		ErrCodeUnauthorized,
	)
	ErrUnauthorized = New(
		"Usuário não autorizado",
		ErrCodeUnauthorized,
	)
)
