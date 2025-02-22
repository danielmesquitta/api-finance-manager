package errs

var (
	ErrOpenFinanceAccountsNotFound = New(
		"Não foi possível encontrar sua conta na integração com o OpenFinance",
		ErrCodeNotFound,
	)
	ErrAccountsAlreadyRegistered = New(
		"Essa conta já está registrada",
		ErrCodeForbidden,
	)
	ErrPremiumAccountsNotFound = New(
		"Não foi possível encontrar contas de usuários com tier premium ou trial",
		ErrCodeNotFound,
	)
)
