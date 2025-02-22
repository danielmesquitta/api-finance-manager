package errs

var (
	ErrInvalidDate = New(
		`Data inválida, utilize uma data no formato "2006-01-02T15:04:05+07:00"`,
		ErrCodeValidation,
	)
	ErrInvalidUUID = New(
		`Utilize um UUID no formato "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"`,
		ErrCodeValidation,
	)
	ErrInvalidBool = New(
		`Utilize um valor verdadeiro ou falso (1, t, T, TRUE, true, True, 0, f, F, FALSE, false, False)`,
		ErrCodeValidation,
	)
)
