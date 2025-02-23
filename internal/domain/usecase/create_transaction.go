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

type CreateTransaction struct {
	v  *validator.Validator
	tr repo.TransactionRepo
}

func NewCreateTransaction(
	v *validator.Validator,
	tr repo.TransactionRepo,
) *CreateTransaction {
	return &CreateTransaction{
		v:  v,
		tr: tr,
	}
}

type CreateTransactionInput struct {
	Name            string     `json:"name"              validate:"required"`
	Amount          int64      `json:"amount"            validate:"required,gt=0"`
	PaymentMethodID uuid.UUID  `json:"payment_method_id" validate:"required"`
	Date            time.Time  `json:"date"              validate:"required"`
	CategoryID      *uuid.UUID `json:"category_id"       validate:"omitempty"`
	UserID          uuid.UUID  `json:"-"                 validate:"required"`
}

func (uc *CreateTransaction) Execute(
	ctx context.Context,
	in CreateTransactionInput,
) error {
	if err := uc.v.Validate(in); err != nil {
		return errs.New(err)
	}

	var params repo.CreateTransactionParams
	if err := copier.Copy(&params, in); err != nil {
		return errs.New(err)
	}

	if err := uc.tr.CreateTransaction(ctx, params); err != nil {
		return errs.New(err)
	}

	return nil
}
