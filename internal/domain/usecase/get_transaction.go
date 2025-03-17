package usecase

import (
	"context"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
	"github.com/google/uuid"
)

type GetTransactionByID struct {
	tr repo.TransactionRepo
}

func NewGetTransaction(
	tr repo.TransactionRepo,
) *GetTransactionByID {
	return &GetTransactionByID{
		tr: tr,
	}
}

type GetTransactionInput struct {
	ID     uuid.UUID `json:"transaction_id"`
	UserID uuid.UUID `json:"user_id"`
}

func (uc *GetTransactionByID) Execute(
	ctx context.Context,
	in GetTransactionInput,
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
