package budget

import (
	"context"
	"time"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/dateutil"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/money"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/ptr"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/validator"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

type GetBudgetUseCase struct {
	v  *validator.Validator
	br repo.BudgetRepo
	tr repo.TransactionRepo
}

func NewGetBudgetUseCase(
	v *validator.Validator,
	br repo.BudgetRepo,
	tr repo.TransactionRepo,
) *GetBudgetUseCase {
	return &GetBudgetUseCase{
		v:  v,
		br: br,
		tr: tr,
	}
}

type GetBudgetUseCaseInput struct {
	UserID uuid.UUID `json:"user_id" validate:"required"`
	Date   time.Time `json:"date"    validate:"required"`
}

type GetBudgetUseCaseBudgetCategories struct {
	entity.BudgetCategory
	Spent     int64                      `json:"spent"`
	Available int64                      `json:"available"`
	Category  entity.TransactionCategory `json:"category"`
}

type GetBudgetUseCaseOutput struct {
	entity.Budget
	Spent                              int64                              `json:"spent"`
	Available                          int64                              `json:"available"`
	AvailablePercentageVariation       int64                              `json:"available_percentage_variation"`
	AvailablePerDay                    int64                              `json:"available_per_day,omitempty"`
	AvailablePerDayPercentageVariation int64                              `json:"available_per_day_percentage_variation,omitempty"`
	ComparisonDates                    dateutil.ComparisonDates           `json:"comparison_dates"`
	BudgetCategories                   []GetBudgetUseCaseBudgetCategories `json:"budget_categories"`
}

func (uc *GetBudgetUseCase) Execute(
	ctx context.Context,
	in GetBudgetUseCaseInput,
) (*GetBudgetUseCaseOutput, error) {
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

	baseTransactionOpts := repo.TransactionOptions{
		IsIgnored: ptr.New(false),
		IsExpense: true,
	}

	g.Go(func() error {
		opts := baseTransactionOpts
		opts.StartDate = cmpDates.StartDate
		opts.EndDate = cmpDates.EndDate

		spentByCategoryID, err = uc.tr.SumTransactionsByCategory(
			gCtx,
			in.UserID,
			opts,
		)
		return err
	})

	g.Go(func() error {
		opts := baseTransactionOpts
		opts.StartDate = cmpDates.ComparisonStartDate
		opts.EndDate = cmpDates.ComparisonEndDate

		spentPreviousMonth, err = uc.tr.SumTransactions(
			gCtx,
			in.UserID,
			opts,
		)
		return err
	})

	if err := g.Wait(); err != nil {
		return nil, errs.New(err)
	}

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

	availablePercentageVariation := money.CalculatePercentageVariation(
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

		availablePerDayPercentageVariation = money.CalculatePercentageVariation(
			availablePerDay, availablePreviousMonthPerDay,
		)
	}

	out := GetBudgetUseCaseOutput{
		Budget:                             *budget,
		Spent:                              spent,
		Available:                          available,
		AvailablePercentageVariation:       availablePercentageVariation,
		AvailablePerDay:                    availablePerDay,
		AvailablePerDayPercentageVariation: availablePerDayPercentageVariation,
		ComparisonDates:                    *cmpDates,
		BudgetCategories:                   []GetBudgetUseCaseBudgetCategories{},
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
			GetBudgetUseCaseBudgetCategories{
				Spent:          spent,
				Available:      available,
				BudgetCategory: budgetCategory,
				Category:       category,
			},
		)
	}

	return &out, nil
}

func (uc *GetBudgetUseCase) calculateAvailablePerDay(
	available int64,
	monthEnd time.Time,
	daysPassed int,
) int64 {
	daysInMonth := monthEnd.Day()
	daysLeft := daysInMonth - daysPassed + 1 // +1 to include today
	availablePerDay := money.FromCents(available) / float64(daysLeft)

	return money.ToCents(availablePerDay)
}
