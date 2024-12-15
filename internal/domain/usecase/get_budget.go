package usecase

import (
	"context"
	"time"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
	"github.com/google/uuid"
)

/*
 * @TODO calculate spent, available and available per day
 */

type GetBudgetUseCase struct {
	br repo.BudgetRepo
}

func NewGetBudgetUseCase(
	br repo.BudgetRepo,
) *GetBudgetUseCase {
	return &GetBudgetUseCase{
		br: br,
	}
}

type GetBudgetUseCaseInput struct {
	UserID uuid.UUID  `json:"-"     validate:"required"`
	Month  time.Month `json:"month" validate:"required,min=1,max=12"`
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
	budget, budgetCategories, categories, err := uc.br.GetBudgetWithCategoriesByUserID(
		ctx,
		in.UserID,
	)
	if err != nil {
		return nil, errs.New(err)
	}
	if budget == nil {
		return nil, errs.ErrBudgetNotFound
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
