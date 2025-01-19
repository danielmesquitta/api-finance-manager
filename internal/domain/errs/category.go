package errs

var (
	ErrCategoriesNotFound = NewWithType(
		"Uma ou mais categorias não foram encontradas",
		ErrTypeNotFound,
	)
)
