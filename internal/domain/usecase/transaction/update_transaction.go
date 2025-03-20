package transaction

import (
	"context"
	"time"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/validator"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

type UpdateTransactionUseCase struct {
	v  *validator.Validator
	tr repo.TransactionRepo
}

func NewUpdateTransactionUseCase(
	v *validator.Validator,
	tr repo.TransactionRepo,
) *UpdateTransactionUseCase {
	return &UpdateTransactionUseCase{
		v:  v,
		tr: tr,
	}
}

type UpdateTransactionUseCaseInput struct {
	ID              uuid.UUID `json:"id"                validate:"required"`
	UserID          uuid.UUID `json:"-"                 validate:"required"`
	Name            string    `json:"name"`
	Amount          int64     `json:"amount"`
	PaymentMethodID uuid.UUID `json:"payment_method_id"`
	Date            time.Time `json:"date"`
	AccountID       uuid.UUID `json:"account_id"`
	InstitutionID   uuid.UUID `json:"institution_id"`
	CategoryID      uuid.UUID `json:"category_id"`
}

func (u *UpdateTransactionUseCase) Execute(
	ctx context.Context,
	in UpdateTransactionUseCaseInput,
) error {
	if err := u.v.Validate(in); err != nil {
		return errs.New(err)
	}

	transaction, err := u.tr.GetTransactionByID(ctx, in.ID)
	if err != nil {
		return errs.New(err)
	}
	if transaction == nil || transaction.UserID != in.UserID {
		return errs.ErrTransactionNotFound
	}

	params := repo.UpdateTransactionParams{}
	if err := copier.Copy(&params, transaction); err != nil {
		return errs.New(err)
	}

	if err := copier.CopyWithOption(&params, in, copier.Option{IgnoreEmpty: true}); err != nil {
		return errs.New(err)
	}

	return u.tr.UpdateTransaction(ctx, params)
}
