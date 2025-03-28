package query

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"

	"github.com/danielmesquitta/api-finance-manager/internal/config/env"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/jmoiron/sqlx"
)

type QueryBuilder struct {
	e  *env.Env
	db *sqlx.DB
}

func NewQueryBuilder(
	e *env.Env,
	db *sqlx.DB,
) *QueryBuilder {
	return &QueryBuilder{
		e:  e,
		db: db,
	}
}

type Join struct {
	Table     exp.Expression
	Condition exp.JoinCondition
}

// buildSearch builds a search expression and an orderable expression like:
// WHERE column1 % indexed_unaccent('search') OR column2 % indexed_unaccent('search') ...
// ORDER BY GREATEST(similarity(column1, 'search'), similarity(column2, 'search'), ...)
func (qb *QueryBuilder) buildSearch(
	search string,
	columns ...string,
) (exp.Expression, exp.Orderable) {
	var whereExps []exp.Expression
	var similarityExps []any
	for _, col := range columns {
		whereExps = append(whereExps,
			goqu.L(col+" % indexed_unaccent(?)", search),
		)
		similarityExps = append(
			similarityExps,
			goqu.Func(
				"similarity",
				goqu.I(col),
				goqu.L("indexed_unaccent(?)", search),
			),
		)
	}

	whereExp := goqu.Or(whereExps...)
	orderExp := goqu.Func("GREATEST", similarityExps...)

	return whereExp, orderExp
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

	return qb.scan(ctx, sql, dest, args...)
}

func (qb *QueryBuilder) ScanRaw(
	ctx context.Context,
	query string,
	dest any,
) error {
	return qb.scan(ctx, query, dest)
}

func (qb *QueryBuilder) scan(
	ctx context.Context,
	query string,
	dest any,
	args ...any,
) error {
	val := reflect.ValueOf(dest)
	if val.Kind() != reflect.Ptr {
		return errs.New("dest must be a pointer")
	}

	elemKind := val.Elem().Kind()

	switch elemKind {
	case reflect.Slice:
		if err := qb.db.SelectContext(ctx, dest, query, args...); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil
			}
			return errs.New(
				fmt.Errorf(
					"failed to execute sql query: sql query: %s: %w",
					query,
					err,
				),
			)
		}
		return nil

	case reflect.Struct:
		if err := qb.db.GetContext(ctx, dest, query, args...); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil
			}
			return errs.New(
				fmt.Errorf(
					"failed to execute sql query: sql query: %s: %w",
					query,
					err,
				),
			)
		}
		return nil

	case reflect.Map:
		row := qb.db.QueryRowxContext(ctx, query, args...)
		m := make(map[string]any)
		if err := row.MapScan(m); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil
			}
			return errs.New(
				fmt.Errorf(
					"failed to execute sql query: sql query: %s: %w",
					query,
					err,
				),
			)
		}
		reflect.ValueOf(dest).Elem().Set(reflect.ValueOf(m))
		return nil

	default:
		row := qb.db.QueryRowContext(ctx, query, args...)
		if err := row.Scan(dest); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil
			}
			return errs.New(
				fmt.Errorf(
					"failed to execute sql query: sql query: %s: %w",
					query,
					err,
				),
			)
		}
		return nil
	}
}

func prepareOptions[T any](
	opts ...T,
) T {
	var zero T
	if len(opts) < 1 {
		return zero
	}
	return opts[0]
}
