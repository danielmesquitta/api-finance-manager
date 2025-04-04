package pgrepo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/db"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/db/sqlc"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
	"github.com/google/uuid"

	"github.com/jinzhu/copier"
)

type PaymentMethodRepo struct {
	db *db.DB
}

func NewPaymentMethodRepo(
	db *db.DB,
) *PaymentMethodRepo {
	return &PaymentMethodRepo{
		db: db,
	}
}

func (r *PaymentMethodRepo) ListPaymentMethods(
	ctx context.Context,
	opts ...repo.PaymentMethodOptions,
) ([]entity.PaymentMethod, error) {
	paymentMethods, err := r.db.ListPaymentMethods(ctx, opts...)
	if err != nil {
		return nil, errs.New(err)
	}

	return paymentMethods, nil
}

func (r *PaymentMethodRepo) CountPaymentMethods(
	ctx context.Context,
	opts ...repo.PaymentMethodOptions,
) (int64, error) {
	return r.db.CountPaymentMethods(ctx, opts...)
}

func (r *PaymentMethodRepo) CreatePaymentMethods(
	ctx context.Context,
	params []repo.CreatePaymentMethodsParams,
) error {
	dbParams := make([]sqlc.CreatePaymentMethodsParams, len(params))
	if err := copier.Copy(&dbParams, params); err != nil {
		return errs.New(err)
	}

	tx := r.db.UseTx(ctx)
	_, err := tx.CreatePaymentMethods(ctx, dbParams)
	if err != nil {
		return errs.New(err)
	}

	return nil
}

func (r *PaymentMethodRepo) GetPaymentMethodByID(
	ctx context.Context,
	id uuid.UUID,
) (*entity.PaymentMethod, error) {
	paymentMethod, err := r.db.GetPaymentMethodByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, errs.New(err)
	}

	result := entity.PaymentMethod{}
	if err := copier.Copy(&result, paymentMethod); err != nil {
		return nil, errs.New(err)
	}

	return &result, nil
}

var _ repo.PaymentMethodRepo = (*PaymentMethodRepo)(nil)
