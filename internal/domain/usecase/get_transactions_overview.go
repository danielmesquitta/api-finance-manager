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

type GetTransactionsOverview struct {
	v  *validator.Validator
	tr repo.TransactionRepo
}

func NewGetTransactionsOverview(
	v *validator.Validator,
	tr repo.TransactionRepo,
) *GetTransactionsOverview {
	return &GetTransactionsOverview{
		v:  v,
		tr: tr,
	}
}

type GetTransactionsOverviewInput struct {
	Date   time.Time `json:"date,omitempty"    validate:"required"`
	UserID uuid.UUID `json:"user_id,omitempty" validate:"required"`
}

type GetTransactionsOverviewOutput struct {
	ComparisonDates                   ComparisonDates `json:"comparison_dates"`
	CurrentBalance                    int64           `json:"current_balance,omitempty"`
	CurrentBalancePercentageVariation int64           `json:"current_balance_percentage_variation,omitempty"`
	MonthlyBalance                    int64           `json:"monthly_balance,omitempty"`
	MonthlyBalancePercentageVariation int64           `json:"monthly_balance_percentage_variation,omitempty"`
}

func (uc *GetTransactionsOverview) Execute(
	ctx context.Context,
	in GetTransactionsOverviewInput,
) (*GetTransactionsOverviewOutput, error) {
	if err := uc.v.Validate(in); err != nil {
		return nil, errs.New(err)
	}

	comparisonDates := calculateComparisonDates(in.Date)

	baseOpts := []repo.TransactionOption{
		repo.WithTransactionIsIgnored(false),
	}

	g, gCtx := errgroup.WithContext(ctx)
	var currentMonthBalance, previousMonthBalance int64

	g.Go(func() error {
		var err error
		opts := append(
			baseOpts,
			repo.WithTransactionDateAfter(comparisonDates.MonthStart),
			repo.WithTransactionDateBefore(
				comparisonDates.MonthComparisonEndDate,
			),
		)
		currentMonthBalance, err = uc.tr.SumTransactions(
			gCtx,
			in.UserID,
			opts...,
		)
		return err
	})

	g.Go(func() error {
		var err error
		opts := append(
			baseOpts,
			repo.WithTransactionDateAfter(comparisonDates.PreviousMonthStart),
			repo.WithTransactionDateBefore(
				comparisonDates.PreviousMonthComparisonEndDate,
			),
		)
		previousMonthBalance, err = uc.tr.SumTransactions(
			gCtx,
			in.UserID,
			opts...,
		)
		return err
	})

	if err := g.Wait(); err != nil {
		return nil, errs.New(err)
	}

	monthlyBalancePercentageVariation := money.FromPercentage(
		1 - (float64(currentMonthBalance) / float64(previousMonthBalance)),
	)

	out := &GetTransactionsOverviewOutput{
		ComparisonDates:                   *comparisonDates,
		MonthlyBalance:                    currentMonthBalance,
		MonthlyBalancePercentageVariation: monthlyBalancePercentageVariation,
	}

	return out, nil
}
