package usecase

import (
	"context"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/money"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/validator"
)

type CalculateCompoundInterest struct {
	v *validator.Validator
}

func NewCalculateCompoundInterest(
	v *validator.Validator,
) *CalculateCompoundInterest {
	return &CalculateCompoundInterest{
		v: v,
	}
}

type CalculateCompoundInterestInput struct {
	InitialDeposit int64               `json:"initial_deposit"`
	MonthlyDeposit int64               `json:"monthly_deposit"`
	Interest       int64               `json:"interest"         validate:"required,min=0,max=10000"`
	InterestType   entity.InterestType `json:"interest_type"    validate:"required,oneof=MONTHLY ANNUAL"`
	PeriodInMonths int                 `json:"period_in_months" validate:"required,min=1"`
}

type CalculateCompoundInterestOutput struct {
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

func (uc *CalculateCompoundInterest) Execute(
	ctx context.Context,
	in CalculateCompoundInterestInput,
) (*CalculateCompoundInterestOutput, error) {
	type Result struct {
		out *CalculateCompoundInterestOutput
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

func (uc *CalculateCompoundInterest) execute(
	in CalculateCompoundInterestInput,
) (*CalculateCompoundInterestOutput, error) {
	if err := uc.validate(in); err != nil {
		return nil, errs.New(err)
	}

	output := &CalculateCompoundInterestOutput{
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

		output.ByMonth[month] = CompoundInterestResult{
			TotalAmount:     money.ToCents(currentBalance),
			TotalInterest:   money.ToCents(totalInterest),
			TotalDeposit:    money.ToCents(totalDeposit),
			MonthlyInterest: money.ToCents(monthlyInterest),
		}
	}

	output.TotalAmount = money.ToCents(currentBalance)
	output.TotalInterest = money.ToCents(totalInterest)
	output.TotalDeposit = money.ToCents(totalDeposit)

	return output, nil
}

func (uc *CalculateCompoundInterest) validate(
	in CalculateCompoundInterestInput,
) error {
	if err := uc.v.Validate(in); err != nil {
		return errs.New(err)
	}

	if in.InitialDeposit == 0 && in.MonthlyDeposit == 0 {
		return errs.ErrInvalidCompoundInterestInput
	}

	return nil
}
