package usecase

import (
	"context"
	"time"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/validator"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
	"github.com/google/uuid"
)

type GetBudgetCategory struct {
	v  *validator.Validator
	br repo.BudgetRepo
	tr repo.TransactionRepo
	cr repo.CategoryRepo
}

func NewGetBudgetCategory(
	v *validator.Validator,
	br repo.BudgetRepo,
	tr repo.TransactionRepo,
	cr repo.CategoryRepo,
) *GetBudgetCategory {
	return &GetBudgetCategory{
		v:  v,
		br: br,
		tr: tr,
		cr: cr,
	}
}

type GetBudgetCategoryInput struct {
	UserID     uuid.UUID `json:"user_id"     validate:"required"`
	Date       string    `json:"date"        validate:"required"`
	CategoryID uuid.UUID `json:"category_id" validate:"required"`
}

type GetBudgetCategoryOutput struct {
	entity.TransactionCategory
	Amount    int64 `json:"amount"`
	Spent     int64 `json:"spent"`
	Available int64 `json:"available"`
}

func (uc *GetBudgetCategory) Execute(
	ctx context.Context,
	in GetBudgetCategoryInput,
) (*GetBudgetCategoryOutput, error) {
	if err := uc.v.Validate(in); err != nil {
		return nil, errs.New(err)
	}

	date, err := time.Parse(time.RFC3339, in.Date)
	if err != nil {
		return nil, errs.ErrInvalidDate
	}

	monthStart := toMonthStart(date)
	monthEnd := toMonthEnd(date)

	budgetCategory, category, err := uc.br.GetBudgetCategory(
		ctx,
		repo.GetBudgetCategoryParams{
			UserID: in.UserID,
			Date:   monthStart,
		},
	)
	if err != nil {
		return nil, errs.New(err)
	}
	if budgetCategory == nil || category == nil {
		return nil, errs.ErrBudgetCategoryNotFound
	}

	transactionOpts := []repo.TransactionOption{
		repo.WithTransactionDateAfter(monthStart),
		repo.WithTransactionDateBefore(monthEnd),
		repo.WithTransactionCategories(in.CategoryID),
		repo.WithTransactionIsIgnored(false),
		repo.WithTransactionIsExpense(true),
	}

	spent, err := uc.tr.SumTransactions(
		ctx,
		in.UserID,
		transactionOpts...,
	)
	if err != nil {
		return nil, errs.New(err)
	}

	amount := budgetCategory.Amount
	available := amount - spent

	out := GetBudgetCategoryOutput{
		TransactionCategory: *category,
		Amount:              budgetCategory.Amount,
		Spent:               spent,
		Available:           available,
	}

	return &out, nil
}
