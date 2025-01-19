package errs

var (
	ErrInvalidDate = NewWithType(
		`Data inválida, utilize uma data no formato "2006-01-02T15:04:05+07:00"`,
		ErrTypeValidation,
	)
	ErrInvalidUUID = NewWithType(
		`Utilize um UUID no formato "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"`,
		ErrTypeValidation,
	)
	ErrInvalidBool = NewWithType(
		`Utilize um valor verdadeiro ou falso (1, t, T, TRUE, true, True, 0, f, F, FALSE, false, False)`,
		ErrTypeValidation,
	)
)
