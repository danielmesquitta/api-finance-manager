package usecase

import (
	"context"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/validator"
)

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
	MonthlyIncome              float64             `validate:"required,min=0"                json:"monthly_income,omitempty"`
	IncomeInvestmentPercentage float64             `validate:"required,min=0,max=100"        json:"income_investment_percentage,omitempty"`
	InitialDeposit             float64             `validate:"required,min=0"                json:"initial_deposit,omitempty"`
	Interest                   float64             `validate:"required,min=0,max=100"        json:"interest,omitempty"`
	InterestType               entity.InterestType `validate:"required,oneof=MONTHLY ANNUAL" json:"interest_type,omitempty"`
	GoalPatrimony              float64             `validate:"required,min=0"                json:"goal_patrimony,omitempty"`
	GoalIncome                 float64             `validate:"required,min=0"                json:"goal_income,omitempty"`
	Age                        int                 `validate:"required,min=0"                json:"age,omitempty"`
	RetirementAge              int                 `validate:"required,min=1"                json:"retirement_age,omitempty"`
	LifeExpectancy             int                 `validate:"min=1"                         json:"life_expectancy,omitempty"`
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

	uc.prepare(&in)

	if err := uc.validate(in); err != nil {
		return nil, errs.New(err)
	}

	resultsOnRetirementDate, err := uc.cci.Execute(
		ctx,
		CalculateCompoundInterestUseCaseInput{
			InitialDeposit: in.InitialDeposit,
			MonthlyDeposit: in.MonthlyIncome * in.IncomeInvestmentPercentage / 100,
			Interest:       in.Interest,
			InterestType:   in.InterestType,
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
			Interest:       in.Interest,
			InterestType:   in.InterestType,
			PeriodInMonths: (in.LifeExpectancy - in.RetirementAge) * 12,
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

func (uc *CalculateRetirementUseCase) prepare(
	in *CalculateRetirementUseCaseInput,
) {
	if in.LifeExpectancy == 0 {
		in.LifeExpectancy = 72
	}
}

func (uc *CalculateRetirementUseCase) validate(
	in CalculateRetirementUseCaseInput,
) error {
	if err := uc.v.Validate(in); err != nil {
		return errs.New(err)
	}

	if in.Age >= in.RetirementAge {
		return errs.ErrInvalidRetirementAge
	}

	if in.RetirementAge >= in.LifeExpectancy {
		return errs.ErrInvalidLifeExpectance
	}

	return nil
}
