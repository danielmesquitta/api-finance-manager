package usecase

import (
	"context"
	"time"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/validator"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

type UpdateTransaction struct {
	v  *validator.Validator
	tr repo.TransactionRepo
}

func NewUpdateTransaction(
	v *validator.Validator,
	tr repo.TransactionRepo,
) *UpdateTransaction {
	return &UpdateTransaction{
		v:  v,
		tr: tr,
	}
}

type UpdateTransactionInput struct {
	ID              uuid.UUID  `json:"id"                validate:"required"`
	UserID          uuid.UUID  `json:"-"                 validate:"required"`
	Name            string     `json:"name"`
	Amount          int64      `json:"amount"`
	PaymentMethodID uuid.UUID  `json:"payment_method_id"`
	Date            time.Time  `json:"date"`
	AccountID       *uuid.UUID `json:"account_id"`
	InstitutionID   *uuid.UUID `json:"institution_id"`
	CategoryID      *uuid.UUID `json:"category_id"`
}

func (u *UpdateTransaction) Execute(
	ctx context.Context,
	in UpdateTransactionInput,
) error {
	if err := u.v.Validate(in); err != nil {
		return errs.New(err)
	}

	getTransactionRepoParams := repo.GetTransactionParams{
		ID:     in.ID,
		UserID: in.UserID,
	}

	transaction, err := u.tr.GetTransaction(ctx, getTransactionRepoParams)
	if err != nil {
		return errs.New(err)
	}
	if transaction == nil {
		return errs.ErrTransactionNotFound
	}

	params := repo.UpdateTransactionParams{}
	if err := copier.CopyWithOption(&params, in, copier.Option{IgnoreEmpty: true}); err != nil {
		return errs.New(err)
	}

	return u.tr.UpdateTransaction(ctx, params)
}
