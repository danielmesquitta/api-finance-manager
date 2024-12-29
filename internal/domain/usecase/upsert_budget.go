package usecase

import (
	"context"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/tx"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/validator"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
	"github.com/google/uuid"
)

type UpsertBudget struct {
	v  *validator.Validator
	tx tx.TX
	br repo.BudgetRepo
}

func NewUpsertBudget(
	v *validator.Validator,
	tx tx.TX,
	br repo.BudgetRepo,
) *UpsertBudget {
	return &UpsertBudget{
		v:  v,
		tx: tx,
		br: br,
	}
}

type UpsertBudgetCategoryInput struct {
	Amount     int64     `json:"amount,omitempty"      validate:"required,gt=0"`
	CategoryID uuid.UUID `json:"category_id,omitempty" validate:"required"`
}

type UpsertBudgetInput struct {
	Amount     int64                       `json:"amount,omitempty"     validate:"required,gt=0"`
	UserID     uuid.UUID                   `json:"user_id,omitempty"    validate:"required"`
	Date       string                      `json:"date,omitempty"       validate:"required"`
	Categories []UpsertBudgetCategoryInput `json:"categories,omitempty" validate:"dive"`
}

func (u *UpsertBudget) Execute(
	ctx context.Context,
	in UpsertBudgetInput,
) error {
	if err := u.validate(in); err != nil {
		return errs.New(err)
	}

	monthStart, err := parseDateToMonthStart(in.Date)
	if err != nil {
		return errs.New(err)
	}

	budget, err := u.br.GetBudget(ctx, repo.GetBudgetParams{
		UserID: in.UserID,
		Date:   monthStart,
	})
	if err != nil {
		return errs.New(err)
	}

	if err := u.tx.Do(ctx, func(ctx context.Context) error {
		if budgetDoesNotExists := budget == nil; budgetDoesNotExists {
			budget, err = u.br.CreateBudget(ctx, repo.CreateBudgetParams{
				Date:   monthStart,
				Amount: in.Amount,
				UserID: in.UserID,
			})
			if err != nil {
				return errs.New(err)
			}
		} else {
			if err := u.br.UpdateBudget(ctx, repo.UpdateBudgetParams{
				Date:   monthStart,
				UserID: in.UserID,
				Amount: in.Amount,
			}); err != nil {
				return errs.New(err)
			}
		}

		if err := u.br.DeleteBudgetCategories(ctx, in.UserID); err != nil {
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

func (u *UpsertBudget) validate(in UpsertBudgetInput) error {
	if err := u.v.Validate(in); err != nil {
		return errs.New(err)
	}

	var totalAmount int64
	for _, c := range in.Categories {
		totalAmount += c.Amount
	}
	if totalAmount > in.Amount {
		return errs.ErrInvalidTotalBudgetCategoryAmount
	}

	return nil
}
