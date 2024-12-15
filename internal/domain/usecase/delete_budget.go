package usecase

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
	budget, err := uc.br.GetBudgetByUserID(ctx, userID)
	if err != nil {
		return err
	}
	if budget == nil {
		return errs.ErrBudgetNotFound
	}

	err = uc.tx.Do(ctx, func(ctx context.Context) error {
		if err := uc.br.DeleteBudgetCategoriesByBudgetID(ctx, budget.ID); err != nil {
			return errs.New(err)
		}

		if err := uc.br.DeleteBudgetByID(ctx, budget.ID); err != nil {
			return errs.New(err)
		}

		return nil
	})
	if err != nil {
		return errs.New(err)
	}

	return nil
}
