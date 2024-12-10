package dto

import "github.com/danielmesquitta/api-finance-manager/internal/domain/usecase"

type CompoundInterestRequestDTO struct {
	usecase.CalculateCompoundInterestUseCaseInput
}

type CompoundInterestResponseDTO struct {
	usecase.CalculateCompoundInterestUseCaseOutput
}

type RetirementRequestDTO struct {
	usecase.CalculateRetirementUseCaseInput
}

type RetirementResponseDTO struct {
	usecase.CalculateRetirementUseCaseOutput
}