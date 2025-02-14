package errs

var (
	ErrOpenFinanceAccountsNotFound = NewWithType(
		"Não foi possível encontrar sua conta na integração com o OpenFinance",
		ErrTypeNotFound,
	)
	ErrAccountsAlreadyRegistered = NewWithType(
		"Essa conta já está registrada",
		ErrTypeForbidden,
	)
	ErrPremiumAccountsNotFound = NewWithType(
		"Não foi possível encontrar contas de usuários com tier premium ou trial",
		ErrTypeNotFound,
	)
)
