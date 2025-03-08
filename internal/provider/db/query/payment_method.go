package query

import (
	"context"
	"strings"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/db"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/georgysavva/scany/v2/pgxscan"
)

func (qb *QueryBuilder) ListPaymentMethods(
	ctx context.Context,
	opts ...repo.PaymentMethodOption,
) ([]entity.PaymentMethod, error) {
	options := repo.PaymentMethodOptions{}
	for _, opt := range opts {
		opt(&options)
	}

	query := goqu.
		From(db.TablePaymentMethod.String()).
		Select(db.TablePaymentMethod.ColumnAll()).
		Where(goqu.I(db.TablePaymentMethod.ColumnDeletedAt()).IsNull())

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
	opts ...repo.PaymentMethodOption,
) (int64, error) {
	options := repo.PaymentMethodOptions{}
	for _, opt := range opts {
		opt(&options)
	}

	query := goqu.
		From(db.TablePaymentMethod.String()).
		Select(goqu.COUNT(db.TablePaymentMethod.ColumnAll())).
		Where(goqu.I(db.TablePaymentMethod.ColumnDeletedAt()).IsNull())

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
	options repo.PaymentMethodOptions,
) (whereExps []goqu.Expression, orderedExps []exp.OrderedExpression) {
	options.Search = strings.TrimSpace(options.Search)
	if options.Search != "" {
		searchExp, distanceExp := qb.buildSearch(
			options.Search,
			db.TablePaymentMethod.ColumnName(),
		)
		whereExps = append(whereExps, searchExp)
		orderedExps = append(orderedExps, distanceExp.Asc())
	}

	orderedExps = append(
		orderedExps,
		goqu.I(db.TablePaymentMethod.ColumnName()).Asc(),
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
