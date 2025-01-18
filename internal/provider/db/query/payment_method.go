package query

import (
	"context"
	"strings"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/georgysavva/scany/v2/pgxscan"
)

func (qb *QueryBuilder) ListPaymentMethods(
	ctx context.Context,
	opts ...repo.ListPaymentMethodsOption,
) ([]entity.PaymentMethod, error) {
	options := repo.ListPaymentMethodsOptions{}
	for _, opt := range opts {
		opt(&options)
	}

	query := goqu.
		From(TablePaymentMethod).
		Select("*").
		Where(goqu.Ex{ColumnPaymentMethodDeletedAt: nil})

	whereExps, orderedExps := qb.buildPaymentMethodExpressions(options)

	query = qb.buildPaymentMethodQuery(
		query,
		options,
		whereExps,
		orderedExps,
	)

	sql, args, err := query.ToSQL()
	if err != nil {
		return nil, errs.New(err)
	}

	var paymentMethods []entity.PaymentMethod
	if err := pgxscan.Select(ctx, qb.db, &paymentMethods, sql, args...); err != nil {
		return nil, errs.New(err)
	}

	return paymentMethods, nil
}

func (qb *QueryBuilder) CountPaymentMethods(
	ctx context.Context,
	opts ...repo.ListPaymentMethodsOption,
) (int64, error) {
	options := repo.ListPaymentMethodsOptions{}
	for _, opt := range opts {
		opt(&options)
	}

	query := goqu.
		From(TablePaymentMethod).
		Select(goqu.COUNT("*")).
		Where(goqu.Ex{ColumnPaymentMethodDeletedAt: nil})

	whereExps, _ := qb.buildPaymentMethodExpressions(options)

	query = qb.buildPaymentMethodQuery(query, options, whereExps, nil)

	sql, args, err := query.ToSQL()
	if err != nil {
		return 0, errs.New(err)
	}

	row := qb.db.QueryRow(ctx, sql, args...)
	var count int64
	if err := row.Scan(&count); err != nil {
		return 0, errs.New(err)
	}

	return count, nil
}

func (qb *QueryBuilder) buildPaymentMethodExpressions(
	options repo.ListPaymentMethodsOptions,
) (whereExps []goqu.Expression, orderedExps []exp.OrderedExpression) {
	options.Search = strings.TrimSpace(options.Search)
	if options.Search != "" {
		searchExp, distanceExp := qb.buildSearch(
			options.Search,
			ColumnPaymentMethodName,
		)
		whereExps = append(whereExps, searchExp)
		orderedExps = append(orderedExps, distanceExp.Asc())
	}

	orderedExps = append(
		orderedExps,
		goqu.I(ColumnPaymentMethodName).Asc(),
	)

	return whereExps, orderedExps
}

func (qb *QueryBuilder) buildPaymentMethodQuery(
	query *goqu.SelectDataset,
	options repo.ListPaymentMethodsOptions,
	whereExps []goqu.Expression,
	orderedExps []exp.OrderedExpression,
) *goqu.SelectDataset {
	if len(whereExps) == 1 {
		query = query.
			Where(whereExps[0])
	} else if len(whereExps) > 0 {
		query = query.
			Where(goqu.And(whereExps...))
	}

	if len(orderedExps) > 0 {
		query = query.
			Order(orderedExps...)
	}

	if options.Limit > 0 {
		query = query.Limit(options.Limit)
	}

	if options.Offset > 0 {
		query = query.Offset(options.Offset)
	}

	return query
}
