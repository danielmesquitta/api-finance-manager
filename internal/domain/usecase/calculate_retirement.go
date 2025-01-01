package usecase

import (
	"context"

	"github.com/shopspring/decimal"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/money"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/validator"
)

type CalculateRetirement struct {
	v   *validator.Validator
	cci *CalculateCompoundInterest
}

func NewCalculateRetirement(
	v *validator.Validator,
	cci *CalculateCompoundInterest,
) *CalculateRetirement {
	return &CalculateRetirement{
		v:   v,
		cci: cci,
	}
}

type CalculateRetirementInput struct {
	MonthlyIncome              int64               `json:"monthly_income"               validate:"required,min=0"`
	IncomeInvestmentPercentage int64               `json:"income_investment_percentage" validate:"required,min=0,max=10000"`
	InitialDeposit             int64               `json:"initial_deposit"`
	Interest                   int64               `json:"interest"                     validate:"required,min=0,max=10000"`
	InterestType               entity.InterestType `json:"interest_type"                validate:"required,oneof=MONTHLY ANNUAL"`
	GoalPatrimony              int64               `json:"goal_patrimony"               validate:"required,min=0"`
	GoalIncome                 int64               `json:"goal_income"                  validate:"required,min=0"`
	Age                        int                 `json:"age"                          validate:"required,min=0"`
	RetirementAge              int                 `json:"retirement_age"               validate:"required,min=1"`
	LifeExpectancy             int                 `json:"life_expectancy"              validate:"required,min=1"`
}

type CalculateRetirementOutput struct {
	PropertyOnRetirement  int64 `json:"property_on_retirement"`
	Heritage              int64 `json:"heritage"`
	AchievedGoalPatrimony bool  `json:"achieved_goal_patrimony"`
	MaxMonthlyExpenses    int64 `json:"max_monthly_expenses"`
	AchievedGoalIncome    bool  `json:"achieved_goal_income"`
	ExceededGoalAmount    int64 `json:"exceeded_goal_amount"`
	ExceededGoal          bool  `json:"exceeded_goal"`
}

func (uc *CalculateRetirement) Execute(
	ctx context.Context,
	in CalculateRetirementInput,
) (*CalculateRetirementOutput, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	if err := uc.validate(in); err != nil {
		return nil, errs.New(err)
	}

	incomeInvestmentPercentage := decimal.New(in.IncomeInvestmentPercentage, -4)

	monthlyDeposit := money.
		FromCents(in.MonthlyIncome).
		Mul(incomeInvestmentPercentage)

	resultsOnRetirementDate, err := uc.cci.Execute(
		ctx,
		CalculateCompoundInterestInput{
			InitialDeposit: in.InitialDeposit,
			MonthlyDeposit: money.ToCents(monthlyDeposit),
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
		CalculateCompoundInterestInput{
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

	heritage := resultsOnExpectedDeathDate.TotalAmount

	exceededGoalAmount := heritage - in.GoalPatrimony

	out := &CalculateRetirementOutput{
		PropertyOnRetirement:  resultsOnRetirementDate.TotalAmount,
		Heritage:              heritage,
		MaxMonthlyExpenses:    maxMonthlyExpenses,
		AchievedGoalPatrimony: resultsOnRetirementDate.TotalAmount >= in.GoalPatrimony,
		AchievedGoalIncome:    maxMonthlyExpenses >= in.GoalIncome,
		ExceededGoalAmount:    exceededGoalAmount,
		ExceededGoal:          exceededGoalAmount > 0,
	}

	return out, nil
}

func (uc *CalculateRetirement) validate(
	in CalculateRetirementInput,
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
