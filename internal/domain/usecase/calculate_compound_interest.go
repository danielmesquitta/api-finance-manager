package usecase

import (
	"context"
	"math"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/moneyutil"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/validator"
)

type CalculateCompoundInterestUseCase struct {
	v validator.Validator
}

func NewCalculateCompoundInterestUseCase(
	v validator.Validator,
) *CalculateCompoundInterestUseCase {
	return &CalculateCompoundInterestUseCase{
		v: v,
	}
}

type CalculateCompoundInterestUseCaseInput struct {
	InitialDeposit float64             `validate:"min=0"                         json:"initial_deposit,omitempty"`
	MonthlyDeposit float64             `                                         json:"monthly_deposit,omitempty"`
	Interest       float64             `validate:"required,min=0,max=100"        json:"interest,omitempty"`
	InterestType   entity.InterestType `validate:"required,oneof=MONTHLY ANNUAL" json:"interest_type,omitempty"`
	PeriodInMonths int                 `validate:"required,min=1"                json:"period_in_months,omitempty"`
}

type CalculateCompoundInterestUseCaseOutput struct {
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

func (uc *CalculateCompoundInterestUseCase) Execute(
	ctx context.Context,
	in CalculateCompoundInterestUseCaseInput,
) (*CalculateCompoundInterestUseCaseOutput, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	if err := uc.v.Validate(in); err != nil {
		return nil, errs.New(err)
	}

	output := &CalculateCompoundInterestUseCaseOutput{
		ByMonth: make(map[int]CompoundInterestResult, in.PeriodInMonths),
	}

	monthlyInterestRate := 0.0
	switch in.InterestType {
	case entity.InterestTypeMonthly:
		monthlyInterestRate = in.Interest / 100
	case entity.InterestTypeAnnual:
		monthlyInterestRate = math.Pow(1+in.Interest/100, 1.0/12) - 1
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
			TotalAmount:     moneyutil.Round2Decimal(currentBalance),
			TotalInterest:   moneyutil.Round2Decimal(totalInterest),
			TotalDeposit:    moneyutil.Round2Decimal(totalDeposit),
			MonthlyInterest: moneyutil.Round2Decimal(monthlyInterest),
		}
	}

	output.TotalAmount = moneyutil.Round2Decimal(currentBalance)
	output.TotalInterest = moneyutil.Round2Decimal(totalInterest)
	output.TotalDeposit = moneyutil.Round2Decimal(totalDeposit)

	return output, nil
}
