package query

import (
	"math"

	"github.com/danielmesquitta/api-finance-manager/internal/config"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/db"
	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
)

type QueryBuilder struct {
	e  *config.Env
	db *db.DB
}

func NewQueryBuilder(
	e *config.Env,
	db *db.DB,
) *QueryBuilder {
	return &QueryBuilder{
		e:  e,
		db: db,
	}
}

func (qb *QueryBuilder) calculateMaxLevenshteinDistance(search string) int {
	searchLength := float64(len(search))
	maxLevenshteinDistance := int(
		math.Floor(qb.e.MaxLevenshteinDistancePercentage * searchLength),
	)
	return maxLevenshteinDistance
}

func (qb *QueryBuilder) buildSearch(
	search string, column Column,
) (exp.Expression, exp.SQLFunctionExpression) {
	unaccentedColumn := goqu.Func("lower", goqu.Func(
		"unaccent",
		goqu.I(string(column)),
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
