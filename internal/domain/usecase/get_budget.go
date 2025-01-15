package usecase

import (
	"context"
	"time"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
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
	UserID uuid.UUID `json:"-"    validate:"required"`
	Date   string    `json:"date" validate:"required"`
}

type GetBudgetCategoryOutput struct {
	entity.BudgetCategory
	Spent     int64           `json:"spent"`
	Available int64           `json:"available"`
	Category  entity.Category `json:"category,omitempty"`
}

type GetBudgetOutput struct {
	entity.Budget
	Spent                              int64                     `json:"spent"`
	Available                          int64                     `json:"available"`
	AvailablePercentageVariation       int64                     `json:"available_percentage_variation"`
	AvailablePerDay                    int64                     `json:"available_per_day,omitempty"`
	AvailablePerDayPercentageVariation int64                     `json:"available_per_day_percentage_variation,omitempty"`
	ComparisonDate                     time.Time                 `json:"comparison_date"`
	BudgetCategories                   []GetBudgetCategoryOutput `json:"budget_categories"`
}

func (uc *GetBudget) Execute(
	ctx context.Context,
	in GetBudgetInput,
) (*GetBudgetOutput, error) {
	if err := uc.v.Validate(in); err != nil {
		return nil, errs.New(err)
	}

	date, err := time.Parse(time.RFC3339, in.Date)
	if err != nil {
		return nil, errs.ErrInvalidDate
	}

	now := time.Now()
	isCurrentMonth := date.Year() == now.Year() && date.Month() == now.Month()

	monthStart := toMonthStart(date)
	monthEnd := toMonthEnd(date)
	monthSameDayAsToday := toMonthDay(monthStart, now.Day())

	previousMonthStart := monthStart.AddDate(0, -1, 0)
	previousMonthEnd := toMonthEnd(previousMonthStart)
	previousMonthSameDayAsToday := toMonthDay(previousMonthStart, now.Day())

	comparisonDate := previousMonthEnd
	if isCurrentMonth {
		comparisonDate = previousMonthSameDayAsToday
	}

	budget, err := uc.br.GetBudget(ctx, repo.GetBudgetParams{
		UserID: in.UserID,
		Date:   monthStart,
	})
	if err != nil {
		return nil, errs.New(err)
	}
	if budget == nil {
		return nil, errs.ErrBudgetNotFound
	}

	g, gCtx := errgroup.WithContext(ctx)
	var budgetCategories []entity.BudgetCategory
	var categories []entity.Category
	var transactions, previousMonthTransactions []entity.TransactionWithCategoryAndInstitution

	g.Go(func() error {
		budgetCategories, categories, err = uc.br.ListBudgetCategories(
			gCtx,
			budget.ID,
		)
		return err
	})

	g.Go(func() error {
		transactions, err = uc.tr.ListTransactionsWithCategoriesAndInstitutions(
			gCtx,
			in.UserID,
			repo.WithTransactionDateAfter(monthStart),
			repo.WithTransactionDateBefore(monthEnd),
		)
		return err
	})

	g.Go(func() error {
		previousMonthTransactions, err = uc.tr.ListTransactionsWithCategoriesAndInstitutions(
			gCtx,
			in.UserID,
			repo.WithTransactionDateAfter(previousMonthStart),
			repo.WithTransactionDateBefore(comparisonDate),
		)
		return err
	})

	if err := g.Wait(); err != nil {
		return nil, errs.New(err)
	}

	var spent int64
	spentByCategoryID := map[uuid.UUID]int64{}
	for _, transaction := range transactions {
		if transaction.Amount > 0 {
			continue
		}
		if _, ok := spentByCategoryID[transaction.CategoryID]; !ok {
			spentByCategoryID[transaction.CategoryID] = 0
		}
		spent -= transaction.Amount
		spentByCategoryID[transaction.CategoryID] -= transaction.Amount
	}

	var spentPreviousMonth int64
	for _, transaction := range previousMonthTransactions {
		if transaction.Amount > 0 {
			continue
		}
		spentPreviousMonth -= transaction.Amount
	}

	available := budget.Amount - spent
	availablePreviousMonth := budget.Amount - spentPreviousMonth

	availablePercentageVariation := money.FromPercentage(
		1 - (float64(available) / float64(availablePreviousMonth)),
	)

	var availablePerDay, availablePerDayPercentageVariation int64
	if isCurrentMonth {
		availablePerDay = uc.calculateAvailablePerDay(
			available,
			monthEnd,
			monthSameDayAsToday.Day(),
		)

		availablePreviousMonthPerDay := uc.calculateAvailablePerDay(
			availablePreviousMonth,
			previousMonthEnd,
			previousMonthSameDayAsToday.Day(),
		)

		availablePerDayPercentageVariation = money.FromPercentage(
			1 - (float64(availablePerDay) / float64(availablePreviousMonthPerDay)),
		)
	}

	out := GetBudgetOutput{
		Budget:                             *budget,
		Spent:                              spent,
		Available:                          available,
		AvailablePercentageVariation:       availablePercentageVariation,
		AvailablePerDay:                    availablePerDay,
		AvailablePerDayPercentageVariation: availablePerDayPercentageVariation,
		ComparisonDate:                     comparisonDate,
		BudgetCategories:                   []GetBudgetCategoryOutput{},
	}

	categoriesByID := make(map[uuid.UUID]entity.Category)
	for _, category := range categories {
		categoriesByID[category.ID] = category
	}

	for _, budgetCategory := range budgetCategories {
		category := categoriesByID[budgetCategory.CategoryID]
		spent := spentByCategoryID[category.ID]
		available := budgetCategory.Amount - spent

		out.BudgetCategories = append(
			out.BudgetCategories,
			GetBudgetCategoryOutput{
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
