package usecase

import (
	"context"
	"time"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/dateutil"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/money"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/validator"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

type GetBudget struct {
	v  *validator.Validator
	br repo.BudgetRepo
	tr repo.TransactionRepo
}

func NewGetBudget(
	v *validator.Validator,
	br repo.BudgetRepo,
	tr repo.TransactionRepo,
) *GetBudget {
	return &GetBudget{
		v:  v,
		br: br,
		tr: tr,
	}
}

type GetBudgetInput struct {
	PaginationInput
	UserID uuid.UUID `json:"user_id" validate:"required"`
	Date   time.Time `json:"date"    validate:"required"`
}

type GetBudgetBudgetCategories struct {
	entity.BudgetCategory
	Spent     int64                      `json:"spent"`
	Available int64                      `json:"available"`
	Category  entity.TransactionCategory `json:"category"`
}

type GetBudgetOutput struct {
	entity.Budget
	Spent                              int64                       `json:"spent"`
	Available                          int64                       `json:"available"`
	AvailablePercentageVariation       int64                       `json:"available_percentage_variation"`
	AvailablePerDay                    int64                       `json:"available_per_day,omitempty"`
	AvailablePerDayPercentageVariation int64                       `json:"available_per_day_percentage_variation,omitempty"`
	ComparisonDates                    dateutil.ComparisonDates    `json:"comparison_dates"`
	BudgetCategories                   []GetBudgetBudgetCategories `json:"budget_categories"`
}

func (uc *GetBudget) Execute(
	ctx context.Context,
	in GetBudgetInput,
) (*GetBudgetOutput, error) {
	if err := uc.v.Validate(in); err != nil {
		return nil, errs.New(err)
	}

	cmpDates := dateutil.CalculateComparisonDates(
		dateutil.ToMonthStart(in.Date),
		dateutil.ToMonthEnd(in.Date),
	)

	budget, err := uc.br.GetBudget(ctx, repo.GetBudgetParams{
		UserID: in.UserID,
		Date:   cmpDates.StartDate,
	})
	if err != nil {
		return nil, errs.New(err)
	}
	if budget == nil {
		return nil, errs.ErrBudgetNotFound
	}

	g, gCtx := errgroup.WithContext(ctx)
	var (
		budgetCategories   []entity.BudgetCategory
		categories         []entity.TransactionCategory
		spentPreviousMonth int64
		spentByCategoryID  map[uuid.UUID]int64
	)

	g.Go(func() error {
		budgetCategories, categories, err = uc.br.ListBudgetCategories(
			gCtx,
			budget.ID,
		)
		return err
	})

	baseTransactionOpts := []repo.TransactionOption{
		repo.WithTransactionIsIgnored(false),
		repo.WithTransactionIsExpense(true),
	}

	g.Go(func() error {
		opts := append(
			baseTransactionOpts,
			repo.WithTransactionDateAfter(cmpDates.StartDate),
			repo.WithTransactionDateBefore(cmpDates.EndDate),
		)
		spentByCategoryID, err = uc.tr.SumTransactionsByCategory(
			gCtx,
			in.UserID,
			opts...,
		)
		return err
	})

	g.Go(func() error {
		opts := append(
			baseTransactionOpts,
			repo.WithTransactionDateAfter(cmpDates.ComparisonStartDate),
			repo.WithTransactionDateBefore(cmpDates.ComparisonEndDate),
		)
		spentPreviousMonth, err = uc.tr.SumTransactions(
			gCtx,
			in.UserID,
			opts...,
		)
		return err
	})

	if err := g.Wait(); err != nil {
		return nil, errs.New(err)
	}

	// Invert spentPreviousMonth and spentByCategoryID values to make them positive
	spentPreviousMonth = -1 * spentPreviousMonth
	for categoryID, spent := range spentByCategoryID {
		spentByCategoryID[categoryID] = -1 * spent
	}

	var spent int64
	for _, amount := range spentByCategoryID {
		spent += amount
	}

	available := budget.Amount - spent
	availablePreviousMonth := budget.Amount - spentPreviousMonth

	availablePercentageVariation := calculatePercentageVariation(
		available, availablePreviousMonth,
	)

	now := time.Now()
	isCurrentMonth := cmpDates.StartDate.Month() == now.Month() &&
		cmpDates.StartDate.Year() == now.Year()

	var availablePerDay, availablePerDayPercentageVariation int64
	if isCurrentMonth {
		availablePerDay = uc.calculateAvailablePerDay(
			available,
			dateutil.ToMonthEnd(cmpDates.EndDate),
			cmpDates.EndDate.Day(),
		)

		availablePreviousMonthPerDay := uc.calculateAvailablePerDay(
			availablePreviousMonth,
			dateutil.ToMonthEnd(cmpDates.ComparisonEndDate),
			cmpDates.ComparisonEndDate.Day(),
		)

		availablePerDayPercentageVariation = calculatePercentageVariation(
			availablePerDay, availablePreviousMonthPerDay,
		)
	}

	out := GetBudgetOutput{
		Budget:                             *budget,
		Spent:                              spent,
		Available:                          available,
		AvailablePercentageVariation:       availablePercentageVariation,
		AvailablePerDay:                    availablePerDay,
		AvailablePerDayPercentageVariation: availablePerDayPercentageVariation,
		ComparisonDates:                    *cmpDates,
		BudgetCategories:                   []GetBudgetBudgetCategories{},
	}

	categoriesByID := map[uuid.UUID]entity.TransactionCategory{}
	for _, category := range categories {
		categoriesByID[category.ID] = category
	}

	for _, budgetCategory := range budgetCategories {
		category := categoriesByID[budgetCategory.CategoryID]
		spent := spentByCategoryID[category.ID]
		available := budgetCategory.Amount - spent

		out.BudgetCategories = append(
			out.BudgetCategories,
			GetBudgetBudgetCategories{
				Spent:          spent,
				Available:      available,
				BudgetCategory: budgetCategory,
				Category:       category,
			},
		)
	}

	return &out, nil
}

func (uc *GetBudget) calculateAvailablePerDay(
	available int64,
	monthEnd time.Time,
	daysPassed int,
) int64 {
	daysInMonth := monthEnd.Day()
	daysLeft := daysInMonth - daysPassed + 1 // +1 to include today
	availablePerDay := money.FromCents(available) / float64(daysLeft)

	return money.ToCents(availablePerDay)
}
