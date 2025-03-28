package account

import (
	"context"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/dateutil"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/money"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/validator"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"golang.org/x/sync/errgroup"
)

type GetAccountsBalanceUseCase struct {
	v   *validator.Validator
	tr  repo.TransactionRepo
	abr repo.AccountBalanceRepo
}

func NewGetAccountsBalanceUseCase(
	v *validator.Validator,
	tr repo.TransactionRepo,
	abr repo.AccountBalanceRepo,
) *GetAccountsBalanceUseCase {
	return &GetAccountsBalanceUseCase{
		v:   v,
		tr:  tr,
		abr: abr,
	}
}

type GetAccountsBalanceUseCaseInput struct {
	repo.TransactionOptions
	UserID uuid.UUID `json:"user_id" validate:"required"`
}

type GetAccountsBalanceUseCaseOutput struct {
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

func (uc *GetAccountsBalanceUseCase) Execute(
	ctx context.Context,
	in GetAccountsBalanceUseCaseInput,
) (*GetAccountsBalanceUseCaseOutput, error) {
	if err := uc.v.Validate(in); err != nil {
		return nil, errs.New(err)
	}

	balanceOpts := repo.AccountBalanceOptions{}
	if err := copier.Copy(&balanceOpts, in.TransactionOptions); err != nil {
		return nil, errs.New(err)
	}

	cmpDates := dateutil.CalculateComparisonDates(in.StartDate, in.EndDate)

	g, gCtx := errgroup.WithContext(ctx)
	var currentBalance, previousBalance, currentIncome,
		previousIncome, currentExpense, previousExpense int64

	g.Go(func() (err error) {
		currentBalance, err = uc.abr.GetUserBalanceOnDate(
			gCtx,
			in.UserID,
			cmpDates.EndDate,
			balanceOpts,
		)
		return err
	})

	g.Go(func() (err error) {
		previousBalance, err = uc.abr.GetUserBalanceOnDate(
			gCtx,
			in.UserID,
			cmpDates.ComparisonEndDate,
			balanceOpts,
		)
		return err
	})

	g.Go(func() (err error) {
		opts := in.TransactionOptions
		opts.StartDate = cmpDates.StartDate
		opts.EndDate = cmpDates.EndDate
		opts.IsIncome = true

		currentIncome, err = uc.tr.SumTransactions(
			gCtx,
			in.UserID,
			opts,
		)
		return err
	})

	g.Go(func() (err error) {
		opts := in.TransactionOptions
		opts.StartDate = cmpDates.ComparisonStartDate
		opts.EndDate = cmpDates.ComparisonEndDate
		opts.IsIncome = true

		previousIncome, err = uc.tr.SumTransactions(
			gCtx,
			in.UserID,
			opts,
		)
		return err
	})

	g.Go(func() (err error) {
		opts := in.TransactionOptions
		opts.StartDate = cmpDates.StartDate
		opts.EndDate = cmpDates.EndDate
		opts.IsExpense = true

		currentExpense, err = uc.tr.SumTransactions(
			gCtx,
			in.UserID,
			opts,
		)
		return err
	})

	g.Go(func() (err error) {
		opts := in.TransactionOptions
		opts.StartDate = cmpDates.ComparisonStartDate
		opts.EndDate = cmpDates.ComparisonEndDate
		opts.IsExpense = true

		previousExpense, err = uc.tr.SumTransactions(
			gCtx,
			in.UserID,
			opts,
		)
		return err
	})

	if err := g.Wait(); err != nil {
		return nil, errs.New(err)
	}

	balancePercentageVariation := money.CalculatePercentageVariation(
		currentBalance,
		previousBalance,
	)

	incomePercentageVariation := money.CalculatePercentageVariation(
		currentIncome,
		previousIncome,
	)

	expensePercentageVariation := money.CalculatePercentageVariation(
		currentExpense,
		previousExpense,
	)

	out := &GetAccountsBalanceUseCaseOutput{
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
