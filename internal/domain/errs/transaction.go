package errs

var (
	ErrTransactionNotFound = New(
		"Transação não encontrada",
		ErrCodeNotFound,
	)
	ErrInvalidDateRange = New(
		"Intervalo de datas inválido",
		ErrCodeValidation,
	)
)
