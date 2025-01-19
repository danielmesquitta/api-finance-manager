package errs

var (
	ErrTransactionNotFound = NewWithType(
		"Transação não encontrada",
		ErrTypeNotFound,
	)
)
