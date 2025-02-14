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
	"github.com/google/uuid"
)

func (qb *QueryBuilder) ListInstitutions(
	ctx context.Context,
	opts ...repo.InstitutionOption,
) ([]entity.Institution, error) {
	options := repo.InstitutionOptions{}
	for _, opt := range opts {
		opt(&options)
	}

	query := goqu.
		From(tableInstitution.String()).
		Select(fmt.Sprintf("%s.*", tableInstitution)).
		Where(goqu.I(tableInstitution.ColumnDeletedAt()).IsNull())

	qb.buildInstitutionJoins(query, options)

	whereExps, orderedExps := qb.buildInstitutionExpressions(options)

	query = qb.buildInstitutionQuery(query, options, whereExps, orderedExps)

	sql, args, err := query.ToSQL()
	if err != nil {
		return nil, errs.New(err)
	}

	var institutions []entity.Institution
	if err := pgxscan.Select(ctx, qb.db, &institutions, sql, args...); err != nil {
		return nil, errs.New(err)
	}

	return institutions, nil
}

func (qb *QueryBuilder) CountInstitutions(
	ctx context.Context,
	opts ...repo.InstitutionOption,
) (int64, error) {
	options := repo.InstitutionOptions{}
	for _, opt := range opts {
		opt(&options)
	}

	query := goqu.
		From(tableInstitution.String()).
		Select(goqu.COUNT("*")).
		Where(goqu.I(tableInstitution.ColumnDeletedAt()).IsNull())

	query = qb.buildInstitutionJoins(query, options)

	whereExps, _ := qb.buildInstitutionExpressions(options)

	query = qb.buildInstitutionQuery(query, options, whereExps, nil)

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

func (qb *QueryBuilder) buildInstitutionJoins(
	query *goqu.SelectDataset,
	options repo.InstitutionOptions,
) *goqu.SelectDataset {
	if options.UserID != uuid.Nil {
		query = query.Join(
			goqu.I(tableAccount.String()),
			goqu.On(
				goqu.I(tableAccount.ColumnInstitutionID()).
					Eq(goqu.I(tableInstitution.ColumnID())),
			),
		)
	}
	return query
}

func (qb *QueryBuilder) buildInstitutionExpressions(
	options repo.InstitutionOptions,
) (whereExps []goqu.Expression, orderedExps []exp.OrderedExpression) {
	options.Search = strings.TrimSpace(options.Search)
	if options.Search != "" {
		searchExp, distanceExp := qb.buildSearch(
			options.Search,
			tableInstitution.ColumnName(),
		)
		whereExps = append(whereExps, searchExp)
		orderedExps = append(orderedExps, distanceExp.Asc())
	}

	if options.UserID != uuid.Nil {
		whereExps = append(
			whereExps,
			goqu.I(tableAccount.ColumnUserID()).
				Eq(options.UserID),
		)
	}

	orderedExps = append(
		orderedExps,
		goqu.I(tableInstitution.ColumnName()).Asc(),
	)

	return whereExps, orderedExps
}

func (qb *QueryBuilder) buildInstitutionQuery(
	query *goqu.SelectDataset,
	options repo.InstitutionOptions,
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
