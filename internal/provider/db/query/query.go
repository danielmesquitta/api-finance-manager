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

func (qb *QueryBuilder) buildSearch(
	search, column string,
) (exp.Expression, exp.Orderable) {
	searchQuery := fmt.Sprintf(
		"%s @@ plainto_tsquery('portuguese', unaccent(?))",
		column,
	)

	whereExp := goqu.L(searchQuery, search)

	orderExp := goqu.Func(
		"ts_rank",
		goqu.I(column),
		goqu.L("plainto_tsquery('portuguese', unaccent(?))", search),
	)

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
			return errs.New(err)
		}
		return nil

	case reflect.Struct:
		if err := qb.db.GetContext(ctx, dest, query, args...); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil
			}
			return errs.New(err)
		}
		return nil

	default:
		row := qb.db.QueryRowContext(ctx, query, args...)
		if err := row.Scan(dest); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil
			}
			return errs.New(err)
		}
		return nil
	}
}
