package usecase

import (
	"context"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/validator"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
	"github.com/google/uuid"
)

/*
 * @TODO calculate spent, available and available per day
 */

type GetBudgetUseCase struct {
	v  *validator.Validator
	br repo.BudgetRepo
}

func NewGetBudgetUseCase(
	v *validator.Validator,
	br repo.BudgetRepo,
) *GetBudgetUseCase {
	return &GetBudgetUseCase{
		v:  v,
		br: br,
	}
}

type GetBudgetUseCaseInput struct {
	UserID uuid.UUID `json:"-"    validate:"required"`
	Date   string    `json:"date" validate:"required"`
}

type GetBudgetBudgetCategory struct {
	entity.BudgetCategory
	Spent    float64         `json:"spent,omitempty"`
	Category entity.Category `json:"category,omitempty"`
}

type GetBudgetUseCaseOutput struct {
	entity.Budget
	Spent            float64                   `json:"spent,omitempty"`
	Available        float64                   `json:"available,omitempty"`
	AvailablePerDay  float64                   `json:"available_per_day,omitempty"`
	BudgetCategories []GetBudgetBudgetCategory `json:"budget_categories,omitempty"`
}

func (uc *GetBudgetUseCase) Execute(
	ctx context.Context,
	in GetBudgetUseCaseInput,
) (*GetBudgetUseCaseOutput, error) {
	if err := uc.v.Validate(in); err != nil {
		return nil, errs.New(err)
	}

	monthStart, err := parseDateToMonthStart(in.Date)
	if err != nil {
		return nil, errs.New(err)
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

	budgetCategories, categories, err := uc.br.GetBudgetCategories(
		ctx,
		budget.ID,
	)
	if err != nil {
		return nil, errs.New(err)
	}

	out := GetBudgetUseCaseOutput{
		Budget: *budget,
	}

	categoriesByID := make(map[uuid.UUID]entity.Category)
	for _, category := range categories {
		categoriesByID[category.ID] = category
	}

	for _, budgetCategory := range budgetCategories {
		category := categoriesByID[budgetCategory.CategoryID]
		out.BudgetCategories = append(
			out.BudgetCategories,
			GetBudgetBudgetCategory{
				BudgetCategory: budgetCategory,
				Category:       category,
			},
		)
	}

	return &out, nil
}
