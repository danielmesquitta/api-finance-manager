package dto

import "github.com/danielmesquitta/api-finance-manager/internal/domain/usecase/calc"

type CompoundInterestRequest struct {
	calc.CalculateCompoundInterestUseCaseInput
}

type CompoundInterestResponse struct {
	calc.CalculateCompoundInterestUseCaseOutput
}

type EmergencyReserveRequest struct {
	calc.CalculateEmergencyReserveUseCaseInput
}

type EmergencyReserveResponse struct {
	calc.CalculateEmergencyReserveUseCaseOutput
}

type RetirementRequest struct {
	calc.CalculateRetirementUseCaseInput
}

type RetirementResponse struct {
	calc.CalculateRetirementUseCaseOutput
}

type SimpleInterestRequest struct {
	calc.CalculateSimpleInterestUseCaseInput
}

type SimpleInterestResponse struct {
	calc.CalculateSimpleInterestUseCaseOutput
}

type CashVsInstallmentsRequest struct {
	calc.CalculateCashVsInstallmentsUseCaseInput
}

type CashVsInstallmentsResponse struct {
	calc.CalculateCashVsInstallmentsUseCaseOutput
}
