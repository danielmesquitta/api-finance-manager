package usecase

import (
	"context"
	"time"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/validator"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"golang.org/x/sync/errgroup"
)

type CreateTransaction struct {
	v   *validator.Validator
	tr  repo.TransactionRepo
	ur  repo.UserRepo
	tcr repo.TransactionCategoryRepo
	pmr repo.PaymentMethodRepo
}

func NewCreateTransaction(
	v *validator.Validator,
	tr repo.TransactionRepo,
	ur repo.UserRepo,
	tcr repo.TransactionCategoryRepo,
	pmr repo.PaymentMethodRepo,
) *CreateTransaction {
	return &CreateTransaction{
		v:   v,
		tr:  tr,
		ur:  ur,
		tcr: tcr,
		pmr: pmr,
	}
}

type CreateTransactionInput struct {
	UserID          uuid.UUID  `json:"-"                 validate:"required"`
	Name            string     `json:"name"              validate:"required"`
	Amount          int64      `json:"amount"            validate:"required"`
	PaymentMethodID uuid.UUID  `json:"payment_method_id" validate:"required"`
	Date            time.Time  `json:"date"              validate:"required"`
	CategoryID      *uuid.UUID `json:"category_id"       validate:"omitempty"`
}

func (uc *CreateTransaction) Execute(
	ctx context.Context,
	in CreateTransactionInput,
) error {
	if err := uc.v.Validate(in); err != nil {
		return errs.New(err)
	}

	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		user, err := uc.ur.GetUserByID(gCtx, in.UserID)
		if err != nil {
			return err
		}
		if user == nil {
			return errs.ErrUserNotFound
		}
		return nil
	})

	g.Go(func() error {
		paymentMethod, err := uc.pmr.GetPaymentMethodByID(
			ctx,
			in.PaymentMethodID,
		)
		if err != nil {
			return err
		}
		if paymentMethod == nil {
			return errs.ErrPaymentMethodNotFound
		}
		return nil
	})

	g.Go(func() error {
		var (
			category *entity.TransactionCategory
			err      error
		)
		if in.CategoryID == nil {
			category, err = uc.tcr.GetDefaultTransactionCategory(ctx)
		} else {
			category, err = uc.tcr.GetTransactionCategoryByID(ctx, *in.CategoryID)
		}
		if err != nil {
			return errs.New(err)
		}
		if category == nil {
			return errs.ErrCategoryNotFound
		}
		in.CategoryID = &category.ID
		return nil
	})

	if err := g.Wait(); err != nil {
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
