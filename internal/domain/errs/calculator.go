package errs

var (
	ErrInvalidRetirementAge = New(
		"Sua idade atual deve ser menor que a idade da aposentadoria",
		ErrCodeValidation,
	)
	ErrInvalidLifeExpectance = New(
		"Sua expectativa de vida deve ser maior que 0",
		ErrCodeValidation,
	)
	ErrInvalidCompoundInterestInput = New(
		"Pelo menos um dos valores de depósito inicial e mensal deve ser diferente de 0",
		ErrCodeValidation,
	)
)
