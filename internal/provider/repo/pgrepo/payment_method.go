package pgrepo

import (
	"context"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/db"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/db/query"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/db/sqlc"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"

	"github.com/jinzhu/copier"
)

type PaymentMethodPgRepo struct {
	db *db.DB
	qb *query.QueryBuilder
}

func NewPaymentMethodPgRepo(
	db *db.DB,
	qb *query.QueryBuilder,
) *PaymentMethodPgRepo {
	return &PaymentMethodPgRepo{
		db: db,
		qb: qb,
	}
}

func (r *PaymentMethodPgRepo) ListPaymentMethods(
	ctx context.Context,
	opts ...repo.PaymentMethodOption,
) ([]entity.PaymentMethod, error) {
	paymentMethods, err := r.qb.ListPaymentMethods(ctx, opts...)
	if err != nil {
		return nil, errs.New(err)
	}

	return paymentMethods, nil
}

func (r *PaymentMethodPgRepo) CountPaymentMethods(
	ctx context.Context,
	opts ...repo.PaymentMethodOption,
) (int64, error) {
	return r.qb.CountPaymentMethods(ctx, opts...)
}

func (r *PaymentMethodPgRepo) CreatePaymentMethods(
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

var _ repo.PaymentMethodRepo = &PaymentMethodPgRepo{}
