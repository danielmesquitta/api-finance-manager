package usecase

import (
	"context"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/money"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/validator"
)

type CalculateCashVsInstallments struct {
	v  *validator.Validator
	ci *CalculateCompoundInterest
}

func NewCalculateCashVsInstallments(
	v *validator.Validator,
	ci *CalculateCompoundInterest,
) *CalculateCashVsInstallments {
	return &CalculateCashVsInstallments{
		v:  v,
		ci: ci,
	}
}

type CalculateCashVsInstallmentsInput struct {
	PurchaseValue      int64               `json:"purchase_value"       validate:"required,min=0"`
	CashDiscount       int64               `json:"cash_discount"        validate:"omitempty,min=0"`
	Installments       int                 `json:"installments"         validate:"required,min=1"`
	CreditCardCashback int64               `json:"cashback"             validate:"omitempty,min=0,max=10000"`
	CreditCardInterest int64               `json:"credit_card_interest" validate:"omitempty,min=0,max=10000"`
	Interest           int64               `json:"interest"             validate:"required,min=0,max=10000"`
	InterestType       entity.InterestType `json:"interest_type"        validate:"required,oneof=MONTHLY ANNUAL"`
}

type CalculateCashVsInstallmentsOutput struct {
	SavingsWithCash       int64            `json:"savings_with_cash"`
	SavingsWithCreditCard int64            `json:"savings_with_credit_card"`
	CashFlowByMonth       map[int]CashFlow `json:"cash_flow_by_month"`
}

type CashFlow struct {
	Cash       int64 `json:"cash"`
	CreditCard int64 `json:"credit_card"`
}

func (uc *CalculateCashVsInstallments) Execute(
	ctx context.Context,
	in CalculateCashVsInstallmentsInput,
) (*CalculateCashVsInstallmentsOutput, error) {
	type Result struct {
		out *CalculateCashVsInstallmentsOutput
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

func (uc *CalculateCashVsInstallments) execute(
	ctx context.Context,
	in CalculateCashVsInstallmentsInput,
) (*CalculateCashVsInstallmentsOutput, error) {
	if err := uc.v.Validate(in); err != nil {
		return nil, errs.New(err)
	}

	purchaseValue := money.FromCents(in.PurchaseValue)
	cashDiscount := money.ToPercentage(in.CashDiscount)
	initialSavingsWithCash := purchaseValue * cashDiscount

	savingsWithCash, err := uc.ci.Execute(ctx, CalculateCompoundInterestInput{
		InitialDeposit: money.ToCents(initialSavingsWithCash),
		MonthlyDeposit: 0,
		Interest:       in.Interest,
		InterestType:   in.InterestType,
		PeriodInMonths: in.Installments,
	})
	if err != nil {
		return nil, errs.New(err)
	}

	creditCardInterest := money.ToPercentage(in.CreditCardInterest)
	creditCardPurchaseValue := purchaseValue * (1 + creditCardInterest)
	installmentValue := creditCardPurchaseValue / float64(in.Installments)

	cashbackPercentage := money.ToPercentage(in.CreditCardCashback)
	creditCardCashback := creditCardPurchaseValue * cashbackPercentage

	creditCardInitialDeposit := purchaseValue + creditCardCashback

	savingsWithCreditCard, err := uc.ci.Execute(
		ctx,
		CalculateCompoundInterestInput{
			InitialDeposit: money.ToCents(creditCardInitialDeposit),
			MonthlyDeposit: money.ToCents(-1 * installmentValue),
			Interest:       in.Interest,
			InterestType:   in.InterestType,
			PeriodInMonths: in.Installments,
		},
	)
	if err != nil {
		return nil, errs.New(err)
	}

	cashFlowByMonth := map[int]CashFlow{
		0: {
			Cash:       money.ToCents(initialSavingsWithCash),
			CreditCard: money.ToCents(creditCardInitialDeposit),
		},
	}
	for month := 1; month <= in.Installments; month++ {
		cashFlowByMonth[month] = CashFlow{
			Cash:       savingsWithCash.ByMonth[month].TotalAmount,
			CreditCard: savingsWithCreditCard.ByMonth[month].TotalAmount,
		}
	}

	out := &CalculateCashVsInstallmentsOutput{
		SavingsWithCash:       savingsWithCash.TotalAmount,
		SavingsWithCreditCard: savingsWithCreditCard.TotalAmount,
		CashFlowByMonth:       cashFlowByMonth,
	}

	return out, nil
}
