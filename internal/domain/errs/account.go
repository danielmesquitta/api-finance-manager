package errs

var (
	ErrAccountsAlreadyRegistered = NewWithType(
		"Essa conta já está registrada",
		ErrTypeForbidden,
	)
	ErrPremiumAccountsNotFound = NewWithType(
		"Não foi possível encontrar contas de usuários com tier premium ou trial",
		ErrTypeNotFound,
	)
)
