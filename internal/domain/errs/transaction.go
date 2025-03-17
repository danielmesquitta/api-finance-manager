package errs

var (
	ErrTransactionNotFound = New(
		"Transação não encontrada",
		ErrCodeNotFound,
	)
)
