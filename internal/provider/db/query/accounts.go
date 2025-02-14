package query

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/georgysavva/scany/v2/pgxscan"
)

func (qb *QueryBuilder) ListAccounts(
	ctx context.Context,
	opts ...repo.AccountOption,
) ([]entity.Account, error) {
	options := repo.AccountOptions{}
	for _, opt := range opts {
		opt(&options)
	}

	query := goqu.
		From(tableAccount.String()).
		Select(fmt.Sprintf("%s.*", tableAccount)).
		Where(goqu.I(tableAccount.ColumnDeletedAt()).IsNull())

	joins := qb.buildAccountJoins(options)

	whereExps, orderedExps := qb.buildAccountExpressions(options)

	query = qb.buildAccountsQuery(query, options, whereExps, joins, orderedExps)

	sql, args, err := query.ToSQL()
	if err != nil {
		return nil, errs.New(err)
	}

	var accounts []entity.Account
	if err := pgxscan.Select(ctx, qb.db, &accounts, sql, args...); err != nil {
		return nil, errs.New(err)
	}

	return accounts, nil
}

func (qb *QueryBuilder) ListFullAccounts(
	ctx context.Context,
	opts ...repo.AccountOption,
) ([]entity.FullAccount, error) {
	options := repo.AccountOptions{}
	for _, opt := range opts {
		opt(&options)
	}

	query := goqu.
		From(tableAccount.String()).
		Select(
			fmt.Sprintf("%s.*", tableAccount),
			goqu.I(tableUser.ColumnSynchronizedAt()).As("synchronized_at"),
			goqu.I(tableUser.ColumnOpenFinanceID()).As("open_finance_id"),
		).
		Where(goqu.I(tableAccount.ColumnDeletedAt()).IsNull())

	joins := qb.buildAccountJoins(options, tableUser.String())

	whereExps, orderedExps := qb.buildAccountExpressions(options)

	query = qb.buildAccountsQuery(query, options, whereExps, joins, orderedExps)

	sql, args, err := query.ToSQL()
	if err != nil {
		return nil, errs.New(err)
	}

	var accounts []entity.FullAccount
	if err := pgxscan.Select(ctx, qb.db, &accounts, sql, args...); err != nil {
		return nil, errs.New(err)
	}

	return accounts, nil
}

func (qb *QueryBuilder) CountAccounts(
	ctx context.Context,
	opts ...repo.AccountOption,
) (int64, error) {
	options := repo.AccountOptions{}
	for _, opt := range opts {
		opt(&options)
	}

	query := goqu.
		From(tableAccount.String()).
		Select(goqu.COUNT(fmt.Sprintf("%s.*", tableAccount))).
		Where(goqu.I(tableAccount.ColumnDeletedAt()).IsNull())

	joins := qb.buildAccountJoins(options)

	whereExps, _ := qb.buildAccountExpressions(options)

	query = qb.buildAccountsQuery(query, options, whereExps, joins, nil)

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

func (qb *QueryBuilder) accountShouldJoinUser(
	options repo.AccountOptions,
	defaultTablesToJoin ...string,
) bool {
	if slices.Contains(defaultTablesToJoin, tableUser.String()) {
		return true
	}

	if options.IsSubscriptionActive != nil || len(options.UserTiers) > 0 {
		return true
	}

	return false
}

type BuildAccountJoinsOptions struct {
	ShouldJoinUser bool
}

func (qb *QueryBuilder) buildAccountJoins(
	options repo.AccountOptions,
	defaultTablesToJoin ...string,
) (joins []Join) {
	shouldJoinUser := qb.accountShouldJoinUser(
		options,
		defaultTablesToJoin...)

	if shouldJoinUser {
		join := Join{
			Table: goqu.I(tableUser.String()),
			Condition: goqu.
				On(
					goqu.I(tableAccount.ColumnUserID()).
						Eq(goqu.I(tableUser.ColumnID())),
				),
		}
		joins = append(joins, join)
	}

	return joins
}

func (qb *QueryBuilder) buildAccountExpressions(
	options repo.AccountOptions,
) (whereExps []goqu.Expression, orderedExps []exp.OrderedExpression) {
	options.Search = strings.TrimSpace(options.Search)
	if options.Search != "" {
		searchExp, distanceExp := qb.buildSearch(
			options.Search,
			tableAccount.ColumnName(),
		)
		whereExps = append(whereExps, searchExp)
		orderedExps = append(orderedExps, distanceExp.Asc())
	}

	if len(options.UserIDs) > 0 {
		exp := goqu.I(tableAccount.ColumnUserID()).In(options.UserIDs)
		whereExps = append(whereExps, exp)
	}

	if len(options.ExternalIDs) > 0 {
		exp := goqu.I(tableAccount.ColumnExternalID()).In(options.ExternalIDs)
		whereExps = append(whereExps, exp)
	}

	if len(options.UserTiers) > 0 {
		exp := goqu.I(tableUser.ColumnTier()).In(options.UserTiers)
		whereExps = append(whereExps, exp)
	}

	if options.IsSubscriptionActive != nil {
		var exp goqu.Expression
		ident := goqu.I(tableUser.ColumnSubscriptionExpiresAt())
		if *options.IsSubscriptionActive {
			exp = ident.Gte(time.Now())
		} else {
			exp = ident.Lt(time.Now())
		}
		whereExps = append(whereExps, exp)
	}

	if shouldJoinUser := qb.accountShouldJoinUser(options); shouldJoinUser {
		exp := goqu.I(tableUser.ColumnDeletedAt()).IsNull()
		whereExps = append(whereExps, exp)
	}

	orderedExps = append(
		orderedExps,
		goqu.I(tableAccount.ColumnName()).Asc(),
	)

	return whereExps, orderedExps
}

func (qb *QueryBuilder) buildAccountsQuery(
	query *goqu.SelectDataset,
	options repo.AccountOptions,
	whereExps []goqu.Expression,
	joins []Join,
	orderedExps []exp.OrderedExpression,
) *goqu.SelectDataset {
	for _, join := range joins {
		query = query.LeftJoin(join.Table, join.Condition)
	}

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
