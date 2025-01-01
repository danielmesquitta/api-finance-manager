package usecase

import (
	"context"

	"github.com/shopspring/decimal"

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
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	if err := uc.v.Validate(in); err != nil {
		return nil, errs.New(err)
	}

	output := &CalculateSimpleInterestOutput{
		ByMonth: make(map[int]SimpleInterestResult, in.PeriodInMonths),
	}

	interest := decimal.New(in.Interest, -4)
	var monthlyInterestRate decimal.Decimal

	switch in.InterestType {
	case entity.InterestTypeMonthly:
		monthlyInterestRate = interest

	case entity.InterestTypeAnnual:
		monthlyInterestRate = money.ToMonthlyInterestRate(interest)
	}

	currentBalance := money.FromCents(in.InitialDeposit)
	totalDeposit := currentBalance
	var totalInterest decimal.Decimal

	monthlyInterest := money.
		FromCents(in.InitialDeposit).
		Mul(monthlyInterestRate)

	for month := 1; month <= in.PeriodInMonths; month++ {
		currentBalance = currentBalance.Add(monthlyInterest)
		totalInterest = totalInterest.Add(monthlyInterest)

		output.ByMonth[month] = SimpleInterestResult{
			TotalAmount:   money.ToCents(currentBalance),
			TotalInterest: money.ToCents(totalInterest),
			TotalDeposit:  money.ToCents(totalDeposit),
		}
	}

	output.TotalAmount = money.ToCents(currentBalance)
	output.TotalInterest = money.ToCents(totalInterest)
	output.TotalDeposit = money.ToCents(totalDeposit)

	return output, nil
}
