package dto

import "github.com/danielmesquitta/api-finance-manager/internal/domain/usecase"

type CompoundInterestRequest struct {
	usecase.CalculateCompoundInterestInput
}

type CompoundInterestResponse struct {
	usecase.CalculateCompoundInterestOutput
}

type EmergencyReserveRequest struct {
	usecase.CalculateEmergencyReserveInput
}

type EmergencyReserveResponse struct {
	usecase.CalculateEmergencyReserveOutput
}

type RetirementRequest struct {
	usecase.CalculateRetirementInput
}

type RetirementResponse struct {
	usecase.CalculateRetirementOutput
}

type SimpleInterestRequest struct {
	usecase.CalculateSimpleInterestInput
}

type SimpleInterestResponse struct {
	usecase.CalculateSimpleInterestOutput
}

type CashVsInstallmentsRequest struct {
	usecase.CalculateCashVsInstallmentsInput
}

type CashVsInstallmentsResponse struct {
	usecase.CalculateCashVsInstallmentsOutput
}
