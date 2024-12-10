package usecase

import (
	"context"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/moneyutil"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/validator"
)

type CalculateSimpleInterestUseCase struct {
	v validator.Validator
}

func NewCalculateSimpleInterestUseCase(
	v validator.Validator,
) *CalculateSimpleInterestUseCase {
	return &CalculateSimpleInterestUseCase{
		v: v,
	}
}

type CalculateSimpleInterestUseCaseInput struct {
	InitialDeposit float64             `validate:"min=0"`
	Interest       float64             `validate:"required,min=0,max=100"`
	InterestType   entity.InterestType `validate:"required,oneof=MONTHLY ANNUAL"`
	PeriodInMonths int                 `validate:"required,min=1"`
}

type CalculateSimpleInterestUseCaseOutput struct {
	SimpleInterestResult
	ByMonth map[int]SimpleInterestResult `json:"by_month,omitempty"`
}

type SimpleInterestResult struct {
	TotalAmount   float64 `json:"total_amount,omitempty"`
	TotalInterest float64 `json:"total_interest,omitempty"`
	TotalDeposit  float64 `json:"total_deposit,omitempty"`
}

func (uc *CalculateSimpleInterestUseCase) Execute(
	ctx context.Context,
	in CalculateSimpleInterestUseCaseInput,
) (*CalculateSimpleInterestUseCaseOutput, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	if err := uc.v.Validate(in); err != nil {
		return nil, errs.New(err)
	}

	output := &CalculateSimpleInterestUseCaseOutput{
		ByMonth: make(map[int]SimpleInterestResult, in.PeriodInMonths),
	}

	monthlyInterestRate := 0.0
	switch in.InterestType {
	case entity.InterestTypeMonthly:
		monthlyInterestRate = in.Interest / 100
	case entity.InterestTypeAnnual:
		monthlyInterestRate = in.Interest / 100 / 12
	}

	currentBalance := in.InitialDeposit
	totalDeposit := in.InitialDeposit
	totalInterest := 0.0

	monthlyInterest := in.InitialDeposit * monthlyInterestRate

	for month := 1; month <= in.PeriodInMonths; month++ {
		currentBalance += monthlyInterest
		totalInterest += monthlyInterest

		output.ByMonth[month] = SimpleInterestResult{
			TotalAmount:   moneyutil.Round2Decimal(currentBalance),
			TotalInterest: moneyutil.Round2Decimal(totalInterest),
			TotalDeposit:  moneyutil.Round2Decimal(totalDeposit),
		}
	}

	output.TotalAmount = moneyutil.Round2Decimal(currentBalance)
	output.TotalInterest = moneyutil.Round2Decimal(totalInterest)
	output.TotalDeposit = moneyutil.Round2Decimal(totalDeposit)

	return output, nil
}
