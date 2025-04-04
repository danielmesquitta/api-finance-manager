package query

import (
	"context"
	"strings"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/db/schema"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
)

func (qb *QueryBuilder) ListPaymentMethods(
	ctx context.Context,
	opts ...repo.PaymentMethodOptions,
) ([]entity.PaymentMethod, error) {
	options := prepareOptions(opts...)

	query := goqu.
		From(schema.PaymentMethod.String()).
		Select(schema.PaymentMethod.All()).
		Where(goqu.I(schema.PaymentMethod.DeletedAt()).IsNull())

	whereExps, orderedExps := qb.buildPaymentMethodExpressions(options)

	query = qb.buildPaymentMethodQuery(
		query,
		options,
		whereExps,
		orderedExps,
	)

	var paymentMethods []entity.PaymentMethod
	if err := qb.Scan(ctx, query, &paymentMethods); err != nil {
		return nil, errs.New(err)
	}

	return paymentMethods, nil
}

func (qb *QueryBuilder) CountPaymentMethods(
	ctx context.Context,
	opts ...repo.PaymentMethodOptions,
) (int64, error) {
	options := prepareOptions(opts...)

	query := goqu.
		From(schema.PaymentMethod.String()).
		Select(goqu.COUNT(schema.PaymentMethod.All())).
		Where(goqu.I(schema.PaymentMethod.DeletedAt()).IsNull())

	whereExps, _ := qb.buildPaymentMethodExpressions(options)

	query = qb.buildPaymentMethodQuery(query, options, whereExps, nil)

	var count int64
	if err := qb.Scan(ctx, query, &count); err != nil {
		return 0, errs.New(err)
	}

	return count, nil
}

func (qb *QueryBuilder) buildPaymentMethodExpressions(
	options repo.PaymentMethodOptions,
) (whereExps []goqu.Expression, orderedExps []exp.OrderedExpression) {
	options.Search = strings.TrimSpace(options.Search)
	if options.Search != "" {
		searchExp, orderExp := qb.buildSearch(
			options.Search,
			schema.PaymentMethod.Name(),
		)
		whereExps = append(whereExps, searchExp)
		orderedExps = append(orderedExps, orderExp.Desc())
	}

	orderedExps = append(
		orderedExps,
		goqu.I(schema.PaymentMethod.Name()).Asc(),
	)

	return whereExps, orderedExps
}

func (qb *QueryBuilder) buildPaymentMethodQuery(
	query *goqu.SelectDataset,
	options repo.PaymentMethodOptions,
	whereExps []goqu.Expression,
	orderedExps []exp.OrderedExpression,
) *goqu.SelectDataset {
	if len(whereExps) > 0 {
		query = query.Where(whereExps...)
	}

	if len(orderedExps) > 0 {
		query = query.Order(orderedExps...)
	}

	if options.Limit > 0 {
		query = query.Limit(options.Limit)
	}

	if options.Offset > 0 {
		query = query.Offset(options.Offset)
	}

	return query
}
