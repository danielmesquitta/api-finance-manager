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
}

func NewGetBudgetCategory(
	v *validator.Validator,
	br repo.BudgetRepo,
	tr repo.TransactionRepo,
) *GetBudgetCategory {
	return &GetBudgetCategory{
		v:  v,
		br: br,
		tr: tr,
	}
}

type GetBudgetCategoryInput struct {
	UserID     uuid.UUID `json:"user_id"     validate:"required"`
	Date       string    `json:"date"        validate:"required"`
	CategoryID uuid.UUID `json:"category_id" validate:"required"`
}

type GetBudgetCategoryOutput struct {
	Spent        int64                                          `json:"spent"`
	Available    int64                                          `json:"available"`
	Transactions []entity.TransactionWithCategoryAndInstitution `json:"transactions"`
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

	transactions, err := uc.tr.ListTransactionsWithCategoriesAndInstitutions(
		ctx,
		in.UserID,
		repo.WithTransactionDateAfter(monthStart),
		repo.WithTransactionDateBefore(monthEnd),
		repo.WithTransactionIsIgnored(false),
		repo.WithTransactionCategory(in.CategoryID),
	)
	if err != nil {
		return nil, errs.New(err)
	}

	var spent int64
	for _, transaction := range transactions {
		if transaction.Amount > 0 {
			continue
		}

		spent -= transaction.Amount
	}

	available := budget.Amount - spent

	out := GetBudgetCategoryOutput{
		Spent:        spent,
		Available:    available,
		Transactions: transactions,
	}

	return &out, nil
}
