package calc

import (
	"context"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/money"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/validator"
)

type CalculateSimpleInterestUseCase struct {
	v *validator.Validator
}

func NewCalculateSimpleInterestUseCase(
	v *validator.Validator,
) *CalculateSimpleInterestUseCase {
	return &CalculateSimpleInterestUseCase{
		v: v,
	}
}

type CalculateSimpleInterestUseCaseInput struct {
	InitialDeposit int64               `json:"initial_deposit"  validate:"min=0"`
	Interest       int64               `json:"interest"         validate:"required,min=0,max=10000"`
	InterestType   entity.InterestType `json:"interest_type"    validate:"required,oneof=MONTHLY ANNUAL"`
	PeriodInMonths int                 `json:"period_in_months" validate:"required,min=1"`
}

type CalculateSimpleInterestUseCaseOutput struct {
	SimpleInterestResult
	ByMonth map[int]SimpleInterestResult `json:"by_month"`
}

type SimpleInterestResult struct {
	TotalAmount   int64 `json:"total_amount"`
	TotalInterest int64 `json:"total_interest"`
	TotalDeposit  int64 `json:"total_deposit"`
}

func (uc *CalculateSimpleInterestUseCase) Execute(
	ctx context.Context,
	in CalculateSimpleInterestUseCaseInput,
) (*CalculateSimpleInterestUseCaseOutput, error) {
	type Result struct {
		out *CalculateSimpleInterestUseCaseOutput
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

func (uc *CalculateSimpleInterestUseCase) execute(
	in CalculateSimpleInterestUseCaseInput,
) (*CalculateSimpleInterestUseCaseOutput, error) {
	if err := uc.v.Validate(in); err != nil {
		return nil, errs.New(err)
	}

	out := &CalculateSimpleInterestUseCaseOutput{
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
