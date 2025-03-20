package transaction

import (
	"context"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
	"github.com/google/uuid"
)

type GetTransactionUseCase struct {
	tr repo.TransactionRepo
}

func NewGetTransactionUseCase(
	tr repo.TransactionRepo,
) *GetTransactionUseCase {
	return &GetTransactionUseCase{
		tr: tr,
	}
}

type GetTransactionUseCaseInput struct {
	ID     uuid.UUID `json:"transaction_id"`
	UserID uuid.UUID `json:"user_id"`
}

func (uc *GetTransactionUseCase) Execute(
	ctx context.Context,
	in GetTransactionUseCaseInput,
) (*entity.FullTransaction, error) {
	transaction, err := uc.tr.GetTransactionByID(ctx, in.ID)
	if err != nil {
		return nil, errs.New(err)
	}
	if transaction == nil || transaction.UserID != in.UserID {
		return nil, errs.ErrTransactionNotFound
	}

	return transaction, nil
}
