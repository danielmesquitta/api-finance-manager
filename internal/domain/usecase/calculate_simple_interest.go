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
	InitialDeposit int64               `validate:"min=0"                         json:"initial_deposit,omitempty"`
	Interest       float64             `validate:"required,min=0,max=100"        json:"interest,omitempty"`
	InterestType   entity.InterestType `validate:"required,oneof=MONTHLY ANNUAL" json:"interest_type,omitempty"`
	PeriodInMonths int                 `validate:"required,min=1"                json:"period_in_months,omitempty"`
}

type CalculateSimpleInterestOutput struct {
	SimpleInterestResult
	ByMonth map[int]SimpleInterestResult `json:"by_month,omitempty"`
}

type SimpleInterestResult struct {
	TotalAmount   int64 `json:"total_amount,omitempty"`
	TotalInterest int64 `json:"total_interest,omitempty"`
	TotalDeposit  int64 `json:"total_deposit,omitempty"`
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

	monthlyInterestRate := 0.0
	switch in.InterestType {
	case entity.InterestTypeMonthly:
		monthlyInterestRate = in.Interest / 100
	case entity.InterestTypeAnnual:
		monthlyInterestRate = in.Interest / 100 / 12
	}

	currentBalance := money.FromCents(in.InitialDeposit)
	totalDeposit := currentBalance
	totalInterest := 0.0

	monthlyInterest := money.FromCents(in.InitialDeposit) * monthlyInterestRate

	for month := 1; month <= in.PeriodInMonths; month++ {
		currentBalance += monthlyInterest
		totalInterest += monthlyInterest

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
