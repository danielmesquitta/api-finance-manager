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
	opts ...repo.ListInstitutionsOption,
) ([]entity.Institution, error) {
	options := repo.ListInstitutionsOptions{}
	for _, opt := range opts {
		opt(&options)
	}

	query := goqu.
		From(TableInstitution).
		Select(fmt.Sprintf("%s.*", TableInstitution))

	qb.setInstitutionJoins(query, options)

	whereExps, orderedExps := qb.buildInstitutionExpressions(options)

	qb.setInstitutionExpressions(query, options, whereExps, orderedExps)

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
	opts ...repo.ListInstitutionsOption,
) (int64, error) {
	options := repo.ListInstitutionsOptions{}
	for _, opt := range opts {
		opt(&options)
	}

	query := goqu.
		From(TableInstitution).
		Select(goqu.COUNT("*"))

	qb.setInstitutionJoins(query, options)

	whereExps, _ := qb.buildInstitutionExpressions(options)

	qb.setInstitutionExpressions(query, options, whereExps, nil)

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

func (qb *QueryBuilder) setInstitutionJoins(
	query *goqu.SelectDataset,
	options repo.ListInstitutionsOptions,
) {
	if options.UserID != uuid.Nil {
		query.Join(
			goqu.I(TableAccount),
			goqu.On(
				goqu.I(fmt.Sprintf("%s.%s", TableAccount, ColumnAccountInstitutionID)).
					Eq(goqu.I(fmt.Sprintf("%s.%s", TableInstitution, ColumnInstitutionID))),
			),
		)
	}
}

func (qb *QueryBuilder) buildInstitutionExpressions(
	options repo.ListInstitutionsOptions,
) (whereExps []goqu.Expression, orderedExps []exp.OrderedExpression) {
	options.Search = strings.TrimSpace(options.Search)
	if options.Search != "" {
		searchExp, distanceExp := qb.buildSearch(
			options.Search,
			ColumnInstitutionName,
		)
		whereExps = append(whereExps, searchExp)
		orderedExps = append(orderedExps, distanceExp.Asc())
	}

	if options.UserID != uuid.Nil {
		whereExps = append(
			whereExps,
			goqu.I(fmt.Sprintf("%s.%s", TableAccount, ColumnAccountUserID)).
				Eq(options.UserID),
		)
	}

	orderedExps = append(
		orderedExps,
		goqu.I(ColumnInstitutionName).Asc(),
	)

	return whereExps, orderedExps
}

func (qb *QueryBuilder) setInstitutionExpressions(
	query *goqu.SelectDataset,
	options repo.ListInstitutionsOptions,
	whereExps []goqu.Expression,
	orderedExps []exp.OrderedExpression,
) {
	if len(whereExps) == 1 {
		query.
			Where(whereExps[0])
	} else if len(whereExps) > 0 {
		query.
			Where(goqu.And(whereExps...))
	}

	if len(orderedExps) > 0 {
		query.
			Order(orderedExps...)
	}

	if options.Limit > 0 {
		query.Limit(options.Limit)
	}

	if options.Offset > 0 {
		query.Offset(options.Offset)
	}
}
