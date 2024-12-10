package dto

import "github.com/danielmesquitta/api-finance-manager/internal/domain/usecase"

type CompoundInterestRequestDTO struct {
	usecase.CalculateCompoundInterestUseCaseInput
}

type CompoundInterestResponseDTO struct {
	usecase.CalculateCompoundInterestUseCaseOutput
}

type EmergencyReserveRequestDTO struct {
	usecase.CalculateEmergencyReserveUseCaseInput
}

type EmergencyReserveResponseDTO struct {
	usecase.CalculateEmergencyReserveUseCaseOutput
}

type RetirementRequestDTO struct {
	usecase.CalculateRetirementUseCaseInput
}

type RetirementResponseDTO struct {
	usecase.CalculateRetirementUseCaseOutput
}

type SimpleInterestRequestDTO struct {
	usecase.CalculateSimpleInterestUseCaseInput
}

type SimpleInterestResponseDTO struct {
	usecase.CalculateSimpleInterestUseCaseOutput
}
