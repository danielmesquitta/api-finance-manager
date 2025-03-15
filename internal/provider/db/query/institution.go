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

func (qb *QueryBuilder) ListInstitutions(
	ctx context.Context,
	opts ...repo.InstitutionOption,
) ([]entity.Institution, error) {
	options := repo.InstitutionOptions{}
	for _, opt := range opts {
		opt(&options)
	}

	query := goqu.
		From(schema.Institution.Table()).
		Select(schema.Institution.ColumnAll()).
		Distinct(schema.Institution.ColumnID()).
		Where(goqu.I(schema.Institution.ColumnDeletedAt()).IsNull())

	query = qb.buildInstitutionJoins(query, options)

	orderedExps := []exp.OrderedExpression{
		goqu.I(schema.Institution.ColumnID()).Asc(),
	}

	whereExps, auxOrderedExps := qb.buildInstitutionExpressions(options)

	orderedExps = append(orderedExps, auxOrderedExps...)

	query = qb.buildInstitutionQuery(query, options, whereExps, orderedExps)

	var institutions []entity.Institution
	if err := qb.Scan(ctx, query, &institutions); err != nil {
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
		From(schema.Institution.Table()).
		Select(
			goqu.COUNT(
				goqu.DISTINCT(schema.Institution.ColumnAll()),
			),
		).
		Where(goqu.I(schema.Institution.ColumnDeletedAt()).IsNull())

	query = qb.buildInstitutionJoins(query, options)

	whereExps, _ := qb.buildInstitutionExpressions(options)

	query = qb.buildInstitutionQuery(query, options, whereExps, nil)

	var count int64
	if err := qb.Scan(ctx, query, &count); err != nil {
		return 0, errs.New(err)
	}

	return count, nil
}

func (qb *QueryBuilder) buildInstitutionJoins(
	query *goqu.SelectDataset,
	options repo.InstitutionOptions,
) *goqu.SelectDataset {
	if len(options.UserIDs) > 0 {
		query = query.Join(
			goqu.I(schema.UserInstitution.Table()),
			goqu.On(
				goqu.I(schema.UserInstitution.ColumnInstitutionID()).
					Eq(goqu.I(schema.Institution.ColumnID())),
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
			schema.Institution.ColumnName(),
		)
		whereExps = append(whereExps, searchExp)
		orderedExps = append(orderedExps, distanceExp.Asc())
	}

	if len(options.UserIDs) > 0 {
		whereExps = append(
			whereExps,
			goqu.I(schema.UserInstitution.ColumnUserID()).
				In(options.UserIDs),
		)
	}

	orderedExps = append(
		orderedExps,
		goqu.I(schema.Institution.ColumnName()).Asc(),
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
		query = query.Order(
			orderedExps...,
		)
	}

	if options.Limit > 0 {
		query = query.Limit(options.Limit)
	}

	if options.Offset > 0 {
		query = query.Offset(options.Offset)
	}

	return query
}
