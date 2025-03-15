package query

import (
	"context"
	"errors"
	"log"
	"math"
	"reflect"

	"github.com/danielmesquitta/api-finance-manager/internal/config"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type QueryBuilder struct {
	e *config.Env
	p *pgxpool.Pool
}

func NewQueryBuilder(
	e *config.Env,
	p *pgxpool.Pool,
) *QueryBuilder {
	return &QueryBuilder{
		e: e,
		p: p,
	}
}

type Join struct {
	Table     exp.Expression
	Condition exp.JoinCondition
}

func (qb *QueryBuilder) calculateMaxLevenshteinDistance(search string) int {
	searchLength := float64(len(search))
	maxLevenshteinDistance := int(
		math.Floor(qb.e.MaxLevenshteinDistancePercentage * searchLength),
	)
	return maxLevenshteinDistance
}

func (qb *QueryBuilder) buildSearch(
	search, column string,
) (exp.Expression, exp.SQLFunctionExpression) {
	unaccentedColumn := goqu.Func("lower", goqu.Func(
		"unaccent",
		goqu.I(column),
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

func (qb *QueryBuilder) Scan(
	ctx context.Context,
	query *goqu.SelectDataset,
	dest any,
) error {
	sql, args, err := query.ToSQL()
	if err != nil {
		return errs.New(err)
	}

	if qb.e.Environment == config.EnvironmentTest {
		log.Printf("Query: %s", sql)
	}

	val := reflect.ValueOf(dest)
	if val.Kind() != reflect.Ptr {
		return errs.New("dest must be a pointer")
	}

	elemKind := val.Elem().Kind()

	switch elemKind {
	case reflect.Slice:
		if err := pgxscan.Select(ctx, qb.p, dest, sql, args...); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil
			}
			return errs.New(err)
		}
		return nil

	case reflect.Struct:
		if err := pgxscan.Get(ctx, qb.p, dest, sql, args...); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil
			}
			return errs.New(err)
		}
		return nil

	default:
		row := qb.p.QueryRow(ctx, sql, args...)
		if err := row.Scan(dest); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil
			}
			return errs.New(err)
		}
		return nil
	}
}
