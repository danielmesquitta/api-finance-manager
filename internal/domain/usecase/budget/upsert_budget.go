package budget

import (
	"context"
	"time"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/dateutil"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/tx"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/validator"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

type UpsertBudgetUseCase struct {
	v  *validator.Validator
	tx tx.TX
	br repo.BudgetRepo
	cr repo.TransactionCategoryRepo
}

func NewUpsertBudgetUseCase(
	v *validator.Validator,
	tx tx.TX,
	br repo.BudgetRepo,
	cr repo.TransactionCategoryRepo,
) *UpsertBudgetUseCase {
	return &UpsertBudgetUseCase{
		v:  v,
		tx: tx,
		br: br,
		cr: cr,
	}
}

type UpsertBudgetUseCaseCategoryInput struct {
	Amount     int64     `json:"amount"      validate:"required,gt=0"`
	CategoryID uuid.UUID `json:"category_id" validate:"required"`
}

type UpsertBudgetUseCaseInput struct {
	Amount     int64                              `json:"amount"     validate:"required,gt=0"`
	UserID     uuid.UUID                          `json:"-"          validate:"required"`
	Date       string                             `json:"date"       validate:"required"`
	Categories []UpsertBudgetUseCaseCategoryInput `json:"categories" validate:"dive"`
}

func (u *UpsertBudgetUseCase) Execute(
	ctx context.Context,
	in UpsertBudgetUseCaseInput,
) error {
	if err := u.validate(in); err != nil {
		return errs.New(err)
	}

	date, err := time.Parse(time.RFC3339, in.Date)
	if err != nil {
		return errs.ErrInvalidDate
	}

	monthStart := dateutil.ToMonthStart(date)

	categoryIDs := []uuid.UUID{}
	for _, c := range in.Categories {
		categoryIDs = append(categoryIDs, c.CategoryID)
	}

	g, gCtx := errgroup.WithContext(ctx)
	var (
		budget          *entity.Budget
		categoriesCount int64
	)

	g.Go(func() error {
		budget, err = u.br.GetBudget(gCtx, repo.GetBudgetParams{
			UserID: in.UserID,
			Date:   monthStart,
		})
		if err != nil {
			return errs.New(err)
		}
		return nil
	})

	g.Go(func() error {
		categoriesCount, err = u.cr.CountTransactionCategories(
			gCtx,
			repo.WithTransactionCategoryIDs(categoryIDs),
		)
		if err != nil {
			return errs.New(err)
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return errs.New(err)
	}

	if categoriesCount != int64(len(in.Categories)) {
		return errs.ErrCategoriesNotFound
	}

	err = u.tx.Do(ctx, func(ctx context.Context) error {
		if budget == nil || !budget.Date.Equal(monthStart) {
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

			if err := u.br.DeleteBudgetCategories(ctx, budget.ID); err != nil {
				return errs.New(err)
			}
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
	})

	if err != nil {
		return errs.New(err)
	}

	return nil
}

func (u *UpsertBudgetUseCase) validate(in UpsertBudgetUseCaseInput) error {
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
