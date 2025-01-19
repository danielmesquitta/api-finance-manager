package errs

var (
	ErrInvalidRetirementAge = NewWithType(
		"Sua idade atual deve ser menor que a idade da aposentadoria",
		ErrTypeValidation,
	)
	ErrInvalidLifeExpectance = NewWithType(
		"Sua expectativa de vida deve ser maior que 0",
		ErrTypeValidation,
	)
	ErrInvalidCompoundInterestInput = NewWithType(
		"Pelo menos um dos valores de depósito inicial e mensal deve ser diferente de 0",
		ErrTypeValidation,
	)
)
