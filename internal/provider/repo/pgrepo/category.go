package pgrepo

import (
	"context"
	"fmt"
	"strings"

	"github.com/danielmesquitta/api-finance-manager/internal/config"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/db"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/db/sqlc"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jinzhu/copier"
)

type CategoryPgRepo struct {
	e *config.Env
	q *db.Queries
}

func NewCategoryPgRepo(
	e *config.Env,
	q *db.Queries,
) *CategoryPgRepo {
	return &CategoryPgRepo{
		e: e,
		q: q,
	}
}

func (r *CategoryPgRepo) ListCategories(
	ctx context.Context,
) ([]entity.Category, error) {
	categories, err := r.q.ListCategories(ctx)
	if err != nil {
		return nil, errs.New(err)
	}

	results := []entity.Category{}
	if err := copier.Copy(&results, categories); err != nil {
		return nil, errs.New(err)
	}

	return results, nil
}

func (r *CategoryPgRepo) CreateCategories(
	ctx context.Context,
	params []repo.CreateCategoriesParams,
) error {
	dbParams := make([]sqlc.CreateCategoriesParams, len(params))
	if err := copier.Copy(&dbParams, params); err != nil {
		return errs.New(err)
	}

	tx := r.q.UseTx(ctx)
	_, err := tx.CreateCategories(ctx, dbParams)
	if err != nil {
		return errs.New(err)
	}

	return nil
}

func (r *CategoryPgRepo) SearchCategories(
	ctx context.Context,
	params repo.SearchCategoriesParams,
) ([]entity.Category, error) {
	query := goqu.From(string(TableCategory))
	search := strings.TrimSpace(params.Search)

	if search != "" {
		whereExp, distanceExp := r.buildSearchCategoriesWhere(search)
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
	if err := pgxscan.Select(ctx, r.q, &categories, sql, args...); err != nil {
		return nil, errs.New(err)
	}

	return categories, nil
}

func (r *CategoryPgRepo) CountSearchCategories(
	ctx context.Context,
	search string,
) (int64, error) {
	query := goqu.From(string(TableCategory))
	search = strings.TrimSpace(search)
	if search != "" {
		where, _ := r.buildSearchCategoriesWhere(search)
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

	row := r.q.QueryRow(ctx, sql, args...)
	var count int64
	if err := row.Scan(&count); err != nil {
		return 0, errs.New(err)
	}

	return count, nil
}

func (r *CategoryPgRepo) buildSearchCategoriesWhere(
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
	maxLevenshteinDistance := calculateMaxLevenshteinDistance(
		search,
		r.e.MaxLevenshteinDistancePercentage,
	)
	likeInput := goqu.Func("concat", "%", unaccentedSearch, "%")
	whereExp := goqu.Or(
		unaccentedColumn.Like(likeInput),
		distanceExp.Lte(maxLevenshteinDistance),
	)
	return whereExp, distanceExp
}

var _ repo.CategoryRepo = &CategoryPgRepo{}
