package calc

import (
	"context"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/money"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/validator"
)

type CalculateRetirementUseCase struct {
	v   *validator.Validator
	cci *CalculateCompoundInterestUseCase
}

func NewCalculateRetirementUseCase(
	v *validator.Validator,
	cci *CalculateCompoundInterestUseCase,
) *CalculateRetirementUseCase {
	return &CalculateRetirementUseCase{
		v:   v,
		cci: cci,
	}
}

type CalculateRetirementUseCaseInput struct {
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

type CalculateRetirementUseCaseOutput struct {
	PropertyOnRetirement  int64 `json:"property_on_retirement"`
	Heritage              int64 `json:"heritage"`
	AchievedGoalPatrimony bool  `json:"achieved_goal_patrimony"`
	MaxMonthlyExpenses    int64 `json:"max_monthly_expenses"`
	AchievedGoalIncome    bool  `json:"achieved_goal_income"`
	ExceededGoalAmount    int64 `json:"exceeded_goal_amount"`
	ExceededGoal          bool  `json:"exceeded_goal"`
}

func (uc *CalculateRetirementUseCase) Execute(
	ctx context.Context,
	in CalculateRetirementUseCaseInput,
) (*CalculateRetirementUseCaseOutput, error) {
	type Result struct {
		out *CalculateRetirementUseCaseOutput
		err error
	}

	ch := make(chan Result)
	defer close(ch)

	go func() {
		out, err := uc.execute(ctx, in)
		ch <- Result{out, err}
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case result := <-ch:
		return result.out, result.err
	}
}

func (uc *CalculateRetirementUseCase) execute(
	ctx context.Context,
	in CalculateRetirementUseCaseInput,
) (*CalculateRetirementUseCaseOutput, error) {
	if err := uc.validate(in); err != nil {
		return nil, errs.New(err)
	}

	incomeInvestmentPercentage := money.ToPercentage(
		in.IncomeInvestmentPercentage,
	)
	monthlyIncome := money.FromCents(in.MonthlyIncome)
	monthlyDeposit := monthlyIncome * incomeInvestmentPercentage

	resultsOnRetirementDate, err := uc.cci.Execute(
		ctx,
		CalculateCompoundInterestUseCaseInput{
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

	heritage := resultsOnExpectedDeathDate.TotalAmount

	exceededGoalAmount := heritage - in.GoalPatrimony

	out := &CalculateRetirementUseCaseOutput{
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
