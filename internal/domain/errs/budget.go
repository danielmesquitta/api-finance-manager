package errs

var (
	ErrInvalidTotalBudgetCategoryAmount = New(
		"O valor total das categorias do orçamento deve ser menor ou igual ao valor do orçamento",
		ErrCodeValidation,
	)
	ErrBudgetNotFound = New(
		"Você não possui um orçamento cadastrado",
		ErrCodeNotFound,
	)
	ErrBudgetCategoryNotFound = New(
		"Você não possui um orçamento cadastrado para essa categoria",
		ErrCodeNotFound,
	)
)
