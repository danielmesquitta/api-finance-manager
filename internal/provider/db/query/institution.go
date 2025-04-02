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
	opts ...repo.InstitutionOptions,
) ([]entity.Institution, error) {
	options := prepareOptions(opts...)

	query := goqu.
		From(schema.Institution.String()).
		Select(schema.Institution.All()).
		Distinct(schema.Institution.ID()).
		Where(goqu.I(schema.Institution.DeletedAt()).IsNull())

	query = qb.buildInstitutionJoins(query, options)

	orderedExps := []exp.OrderedExpression{
		goqu.I(schema.Institution.ID()).Asc(),
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
	opts ...repo.InstitutionOptions,
) (int64, error) {
	options := prepareOptions(opts...)

	query := goqu.
		From(schema.Institution.String()).
		Select(
			goqu.COUNT(
				goqu.DISTINCT(schema.Institution.All()),
			),
		).
		Where(goqu.I(schema.Institution.DeletedAt()).IsNull())

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
			goqu.I(schema.UserInstitution.String()),
			goqu.On(
				goqu.I(schema.UserInstitution.InstitutionID()).
					Eq(goqu.I(schema.Institution.ID())),
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
		searchExp, orderExp := qb.buildSearch(
			options.Search,
			schema.Institution.Name(),
		)
		whereExps = append(whereExps, searchExp)
		orderedExps = append(orderedExps, orderExp.Desc())
	}

	if len(options.UserIDs) > 0 {
		whereExps = append(
			whereExps,
			goqu.I(schema.UserInstitution.UserID()).
				In(options.UserIDs),
		)
	}

	orderedExps = append(
		orderedExps,
		goqu.I(schema.Institution.Name()).Asc(),
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
