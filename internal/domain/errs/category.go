package errs

var (
	ErrCategoryNotFound = NewWithType(
		"Categoria não encontrada",
		ErrTypeNotFound,
	)
	ErrCategoriesNotFound = NewWithType(
		"Uma ou mais categorias não foram encontradas",
		ErrTypeNotFound,
	)
)
