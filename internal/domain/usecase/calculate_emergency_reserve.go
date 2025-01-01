package usecase

import (
	"context"

	"github.com/shopspring/decimal"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/money"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/validator"
)

type CalculateEmergencyReserve struct {
	v *validator.Validator
}

func NewCalculateEmergencyReserve(
	v *validator.Validator,
) *CalculateEmergencyReserve {
	return &CalculateEmergencyReserve{
		v: v,
	}
}

type CalculateEmergencyReserveInput struct {
	JobType                  entity.JobType `json:"job_type"                   validate:"required,oneof=ENTREPRENEUR EMPLOYEE CIVIL_SERVANT"`
	MonthlyExpenses          int64          `json:"monthly_expenses"           validate:"min=0"`
	MonthlyIncome            int64          `json:"monthly_income"             validate:"min=0"`
	MonthlySavingsPercentage int64          `json:"monthly_savings_percentage" validate:"min=0,max=10000"`
}

type CalculateEmergencyReserveOutput struct {
	RecommendedReserveInMonths      int64 `json:"recommended_reserve_in_months"`
	RecommendedReserveInValue       int64 `json:"recommended_reserve_in_value"`
	MonthsToAchieveEmergencyReserve int64 `json:"months_to_achieve_emergency_reserve"`
}

func (uc *CalculateEmergencyReserve) Execute(
	ctx context.Context,
	in CalculateEmergencyReserveInput,
) (*CalculateEmergencyReserveOutput, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	if err := uc.v.Validate(in); err != nil {
		return nil, errs.New(err)
	}

	out := &CalculateEmergencyReserveOutput{}

	switch in.JobType {
	case entity.JobTypeCivilServant:
		out.RecommendedReserveInMonths = 3
	case entity.JobTypeEmployee:
		out.RecommendedReserveInMonths = 6
	case entity.JobTypeEntrepreneur:
		out.RecommendedReserveInMonths = 12
	}

	monthlySavingsPercentage := decimal.New(in.MonthlySavingsPercentage, -4)

	monthlySavings := money.
		FromCents(in.MonthlyIncome).
		Mul(monthlySavingsPercentage)

	recommendedReserveInMonths := decimal.New(out.RecommendedReserveInMonths, 0)

	recommendedReserveInValue := money.
		FromCents(in.MonthlyIncome).
		Mul(recommendedReserveInMonths)

	out.RecommendedReserveInValue = money.ToCents(recommendedReserveInValue)

	out.MonthsToAchieveEmergencyReserve = money.
		FromCents(out.RecommendedReserveInValue).
		Div(monthlySavings).
		Ceil().
		IntPart()

	return out, nil
}
