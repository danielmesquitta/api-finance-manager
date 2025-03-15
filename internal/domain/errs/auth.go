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
	ErrInvalidProvider = New(
		"Provedor de autenticação inválido",
		ErrCodeValidation,
	)
)
