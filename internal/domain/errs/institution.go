package errs

var (
	ErrInstitutionNotFound = New(
		"Instituição não encontrada",
		ErrCodeNotFound,
	)
)
