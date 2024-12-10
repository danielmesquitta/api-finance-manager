package usecase

import (
	"context"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/validator"
)

const lifeExpectance = 72

type CalculateRetirementUseCase struct {
	v   validator.Validator
	cci *CalculateCompoundInterestUseCase
}

func NewCalculateRetirementUseCase(
	v validator.Validator,
	cci *CalculateCompoundInterestUseCase,
) *CalculateRetirementUseCase {
	return &CalculateRetirementUseCase{
		v:   v,
		cci: cci,
	}
}

type CalculateRetirementUseCaseInput struct {
	MonthlyIncome              float64 `validate:"required,min=0"         json:"monthly_income,omitempty"`
	IncomeInvestmentPercentage float64 `validate:"required,min=0,max=100" json:"income_investment_percentage,omitempty"`
	InitialDeposit             float64 `validate:"required,min=0"         json:"initial_deposit,omitempty"`
	AnnualInterest             float64 `validate:"required,min=0,max=100" json:"annual_interest,omitempty"`
	GoalPatrimony              float64 `validate:"required,min=0"         json:"goal_patrimony,omitempty"`
	GoalIncome                 float64 `validate:"required,min=0"         json:"goal_income,omitempty"`
	Age                        int     `validate:"required,min=0"         json:"age,omitempty"`
	RetirementAge              int     `validate:"required,min=1"         json:"retirement_age,omitempty"`
	LifeExpectancy             int     `validate:"required,min=1"         json:"life_expectancy,omitempty"`
}

type CalculateRetirementUseCaseOutput struct {
	PropertyOnRetirement  float64 `json:"property_on_retirement,omitempty"`
	Heritage              float64 `json:"heritage,omitempty"`
	MaxMonthlyExpenses    float64 `json:"max_monthly_expenses,omitempty"`
	AchievedGoalPatrimony bool    `json:"achieved_goal_patrimony,omitempty"`
	AchievedGoalIncome    bool    `json:"achieved_goal_income,omitempty"`
}

func (uc *CalculateRetirementUseCase) Execute(
	ctx context.Context,
	in CalculateRetirementUseCaseInput,
) (*CalculateRetirementUseCaseOutput, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	if err := uc.v.Validate(in); err != nil {
		return nil, errs.New(err)
	}

	if in.Age >= in.RetirementAge {
		return nil, errs.ErrInvalidRetirementAge
	}

	if in.RetirementAge >= in.LifeExpectancy {
		return nil, errs.ErrInvalidLifeExpectance
	}

	resultsOnRetirementDate, err := uc.cci.Execute(
		ctx,
		CalculateCompoundInterestUseCaseInput{
			InitialDeposit: in.InitialDeposit,
			MonthlyDeposit: in.MonthlyIncome * in.IncomeInvestmentPercentage / 100,
			AnnualInterest: in.AnnualInterest,
			PeriodInMonths: (in.RetirementAge - in.Age) * 12,
		},
	)
	if err != nil {
		return nil, errs.New(err)
	}

	resultsOnExpectedDeathDate, err := uc.cci.Execute(
		ctx,
		CalculateCompoundInterestUseCaseInput{
			InitialDeposit: resultsOnRetirementDate.TotalAmount,
			MonthlyDeposit: -1 * in.GoalIncome,
			AnnualInterest: in.AnnualInterest,
			PeriodInMonths: (lifeExpectance - in.RetirementAge) * 12,
		},
	)
	if err != nil {
		return nil, errs.New(err)
	}

	maxMonthlyExpenses := resultsOnRetirementDate.ByMonth[len(resultsOnRetirementDate.ByMonth)-1].MonthlyInterest

	out := &CalculateRetirementUseCaseOutput{
		PropertyOnRetirement:  resultsOnRetirementDate.TotalAmount,
		Heritage:              resultsOnExpectedDeathDate.TotalAmount,
		MaxMonthlyExpenses:    maxMonthlyExpenses,
		AchievedGoalPatrimony: resultsOnRetirementDate.TotalAmount >= in.GoalPatrimony,
		AchievedGoalIncome:    maxMonthlyExpenses >= in.GoalIncome,
	}

	return out, nil
}
