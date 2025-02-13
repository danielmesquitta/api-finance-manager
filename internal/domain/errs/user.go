package errs

var (
	ErrUserNotFound = NewWithType(
		"Usuário não encontrado",
		ErrTypeNotFound,
	)
	ErrPremiumUsersNotFound = NewWithType(
		"Não foi possível encontrar usuários com tier premium ou trial",
		ErrTypeNotFound,
	)
)
