package usecase

import (
	"context"
	"time"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/dateutil"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/validator"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
	"github.com/google/uuid"
)

type ListBudgetCategoryTransactions struct {
	v  *validator.Validator
	lt *ListTransactions
}

func NewListBudgetCategoryTransactions(
	v *validator.Validator,
	lt *ListTransactions,
) *ListBudgetCategoryTransactions {
	return &ListBudgetCategoryTransactions{
		v:  v,
		lt: lt,
	}
}

type ListBudgetCategoryTransactionsInput struct {
	PaginationInput
	UserID     uuid.UUID `json:"user_id"     validate:"required"`
	Date       time.Time `json:"date"        validate:"required"`
	CategoryID uuid.UUID `json:"category_id" validate:"required"`
}

func (uc *ListBudgetCategoryTransactions) Execute(
	ctx context.Context,
	in ListBudgetCategoryTransactionsInput,
) (*entity.PaginatedList[entity.FullTransaction], error) {
	if err := uc.v.Validate(in); err != nil {
		return nil, errs.New(err)
	}

	isIgnored := false
	categoryIDs := []uuid.UUID{in.CategoryID}

	transactions, err := uc.lt.Execute(ctx, ListTransactionsInput{
		PaginationInput: in.PaginationInput,
		UserID:          in.UserID,
		TransactionOptions: repo.TransactionOptions{
			StartDate:   dateutil.ToMonthStart(in.Date),
			EndDate:     dateutil.ToMonthEnd(in.Date),
			CategoryIDs: categoryIDs,
			IsExpense:   true,
			IsIgnored:   &isIgnored,
		},
	})
	if err != nil {
		return nil, err
	}

	for i, transaction := range transactions.Items {
		transactions.Items[i].Amount = -1 * transaction.Amount
	}

	return transactions, nil
}
