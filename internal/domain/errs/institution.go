package errs

var (
	ErrInstitutionNotFound = NewWithType(
		"Instituição não encontrada",
		ErrTypeNotFound,
	)
)
