package errs

var (
	ErrInvalidTotalBudgetCategoryAmount = NewWithType(
		"O valor total das categorias do orçamento deve ser menor ou igual ao valor do orçamento",
		ErrTypeValidation,
	)
	ErrBudgetNotFound = NewWithType(
		"Você não possui um orçamento cadastrado",
		ErrTypeNotFound,
	)
)
