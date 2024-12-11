package usecase

import (
	"context"
	"math"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/money"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/validator"
)

type CalculateEmergencyReserveUseCase struct {
	v validator.Validator
}

func NewCalculateEmergencyReserveUseCase(
	v validator.Validator,
) *CalculateEmergencyReserveUseCase {
	return &CalculateEmergencyReserveUseCase{
		v: v,
	}
}

type CalculateEmergencyReserveUseCaseInput struct {
	JobType                  entity.JobType `validate:"required,oneof=ENTREPRENEUR EMPLOYEE CIVIL_SERVANT"`
	MonthlyExpenses          float64        `validate:"min=0"`
	MonthlyIncome            float64        `validate:"min=0"`
	MonthlySavingsPercentage float64        `validate:"min=0,max=100"`
}

type CalculateEmergencyReserveUseCaseOutput struct {
	RecommendedReserveInMonths      int
	RecommendedReserveInValue       float64
	MonthsToAchieveEmergencyReserve int
}

func (uc *CalculateEmergencyReserveUseCase) Execute(
	ctx context.Context,
	in CalculateEmergencyReserveUseCaseInput,
) (*CalculateEmergencyReserveUseCaseOutput, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	if err := uc.v.Validate(in); err != nil {
		return nil, errs.New(err)
	}

	out := &CalculateEmergencyReserveUseCaseOutput{}

	switch in.JobType {
	case entity.JobTypeCivilServant:
		out.RecommendedReserveInMonths = 6
	case entity.JobTypeEmployee:
		out.RecommendedReserveInMonths = 6
	case entity.JobTypeEntrepreneur:
		out.RecommendedReserveInMonths = 12
	}

	monthlySavings := in.MonthlyIncome * in.MonthlySavingsPercentage / 100

	out.RecommendedReserveInValue = money.Round(
		in.MonthlyExpenses * float64(out.RecommendedReserveInMonths),
	)

	out.MonthsToAchieveEmergencyReserve = int(
		math.Ceil(out.RecommendedReserveInValue / monthlySavings),
	)

	return out, nil
}
