package errs

var (
	ErrAIChatNotFound = New(
		"Usuário não encontrado",
		ErrCodeNotFound,
	)
)
