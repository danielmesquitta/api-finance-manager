package dto

import "github.com/danielmesquitta/api-finance-manager/internal/domain/usecase"

type CompoundInterestRequest struct {
	usecase.CalculateCompoundInterestUseCaseInput
}

type CompoundInterestResponse struct {
	usecase.CalculateCompoundInterestUseCaseOutput
}

type EmergencyReserveRequest struct {
	usecase.CalculateEmergencyReserveUseCaseInput
}

type EmergencyReserveResponse struct {
	usecase.CalculateEmergencyReserveUseCaseOutput
}

type RetirementRequest struct {
	usecase.CalculateRetirementUseCaseInput
}

type RetirementResponse struct {
	usecase.CalculateRetirementUseCaseOutput
}

type SimpleInterestRequest struct {
	usecase.CalculateSimpleInterestUseCaseInput
}

type SimpleInterestResponse struct {
	usecase.CalculateSimpleInterestUseCaseOutput
}
