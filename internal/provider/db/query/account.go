package query

import (
	"context"
	"time"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/db/schema"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
)

func (qb *QueryBuilder) ListAccounts(
	ctx context.Context,
	opts ...repo.AccountOptions,
) ([]entity.Account, error) {
	options := prepareOptions(opts...)

	query := goqu.
		From(schema.Account.String()).
		Select(schema.Account.All()).
		Where(goqu.I(schema.Account.DeletedAt()).IsNull())

	joins := qb.buildAccountJoins(options)

	whereExps, orderedExps := qb.buildAccountExpressions(options)

	query = qb.buildAccountsQuery(query, options, whereExps, joins, orderedExps)

	var accounts []entity.Account
	if err := qb.Scan(ctx, query, &accounts); err != nil {
		return nil, errs.New(err)
	}

	return accounts, nil
}

func (qb *QueryBuilder) ListFullAccounts(
	ctx context.Context,
	opts ...repo.AccountOptions,
) ([]entity.FullAccount, error) {
	options := prepareOptions(opts...)

	query := goqu.
		From(schema.Account.String()).
		Select(
			schema.Account.All(),
			goqu.I(schema.UserInstitution.UserID()),
			goqu.I(schema.UserInstitution.InstitutionID()),
			goqu.I(schema.UserInstitution.ExternalID()).
				As("user_institution_external_id"),
			goqu.I(schema.User.SynchronizedAt()),
		).
		Where(goqu.I(schema.Account.DeletedAt()).IsNull())

	joins := qb.buildAccountJoins(options, true)

	whereExps, orderedExps := qb.buildAccountExpressions(options)

	query = qb.buildAccountsQuery(query, options, whereExps, joins, orderedExps)

	var accounts []entity.FullAccount
	if err := qb.Scan(ctx, query, &accounts); err != nil {
		return nil, errs.New(err)
	}

	return accounts, nil
}

func (qb *QueryBuilder) CountAccounts(
	ctx context.Context,
	opts ...repo.AccountOptions,
) (int64, error) {
	options := prepareOptions(opts...)

	query := goqu.
		From(schema.Account.String()).
		Select(goqu.COUNT(schema.Account.All())).
		Where(goqu.I(schema.Account.DeletedAt()).IsNull())

	joins := qb.buildAccountJoins(options)

	whereExps, _ := qb.buildAccountExpressions(options)

	query = qb.buildAccountsQuery(query, options, whereExps, joins, nil)

	var count int64
	if err := qb.Scan(ctx, query, &count); err != nil {
		return 0, errs.New(err)
	}

	return count, nil
}

func (qb *QueryBuilder) buildAccountJoins(
	options repo.AccountOptions,
	shouldJoinAll ...bool,
) []Join {
	userInstitutionJoin := Join{
		Table: goqu.I(schema.UserInstitution.String()),
		Condition: goqu.
			On(
				goqu.I(schema.Account.UserInstitutionID()).
					Eq(goqu.I(schema.UserInstitution.ID())),
			),
	}

	userJoin := Join{
		Table: goqu.I(schema.User.String()),
		Condition: goqu.
			On(
				goqu.I(schema.UserInstitution.UserID()).
					Eq(goqu.I(schema.User.ID())),
			),
	}

	if len(shouldJoinAll) > 0 && shouldJoinAll[0] {
		return []Join{userInstitutionJoin, userJoin}
	}

	if len(options.UserIDs) > 0 || len(options.UserTiers) > 0 ||
		options.IsSubscriptionActive != nil {
		return []Join{userJoin}
	}

	return []Join{}
}

func (qb *QueryBuilder) buildAccountExpressions(
	options repo.AccountOptions,
) (whereExps []goqu.Expression, orderedExps []exp.OrderedExpression) {
	if len(options.UserIDs) > 0 {
		exp := goqu.I(schema.User.ID()).In(options.UserIDs)
		whereExps = append(whereExps, exp)
	}

	if len(options.ExternalIDs) > 0 {
		exp := goqu.I(schema.Account.ExternalID()).
			In(options.ExternalIDs)
		whereExps = append(whereExps, exp)
	}

	if len(options.UserTiers) > 0 {
		exp := goqu.I(schema.User.Tier()).In(options.UserTiers)
		whereExps = append(whereExps, exp)
	}

	if len(options.Types) > 0 {
		exp := goqu.I(schema.Account.Type()).In(options.Types)
		whereExps = append(whereExps, exp)
	}

	if options.IsSubscriptionActive != nil {
		var exp goqu.Expression
		ident := goqu.I(schema.User.SubscriptionExpiresAt())
		if *options.IsSubscriptionActive {
			exp = ident.Gte(time.Now())
		} else {
			exp = ident.Lt(time.Now())
		}
		whereExps = append(whereExps, exp)
	}

	orderedExps = append(
		orderedExps,
		goqu.I(schema.Account.Name()).Asc(),
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
