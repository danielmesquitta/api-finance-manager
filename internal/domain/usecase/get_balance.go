package usecase

import (
	"context"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/dateutil"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/validator"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
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
	UserID uuid.UUID `json:"user_id" validate:"required"`
}

type GetBalanceOutput struct {
	ComparisonDates            dateutil.ComparisonDates `json:"comparison_dates"`
	CurrentBalance             int64                    `json:"current_balance"`
	PreviousBalance            int64                    `json:"previous_balance"`
	BalancePercentageVariation int64                    `json:"balance_percentage_variation"`
	CurrentIncome              int64                    `json:"current_income"`
	PreviousIncome             int64                    `json:"previous_income"`
	IncomePercentageVariation  int64                    `json:"income_percentage_variation"`
	CurrentExpense             int64                    `json:"current_expense"`
	PreviousExpense            int64                    `json:"previous_expense"`
	ExpensePercentageVariation int64                    `json:"expense_percentage_variation"`
}

func (uc *GetBalance) Execute(
	ctx context.Context,
	in GetBalanceInput,
) (*GetBalanceOutput, error) {
	if err := uc.v.Validate(in); err != nil {
		return nil, errs.New(err)
	}

	cmpDates := dateutil.CalculateComparisonDates(in.StartDate, in.EndDate)

	g, gCtx := errgroup.WithContext(ctx)
	var currentBalance, previousBalance, currentIncome,
		previousIncome, currentExpense, previousExpense int64

	g.Go(func() error {
		var err error
		currentBalance, err = uc.abr.GetUserBalanceOnDate(
			gCtx,
			repo.GetUserBalanceOnDateParams{
				UserID: in.UserID,
				Date:   cmpDates.EndDate,
			},
		)
		return err
	})

	g.Go(func() error {
		var err error
		previousBalance, err = uc.abr.GetUserBalanceOnDate(
			gCtx,
			repo.GetUserBalanceOnDateParams{
				UserID: in.UserID,
				Date:   cmpDates.ComparisonEndDate,
			},
		)
		return err
	})

	g.Go(func() error {
		var err error

		var inOpts repo.TransactionOptions
		if err := copier.Copy(&inOpts, in.TransactionOptions); err != nil {
			return err
		}
		inOpts.StartDate = cmpDates.StartDate
		inOpts.EndDate = cmpDates.EndDate
		inOpts.IsIncome = true
		opts := prepareTransactionOptions(inOpts)

		currentIncome, err = uc.tr.SumTransactions(
			gCtx,
			in.UserID,
			opts...,
		)
		return err
	})

	g.Go(func() error {
		var err error

		var inOpts repo.TransactionOptions
		if err := copier.Copy(&inOpts, in.TransactionOptions); err != nil {
			return err
		}
		inOpts.StartDate = cmpDates.ComparisonStartDate
		inOpts.EndDate = cmpDates.ComparisonEndDate
		inOpts.IsIncome = true
		opts := prepareTransactionOptions(inOpts)

		previousIncome, err = uc.tr.SumTransactions(
			gCtx,
			in.UserID,
			opts...,
		)
		return err
	})

	g.Go(func() error {
		var err error

		var inOpts repo.TransactionOptions
		if err := copier.Copy(&inOpts, in.TransactionOptions); err != nil {
			return err
		}
		inOpts.StartDate = cmpDates.StartDate
		inOpts.EndDate = cmpDates.EndDate
		inOpts.IsExpense = true
		opts := prepareTransactionOptions(inOpts)

		currentExpense, err = uc.tr.SumTransactions(
			gCtx,
			in.UserID,
			opts...,
		)
		return err
	})

	g.Go(func() error {
		var err error

		var inOpts repo.TransactionOptions
		if err := copier.Copy(&inOpts, in.TransactionOptions); err != nil {
			return err
		}
		inOpts.StartDate = cmpDates.ComparisonStartDate
		inOpts.EndDate = cmpDates.ComparisonEndDate
		inOpts.IsExpense = true
		opts := prepareTransactionOptions(inOpts)

		previousExpense, err = uc.tr.SumTransactions(
			gCtx,
			in.UserID,
			opts...,
		)
		return err
	})

	if err := g.Wait(); err != nil {
		return nil, errs.New(err)
	}

	balancePercentageVariation := calculatePercentageVariation(
		currentBalance,
		previousBalance,
	)

	incomePercentageVariation := calculatePercentageVariation(
		currentIncome,
		previousIncome,
	)

	expensePercentageVariation := calculatePercentageVariation(
		currentExpense,
		previousExpense,
	)

	out := &GetBalanceOutput{
		ComparisonDates:            *cmpDates,
		CurrentBalance:             currentBalance,
		PreviousBalance:            previousBalance,
		BalancePercentageVariation: balancePercentageVariation,
		CurrentIncome:              currentIncome,
		PreviousIncome:             previousIncome,
		IncomePercentageVariation:  incomePercentageVariation,
		CurrentExpense:             currentExpense,
		PreviousExpense:            previousExpense,
		ExpensePercentageVariation: expensePercentageVariation,
	}

	return out, nil
}
