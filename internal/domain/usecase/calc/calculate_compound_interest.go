package calc

import (
	"context"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/money"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/validator"
)

type CalculateCompoundInterestUseCase struct {
	v *validator.Validator
}

func NewCalculateCompoundInterestUseCase(
	v *validator.Validator,
) *CalculateCompoundInterestUseCase {
	return &CalculateCompoundInterestUseCase{
		v: v,
	}
}

type CalculateCompoundInterestUseCaseInput struct {
	InitialDeposit int64               `json:"initial_deposit"`
	MonthlyDeposit int64               `json:"monthly_deposit"`
	Interest       int64               `json:"interest"         validate:"required,min=0,max=10000"`
	InterestType   entity.InterestType `json:"interest_type"    validate:"required,oneof=MONTHLY ANNUAL"`
	PeriodInMonths int                 `json:"period_in_months" validate:"required,min=1"`
}

type CalculateCompoundInterestUseCaseOutput struct {
	TotalAmount   int64                          `json:"total_amount"`
	TotalInterest int64                          `json:"total_interest"`
	TotalDeposit  int64                          `json:"total_deposit"`
	ByMonth       map[int]CompoundInterestResult `json:"by_month"`
}

type CompoundInterestResult struct {
	TotalAmount     int64 `json:"total_amount"`
	TotalInterest   int64 `json:"total_interest"`
	TotalDeposit    int64 `json:"total_deposit"`
	MonthlyInterest int64 `json:"monthly_interest"`
}

func (uc *CalculateCompoundInterestUseCase) Execute(
	ctx context.Context,
	in CalculateCompoundInterestUseCaseInput,
) (*CalculateCompoundInterestUseCaseOutput, error) {
	type Result struct {
		out *CalculateCompoundInterestUseCaseOutput
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

func (uc *CalculateCompoundInterestUseCase) execute(
	in CalculateCompoundInterestUseCaseInput,
) (*CalculateCompoundInterestUseCaseOutput, error) {
	if err := uc.validate(in); err != nil {
		return nil, errs.New(err)
	}

	out := &CalculateCompoundInterestUseCaseOutput{
		ByMonth: make(map[int]CompoundInterestResult, in.PeriodInMonths),
	}

	interestRate := money.ToPercentage(in.Interest)
	monthlyInterestRate := 0.0
	switch in.InterestType {
	case entity.InterestTypeMonthly:
		monthlyInterestRate = interestRate
	case entity.InterestTypeAnnual:
		monthlyInterestRate = money.ToMonthlyInterestRate(interestRate)
	}

	currentBalance := money.FromCents(in.InitialDeposit)
	totalDeposit := currentBalance
	totalInterest := 0.0

	for month := 1; month <= in.PeriodInMonths; month++ {
		monthlyDeposit := money.FromCents(in.MonthlyDeposit)
		monthlyInterest := currentBalance * monthlyInterestRate
		currentBalance += monthlyInterest + monthlyDeposit
		totalDeposit += monthlyDeposit
		totalInterest += monthlyInterest

		out.ByMonth[month] = CompoundInterestResult{
			TotalAmount:     money.ToCents(currentBalance),
			TotalInterest:   money.ToCents(totalInterest),
			TotalDeposit:    money.ToCents(totalDeposit),
			MonthlyInterest: money.ToCents(monthlyInterest),
		}
	}

	out.TotalAmount = money.ToCents(currentBalance)
	out.TotalInterest = money.ToCents(totalInterest)
	out.TotalDeposit = money.ToCents(totalDeposit)

	return out, nil
}

func (uc *CalculateCompoundInterestUseCase) validate(
	in CalculateCompoundInterestUseCaseInput,
) error {
	if err := uc.v.Validate(in); err != nil {
		return errs.New(err)
	}

	if in.InitialDeposit == 0 && in.MonthlyDeposit == 0 {
		return errs.ErrInvalidCompoundInterestInput
	}

	return nil
}
