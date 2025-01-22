package usecase

import (
	"context"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

type GetTransaction struct {
	tr repo.TransactionRepo
}

func NewGetTransaction(
	tr repo.TransactionRepo,
) *GetTransaction {
	return &GetTransaction{
		tr: tr,
	}
}

type GetTransactionInput struct {
	UserID        uuid.UUID `json:"user_id"`
	TransactionID uuid.UUID `json:"transaction_id"`
}

func (uc *GetTransaction) Execute(
	ctx context.Context,
	in GetTransactionInput,
) (*entity.FullTransaction, error) {
	var repoParams repo.GetTransactionParams
	if err := copier.Copy(&repoParams, in); err != nil {
		return nil, errs.New(err)
	}

	transaction, err := uc.tr.GetTransaction(ctx, repoParams)
	if err != nil {
		return nil, errs.New(err)
	}
	if transaction == nil {
		return nil, errs.ErrTransactionNotFound
	}

	return transaction, nil
}
