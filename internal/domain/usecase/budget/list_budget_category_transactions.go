package budget

import (
	"context"
	"time"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/usecase"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/usecase/transaction"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/dateutil"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/validator"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
	"github.com/google/uuid"
)

type ListBudgetCategoryTransactionsUseCase struct {
	v  *validator.Validator
	lt *transaction.ListTransactionsUseCase
}

func NewListBudgetCategoryTransactionsUseCase(
	v *validator.Validator,
	lt *transaction.ListTransactionsUseCase,
) *ListBudgetCategoryTransactionsUseCase {
	return &ListBudgetCategoryTransactionsUseCase{
		v:  v,
		lt: lt,
	}
}

type ListBudgetCategoryTransactionsUseCaseInput struct {
	usecase.PaginationInput
	UserID     uuid.UUID `json:"user_id"     validate:"required"`
	Date       time.Time `json:"date"        validate:"required"`
	CategoryID uuid.UUID `json:"category_id" validate:"required"`
}

func (uc *ListBudgetCategoryTransactionsUseCase) Execute(
	ctx context.Context,
	in ListBudgetCategoryTransactionsUseCaseInput,
) (*entity.PaginatedList[entity.FullTransaction], error) {
	if err := uc.v.Validate(in); err != nil {
		return nil, errs.New(err)
	}

	isIgnored := false
	categoryIDs := []uuid.UUID{in.CategoryID}

	transactions, err := uc.lt.Execute(
		ctx,
		transaction.ListTransactionsUseCaseInput{
			PaginationInput: in.PaginationInput,
			UserID:          in.UserID,
			TransactionOptions: repo.TransactionOptions{
				StartDate:   dateutil.ToMonthStart(in.Date),
				EndDate:     dateutil.ToMonthEnd(in.Date),
				CategoryIDs: categoryIDs,
				IsExpense:   true,
				IsIgnored:   &isIgnored,
			},
		},
	)
	if err != nil {
		return nil, err
	}

	for i, trans := range transactions.Items {
		transactions.Items[i].Amount = -1 * trans.Amount
	}

	return transactions, nil
}
