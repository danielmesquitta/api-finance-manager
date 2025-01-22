package errs

var (
	ErrTransactionNotFound = NewWithType(
		"Transação não encontrada",
		ErrTypeNotFound,
	)
	ErrInvalidDateRange = NewWithType(
		"Intervalo de datas inválido",
		ErrTypeValidation,
	)
)
