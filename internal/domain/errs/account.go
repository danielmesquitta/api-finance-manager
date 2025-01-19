package errs

var (
	ErrAccountsAlreadyRegistered = NewWithType(
		"Essa conta já está registrada",
		ErrTypeForbidden,
	)
)
