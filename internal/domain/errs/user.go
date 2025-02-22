package errs

var (
	ErrUserNotFound = New(
		"Usuário não encontrado",
		ErrCodeNotFound,
	)
	ErrPremiumUsersNotFound = New(
		"Não foi possível encontrar usuários com tier premium ou trial",
		ErrCodeNotFound,
	)
)
