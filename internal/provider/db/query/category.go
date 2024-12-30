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

func (qb *QueryBuilder) ListCategories(
	ctx context.Context,
	opts ...repo.ListCategoriesOption,
) ([]entity.Category, error) {
	options := repo.ListCategoriesOptions{}
	for _, opt := range opts {
		opt(&options)
	}

	query := goqu.
		From(string(TableCategory)).
		Select("*")

	var whereExps []goqu.Expression
	var orderedExps []exp.OrderedExpression

	options.Search = strings.TrimSpace(options.Search)
	if options.Search != "" {
		searchExp, distanceExp := qb.buildSearch(
			options.Search,
			ColumnCategoryName,
		)
		whereExps = append(whereExps, searchExp)
		orderedExps = append(orderedExps, distanceExp.Asc())
	}

	orderedExps = append(
		orderedExps,
		goqu.I(string(ColumnCategoryName)).Asc(),
	)

	if len(whereExps) == 1 {
		query.
			Where(whereExps[0])
	} else if len(whereExps) > 0 {
		query.
			Where(goqu.And(whereExps...))
	}

	if len(orderedExps) > 0 {
		query.Order(orderedExps...)
	}

	if options.Limit > 0 {
		query.Limit(options.Limit)
	}

	if options.Offset > 0 {
		query.Offset(options.Offset)
	}

	sql, args, err := query.ToSQL()
	if err != nil {
		return nil, errs.New(err)
	}

	var categories []entity.Category
	if err := pgxscan.Select(ctx, qb.db, &categories, sql, args...); err != nil {
		return nil, errs.New(err)
	}

	return categories, nil
}

func (qb *QueryBuilder) CountCategories(
	ctx context.Context,
	opts ...repo.ListCategoriesOption,
) (int64, error) {
	options := repo.ListCategoriesOptions{}
	for _, opt := range opts {
		opt(&options)
	}

	query := goqu.
		From(string(TableCategory)).
		Select(goqu.COUNT("*"))

	var whereExps []goqu.Expression

	options.Search = strings.TrimSpace(options.Search)
	if options.Search != "" {
		searchExp, _ := qb.buildSearch(options.Search, ColumnCategoryName)
		whereExps = append(whereExps, searchExp)
	}

	if len(whereExps) == 1 {
		query.
			Where(whereExps[0])
	} else if len(whereExps) > 0 {
		query.
			Where(goqu.And(whereExps...))
	}

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
