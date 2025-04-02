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

func (qb *QueryBuilder) ListTransactionCategories(
	ctx context.Context,
	opts ...repo.TransactionCategoryOptions,
) ([]entity.TransactionCategory, error) {
	options := prepareOptions(opts...)

	query := goqu.
		From(schema.TransactionCategory.String()).
		Select(schema.TransactionCategory.All()).
		Where(goqu.I(schema.TransactionCategory.DeletedAt()).IsNull())

	whereExps, orderedExps := qb.buildCategoryExpressions(options)

	query = qb.buildCategoryQuery(
		query,
		options,
		whereExps,
		orderedExps,
	)

	var categories []entity.TransactionCategory
	if err := qb.Scan(ctx, query, &categories); err != nil {
		return nil, errs.New(err)
	}

	return categories, nil
}

func (qb *QueryBuilder) CountTransactionCategories(
	ctx context.Context,
	opts ...repo.TransactionCategoryOptions,
) (int64, error) {
	options := prepareOptions(opts...)

	query := goqu.
		From(schema.TransactionCategory.String()).
		Select(goqu.COUNT(schema.TransactionCategory.All())).
		Where(goqu.I(schema.TransactionCategory.DeletedAt()).IsNull())

	whereExps, _ := qb.buildCategoryExpressions(options)

	query = qb.buildCategoryQuery(query, options, whereExps, nil)

	var count int64
	if err := qb.Scan(ctx, query, &count); err != nil {
		return 0, errs.New(err)
	}

	return count, nil
}

func (qb *QueryBuilder) buildCategoryExpressions(
	options repo.TransactionCategoryOptions,
) (whereExps []goqu.Expression, orderedExps []exp.OrderedExpression) {
	options.Search = strings.TrimSpace(options.Search)
	if options.Search != "" {
		searchExp, orderExp := qb.buildSearch(
			options.Search,
			schema.TransactionCategory.Name(),
		)
		whereExps = append(whereExps, searchExp)
		orderedExps = append(orderedExps, orderExp.Desc())
	}

	if len(options.IDs) > 0 {
		exp := goqu.I(schema.TransactionCategory.ID()).
			In(options.IDs)
		whereExps = append(whereExps, exp)
	}

	orderedExps = append(
		orderedExps,
		goqu.I(schema.TransactionCategory.Name()).Asc(),
	)

	return whereExps, orderedExps
}

func (qb *QueryBuilder) buildCategoryQuery(
	query *goqu.SelectDataset,
	options repo.TransactionCategoryOptions,
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
