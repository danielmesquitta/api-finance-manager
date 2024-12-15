package usecase

import (
	"context"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/tx"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/validator"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
	"github.com/google/uuid"
)

type UpsertBudgetUseCase struct {
	v  *validator.Validator
	tx tx.TX
	br repo.BudgetRepo
}

func NewUpsertBudgetUseCase(
	v *validator.Validator,
	tx tx.TX,
	br repo.BudgetRepo,
) *UpsertBudgetUseCase {
	return &UpsertBudgetUseCase{
		v:  v,
		tx: tx,
		br: br,
	}
}

type UpsertBudgetCategoryInput struct {
	Amount     float64   `json:"amount,omitempty"      validate:"required,gt=0"`
	CategoryID uuid.UUID `json:"category_id,omitempty" validate:"required"`
}

type UpsertBudgetUseCaseInput struct {
	Amount     float64                     `json:"amount,omitempty"     validate:"required,gt=0"`
	UserID     uuid.UUID                   `json:"-"                    validate:"required"`
	Categories []UpsertBudgetCategoryInput `json:"categories,omitempty" validate:"dive"`
}

func (u *UpsertBudgetUseCase) Execute(
	ctx context.Context,
	in UpsertBudgetUseCaseInput,
) error {
	if err := u.validate(in); err != nil {
		return errs.New(err)
	}

	budget, err := u.br.GetBudgetByUserID(ctx, in.UserID)
	if err != nil {
		return errs.New(err)
	}

	if err := u.tx.Do(ctx, func(ctx context.Context) error {
		if budgetDoesNotExists := budget == nil; budgetDoesNotExists {
			budget, err = u.br.CreateBudget(ctx, repo.CreateBudgetParams{
				Amount: in.Amount,
				UserID: in.UserID,
			})
			if err != nil {
				return errs.New(err)
			}
		} else {
			if err := u.br.UpdateBudget(ctx, repo.UpdateBudgetParams{
				UserID: in.UserID,
				Amount: in.Amount,
			}); err != nil {
				return errs.New(err)
			}
		}

		if err := u.br.DeleteBudgetCategoriesByBudgetID(ctx, budget.ID); err != nil {
			return errs.New(err)
		}

		var categories []repo.CreateBudgetCategoriesParams
		for _, c := range in.Categories {
			categories = append(categories, repo.CreateBudgetCategoriesParams{
				Amount:     c.Amount,
				BudgetID:   budget.ID,
				CategoryID: c.CategoryID,
			})
		}

		if err = u.br.CreateBudgetCategories(ctx, categories); err != nil {
			return errs.New(err)
		}

		return nil
	}); err != nil {
		return errs.New(err)
	}

	return nil
}

func (u *UpsertBudgetUseCase) validate(in UpsertBudgetUseCaseInput) error {
	if err := u.v.Validate(in); err != nil {
		return errs.New(err)
	}

	totalAmount := 0.0
	for _, c := range in.Categories {
		totalAmount += c.Amount
	}
	if totalAmount > in.Amount {
		return errs.ErrInvalidTotalBudgetCategoryAmount
	}

	return nil
}
