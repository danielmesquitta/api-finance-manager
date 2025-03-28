package calc

import (
	"context"
	"math"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/money"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/validator"
)

type CalculateEmergencyReserveUseCase struct {
	v *validator.Validator
}

func NewCalculateEmergencyReserveUseCase(
	v *validator.Validator,
) *CalculateEmergencyReserveUseCase {
	return &CalculateEmergencyReserveUseCase{
		v: v,
	}
}

type CalculateEmergencyReserveUseCaseInput struct {
	JobType                  entity.JobType `json:"job_type"                   validate:"required,oneof=ENTREPRENEUR EMPLOYEE CIVIL_SERVANT"`
	MonthlyExpenses          int64          `json:"monthly_expenses"           validate:"min=0"`
	MonthlyIncome            int64          `json:"monthly_income"             validate:"min=0"`
	MonthlySavingsPercentage int64          `json:"monthly_savings_percentage" validate:"min=0,max=10000"`
}

type CalculateEmergencyReserveUseCaseOutput struct {
	RecommendedReserveInMonths      int64 `json:"recommended_reserve_in_months"`
	RecommendedReserveInValue       int64 `json:"recommended_reserve_in_value"`
	MonthsToAchieveEmergencyReserve int64 `json:"months_to_achieve_emergency_reserve"`
}

func (uc *CalculateEmergencyReserveUseCase) Execute(
	ctx context.Context,
	in CalculateEmergencyReserveUseCaseInput,
) (*CalculateEmergencyReserveUseCaseOutput, error) {
	type Result struct {
		out *CalculateEmergencyReserveUseCaseOutput
		err error
	}

	ch := make(chan Result)
	defer close(ch)

	go func() {
		out, err := uc.execute(in)
		ch <- Result{out, err}
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case result := <-ch:
		return result.out, result.err
	}
}

func (uc *CalculateEmergencyReserveUseCase) execute(
	in CalculateEmergencyReserveUseCaseInput,
) (*CalculateEmergencyReserveUseCaseOutput, error) {
	if err := uc.v.Validate(in); err != nil {
		return nil, errs.New(err)
	}

	out := &CalculateEmergencyReserveUseCaseOutput{}

	switch in.JobType {
	case entity.JobTypeCivilServant:
		out.RecommendedReserveInMonths = 3
	case entity.JobTypeEmployee:
		out.RecommendedReserveInMonths = 6
	case entity.JobTypeEntrepreneur:
		out.RecommendedReserveInMonths = 12
	}

	monthlyIncome := money.FromCents(in.MonthlyIncome)
	monthlySavingsPercentage := money.ToPercentage(in.MonthlySavingsPercentage)
	monthlySavings := monthlyIncome * monthlySavingsPercentage

	out.RecommendedReserveInValue = in.MonthlyExpenses * out.RecommendedReserveInMonths
	recommendedReserveInValue := money.FromCents(out.RecommendedReserveInValue)

	out.MonthsToAchieveEmergencyReserve = int64(math.Ceil(
		recommendedReserveInValue / monthlySavings,
	))

	return out, nil
}
