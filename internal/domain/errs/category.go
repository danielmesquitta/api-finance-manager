package errs

var (
	ErrCategoryNotFound = New(
		"Categoria não encontrada",
		ErrCodeNotFound,
	)
	ErrCategoriesNotFound = New(
		"Uma ou mais categorias não foram encontradas",
		ErrCodeNotFound,
	)
)
