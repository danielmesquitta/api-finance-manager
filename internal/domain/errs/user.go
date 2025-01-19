package errs

var (
	ErrUserNotFound = NewWithType(
		"Usuário não encontrado",
		ErrTypeNotFound,
	)
)
