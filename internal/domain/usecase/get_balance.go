package usecase

import (
	"context"
	"time"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/money"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/validator"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

type GetBalance struct {
	v   *validator.Validator
	tr  repo.TransactionRepo
	abr repo.AccountBalanceRepo
}

func NewGetBalance(
	v *validator.Validator,
	tr repo.TransactionRepo,
	abr repo.AccountBalanceRepo,
) *GetBalance {
	return &GetBalance{
		v:   v,
		tr:  tr,
		abr: abr,
	}
}

type GetBalanceInput struct {
	repo.TransactionOptions
	Date   time.Time `json:"date,omitempty"    validate:"required"`
	UserID uuid.UUID `json:"user_id,omitempty" validate:"required"`
}

type GetBalanceOutput struct {
	ComparisonDates                   ComparisonDates `json:"comparison_dates"`
	CurrentBalance                    int64           `json:"current_balance,omitempty"`
	CurrentBalancePercentageVariation int64           `json:"current_balance_percentage_variation,omitempty"`
	MonthlyBalance                    int64           `json:"monthly_balance,omitempty"`
	MonthlyBalancePercentageVariation int64           `json:"monthly_balance_percentage_variation,omitempty"`
}

func (uc *GetBalance) Execute(
	ctx context.Context,
	in GetBalanceInput,
) (*GetBalanceOutput, error) {
	if err := uc.v.Validate(in); err != nil {
		return nil, errs.New(err)
	}

	comparisonDates := calculateComparisonDates(in.Date)

	g, gCtx := errgroup.WithContext(ctx)
	var currentBalance, previousMonthBalance, monthlyBalance, previousMonthlyBalance int64

	g.Go(func() error {
		var err error
		currentBalance, err = uc.abr.GetUserBalanceOnDate(
			gCtx,
			repo.GetUserBalanceOnDateParams{
				UserID: in.UserID,
				Date:   comparisonDates.MonthComparisonEndDate,
			},
		)
		return err
	})

	g.Go(func() error {
		var err error
		previousMonthBalance, err = uc.abr.GetUserBalanceOnDate(
			gCtx,
			repo.GetUserBalanceOnDateParams{
				UserID: in.UserID,
				Date:   comparisonDates.PreviousMonthComparisonEndDate,
			},
		)
		return err
	})

	g.Go(func() error {
		var err error
		in.TransactionOptions.StartDate = comparisonDates.MonthStart
		in.TransactionOptions.EndDate = comparisonDates.MonthComparisonEndDate
		opts := prepareTransactionOptions(in.TransactionOptions, time.Time{})
		monthlyBalance, err = uc.tr.SumTransactions(
			gCtx,
			in.UserID,
			opts...,
		)
		return err
	})

	g.Go(func() error {
		var err error
		in.TransactionOptions.StartDate = comparisonDates.PreviousMonthStart
		in.TransactionOptions.EndDate = comparisonDates.PreviousMonthComparisonEndDate
		opts := prepareTransactionOptions(in.TransactionOptions, time.Time{})
		previousMonthlyBalance, err = uc.tr.SumTransactions(
			gCtx,
			in.UserID,
			opts...,
		)
		return err
	})

	if err := g.Wait(); err != nil {
		return nil, errs.New(err)
	}

	currentBalancePercentageVariation := uc.calculatePercentageVariation(
		currentBalance,
		previousMonthBalance,
	)

	monthlyBalancePercentageVariation := uc.calculatePercentageVariation(
		monthlyBalance,
		previousMonthlyBalance,
	)

	out := &GetBalanceOutput{
		ComparisonDates:                   *comparisonDates,
		CurrentBalance:                    currentBalance,
		CurrentBalancePercentageVariation: currentBalancePercentageVariation,
		MonthlyBalance:                    monthlyBalance,
		MonthlyBalancePercentageVariation: monthlyBalancePercentageVariation,
	}

	return out, nil
}

func (uc *GetBalance) calculatePercentageVariation(
	a, b int64,
) int64 {
	if b == 0 {
		return 0
	}
	return money.FromPercentage(1 - (float64(a) / float64(b)))
}
