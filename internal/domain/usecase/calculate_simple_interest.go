package usecase

import (
	"context"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/money"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/validator"
)

type CalculateSimpleInterest struct {
	v *validator.Validator
}

func NewCalculateSimpleInterest(
	v *validator.Validator,
) *CalculateSimpleInterest {
	return &CalculateSimpleInterest{
		v: v,
	}
}

type CalculateSimpleInterestInput struct {
	InitialDeposit int64               `json:"initial_deposit"  validate:"min=0"`
	Interest       int64               `json:"interest"         validate:"required,min=0,max=10000"`
	InterestType   entity.InterestType `json:"interest_type"    validate:"required,oneof=MONTHLY ANNUAL"`
	PeriodInMonths int                 `json:"period_in_months" validate:"required,min=1"`
}

type CalculateSimpleInterestOutput struct {
	SimpleInterestResult
	ByMonth map[int]SimpleInterestResult `json:"by_month"`
}

type SimpleInterestResult struct {
	TotalAmount   int64 `json:"total_amount"`
	TotalInterest int64 `json:"total_interest"`
	TotalDeposit  int64 `json:"total_deposit"`
}

func (uc *CalculateSimpleInterest) Execute(
	ctx context.Context,
	in CalculateSimpleInterestInput,
) (*CalculateSimpleInterestOutput, error) {
	type Result struct {
		out *CalculateSimpleInterestOutput
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

func (uc *CalculateSimpleInterest) execute(
	in CalculateSimpleInterestInput,
) (*CalculateSimpleInterestOutput, error) {
	if err := uc.v.Validate(in); err != nil {
		return nil, errs.New(err)
	}

	out := &CalculateSimpleInterestOutput{
		ByMonth: make(map[int]SimpleInterestResult, in.PeriodInMonths),
	}

	monthlyInterestRate := 0.0
	interest := money.ToPercentage(in.Interest)
	switch in.InterestType {
	case entity.InterestTypeMonthly:
		monthlyInterestRate = interest
	case entity.InterestTypeAnnual:
		monthlyInterestRate = interest / 12
	}

	currentBalance := money.FromCents(in.InitialDeposit)
	totalDeposit := currentBalance
	totalInterest := 0.0

	monthlyInterest := money.FromCents(in.InitialDeposit) * monthlyInterestRate

	for month := 1; month <= in.PeriodInMonths; month++ {
		currentBalance += monthlyInterest
		totalInterest += monthlyInterest

		out.ByMonth[month] = SimpleInterestResult{
			TotalAmount:   money.ToCents(currentBalance),
			TotalInterest: money.ToCents(totalInterest),
			TotalDeposit:  money.ToCents(totalDeposit),
		}
	}

	out.TotalAmount = money.ToCents(currentBalance)
	out.TotalInterest = money.ToCents(totalInterest)
	out.TotalDeposit = money.ToCents(totalDeposit)

	return out, nil
}
