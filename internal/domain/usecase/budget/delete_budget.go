package budget

import (
	"context"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/tx"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
	"github.com/google/uuid"
)

type DeleteBudgetUseCase struct {
	tx tx.TX
	br repo.BudgetRepo
}

func NewDeleteBudgetUseCase(
	tx tx.TX,
	br repo.BudgetRepo,
) *DeleteBudgetUseCase {
	return &DeleteBudgetUseCase{
		tx: tx,
		br: br,
	}
}

func (uc *DeleteBudgetUseCase) Execute(
	ctx context.Context,
	userID uuid.UUID,
) error {
	err := uc.tx.Do(ctx, func(ctx context.Context) error {
		if err := uc.br.DeleteBudgetCategories(ctx, userID); err != nil {
			return errs.New(err)
		}

		if err := uc.br.DeleteBudgets(ctx, userID); err != nil {
			return errs.New(err)
		}

		return nil
	})
	if err != nil {
		return errs.New(err)
	}

	return nil
}
