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
	opts ...repo.TransactionCategoryOption,
) ([]entity.TransactionCategory, error) {
	options := repo.TransactionCategoryOptions{}
	for _, opt := range opts {
		opt(&options)
	}

	query := goqu.
		From(schema.TransactionCategory.Table()).
		Select(schema.TransactionCategory.ColumnAll()).
		Where(goqu.I(schema.TransactionCategory.ColumnDeletedAt()).IsNull())

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
	opts ...repo.TransactionCategoryOption,
) (int64, error) {
	options := repo.TransactionCategoryOptions{}
	for _, opt := range opts {
		opt(&options)
	}

	query := goqu.
		From(schema.TransactionCategory.Table()).
		Select(goqu.COUNT(schema.TransactionCategory.ColumnAll())).
		Where(goqu.I(schema.TransactionCategory.ColumnDeletedAt()).IsNull())

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
			schema.TransactionCategory.ColumnName(),
		)
		whereExps = append(whereExps, searchExp)
		orderedExps = append(orderedExps, orderExp.Desc())
	}

	if len(options.IDs) > 0 {
		exp := goqu.I(schema.TransactionCategory.ColumnID()).
			In(options.IDs)
		whereExps = append(whereExps, exp)
	}

	orderedExps = append(
		orderedExps,
		goqu.I(schema.TransactionCategory.ColumnName()).Asc(),
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
