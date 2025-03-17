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
	ErrUnauthorizedTier = New(
		"Você não possui permissão para acessar esse recurso, faça o upgrade de sua assinatura",
		ErrCodeUnauthorized,
	)
)
