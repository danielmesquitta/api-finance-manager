package usecase

import (
	"context"
	"log/slog"

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
	InitialDeposit float64             `json:"initial_deposit,omitempty"`
	MonthlyDeposit float64             `json:"monthly_deposit,omitempty"`
	Interest       float64             `json:"interest,omitempty"         validate:"required,min=0,max=100"`
	InterestType   entity.InterestType `json:"interest_type,omitempty"    validate:"required,oneof=MONTHLY ANNUAL"`
	PeriodInMonths int                 `json:"period_in_months,omitempty" validate:"required,min=1"`
}

type CalculateCompoundInterestOutput struct {
	TotalAmount   float64                        `json:"total_amount,omitempty"`
	TotalInterest float64                        `json:"total_interest,omitempty"`
	TotalDeposit  float64                        `json:"total_deposit,omitempty"`
	ByMonth       map[int]CompoundInterestResult `json:"by_month,omitempty"`
}

type CompoundInterestResult struct {
	TotalAmount     float64 `json:"total_amount,omitempty"`
	TotalInterest   float64 `json:"total_interest,omitempty"`
	TotalDeposit    float64 `json:"total_deposit,omitempty"`
	MonthlyInterest float64 `json:"monthly_interest,omitempty"`
}

func (uc *CalculateCompoundInterest) Execute(
	ctx context.Context,
	in CalculateCompoundInterestInput,
) (*CalculateCompoundInterestOutput, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	if err := uc.validate(in); err != nil {
		return nil, errs.New(err)
	}

	output := &CalculateCompoundInterestOutput{
		ByMonth: make(map[int]CompoundInterestResult, in.PeriodInMonths),
	}

	interestRate := in.Interest / 100
	monthlyInterestRate := 0.0
	switch in.InterestType {
	case entity.InterestTypeMonthly:
		monthlyInterestRate = interestRate
	case entity.InterestTypeAnnual:
		monthlyInterestRate = money.CompoundInterestAnnualToMonthlyRate(
			interestRate,
		)
	}

	currentBalance := in.InitialDeposit
	totalDeposit := in.InitialDeposit
	totalInterest := 0.0

	for month := 1; month <= in.PeriodInMonths; month++ {
		monthlyInterest := currentBalance * monthlyInterestRate
		currentBalance += monthlyInterest + in.MonthlyDeposit
		totalDeposit += in.MonthlyDeposit
		totalInterest += monthlyInterest

		output.ByMonth[month] = CompoundInterestResult{
			TotalAmount:     money.Round(currentBalance),
			TotalInterest:   money.Round(totalInterest),
			TotalDeposit:    money.Round(totalDeposit),
			MonthlyInterest: money.Round(monthlyInterest),
		}
	}

	output.TotalAmount = money.Round(currentBalance)
	output.TotalInterest = money.Round(totalInterest)
	output.TotalDeposit = money.Round(totalDeposit)

	slog.Info("CalculateCompoundInterest.Execute", "output", output)

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
