package query

import (
	"context"
	"fmt"
	"strings"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/georgysavva/scany/v2/pgxscan"
)

func (qb *QueryBuilder) SearchCategories(
	ctx context.Context,
	params repo.SearchCategoriesParams,
) ([]entity.Category, error) {
	query := goqu.From(string(TableCategory))
	search := strings.TrimSpace(params.Search)

	if search != "" {
		whereExp, distanceExp := qb.buildSearchCategoriesWhere(search)
		query = query.
			Select("*").
			Where(whereExp).
			Order(
				distanceExp.Asc(),
				goqu.I(string(ColumnCategoryName)).Asc(),
			).
			Limit(params.Limit).
			Offset(params.Offset)
	} else {
		query = query.
			Select("*").
			Order(goqu.I(string(ColumnCategoryName)).Asc()).
			Limit(params.Limit).
			Offset(params.Offset)
	}

	sql, args, err := query.ToSQL()
	if err != nil {
		return nil, errs.New(err)
	}

	fmt.Println(sql)

	var categories []entity.Category
	if err := pgxscan.Select(ctx, qb.db, &categories, sql, args...); err != nil {
		return nil, errs.New(err)
	}

	return categories, nil
}

func (qb *QueryBuilder) CountSearchCategories(
	ctx context.Context,
	search string,
) (int64, error) {
	query := goqu.From(string(TableCategory))
	search = strings.TrimSpace(search)
	if search != "" {
		where, _ := qb.buildSearchCategoriesWhere(search)
		query = query.
			Select(goqu.COUNT("*")).
			Where(where)
	} else {
		query = query.Select(goqu.COUNT("*"))
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

func (qb *QueryBuilder) buildSearchCategoriesWhere(
	search string,
) (exp.Expression, exp.SQLFunctionExpression) {
	unaccentedColumn := goqu.Func("lower", goqu.Func(
		"unaccent",
		goqu.I(string(ColumnCategoryName)),
	))
	searchPlaceholder := goqu.L("?", search)
	unaccentedSearch := goqu.Func(
		"lower",
		goqu.Func("unaccent", searchPlaceholder),
	)
	distanceExp := goqu.Func(
		"levenshtein",
		unaccentedColumn,
		unaccentedSearch,
	)
	maxLevenshteinDistance := qb.calculateMaxLevenshteinDistance(search)
	likeInput := goqu.Func("concat", "%", unaccentedSearch, "%")
	whereExp := goqu.Or(
		unaccentedColumn.Like(likeInput),
		distanceExp.Lte(maxLevenshteinDistance),
	)
	return whereExp, distanceExp
}
