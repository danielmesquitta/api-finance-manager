package errs

var (
	ErrSubscriptionExpired = NewWithType(
		"Assinatura expirada",
		ErrTypeUnauthorized,
	)
	ErrUnauthorized = NewWithType(
		"Usuário não autorizado",
		ErrTypeUnauthorized,
	)
)
