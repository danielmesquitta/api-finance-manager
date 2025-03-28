package budget

import (
	"context"
	"time"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/dateutil"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/ptr"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/validator"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
	"github.com/google/uuid"
)

type GetBudgetCategoryUseCase struct {
	v  *validator.Validator
	br repo.BudgetRepo
	tr repo.TransactionRepo
	cr repo.TransactionCategoryRepo
}

func NewGetBudgetCategoryUseCase(
	v *validator.Validator,
	br repo.BudgetRepo,
	tr repo.TransactionRepo,
	cr repo.TransactionCategoryRepo,
) *GetBudgetCategoryUseCase {
	return &GetBudgetCategoryUseCase{
		v:  v,
		br: br,
		tr: tr,
		cr: cr,
	}
}

type GetBudgetCategoryUseCaseInput struct {
	UserID     uuid.UUID `json:"user_id"     validate:"required"`
	Date       string    `json:"date"        validate:"required"`
	CategoryID uuid.UUID `json:"category_id" validate:"required"`
}

type GetBudgetCategoryUseCaseOutput struct {
	entity.TransactionCategory
	Amount    int64 `json:"amount"`
	Spent     int64 `json:"spent"`
	Available int64 `json:"available"`
}

func (uc *GetBudgetCategoryUseCase) Execute(
	ctx context.Context,
	in GetBudgetCategoryUseCaseInput,
) (*GetBudgetCategoryUseCaseOutput, error) {
	if err := uc.v.Validate(in); err != nil {
		return nil, errs.New(err)
	}

	date, err := time.Parse(time.RFC3339, in.Date)
	if err != nil {
		return nil, errs.ErrInvalidDate
	}

	monthStart := dateutil.ToMonthStart(date)
	monthEnd := dateutil.ToMonthEnd(date)

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

	transactionOpts := repo.TransactionOptions{
		StartDate:   monthStart,
		EndDate:     monthEnd,
		CategoryIDs: []uuid.UUID{in.CategoryID},
		IsIgnored:   ptr.New(false),
		IsExpense:   true,
	}

	spent, err := uc.tr.SumTransactions(
		ctx,
		in.UserID,
		transactionOpts,
	)
	if err != nil {
		return nil, errs.New(err)
	}

	amount := budgetCategory.Amount
	available := amount - spent

	out := GetBudgetCategoryUseCaseOutput{
		TransactionCategory: *category,
		Amount:              budgetCategory.Amount,
		Spent:               spent,
		Available:           available,
	}

	return &out, nil
}
