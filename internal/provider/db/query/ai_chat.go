package query

import (
	"context"
	"strings"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
)

func (qb *QueryBuilder) ListAIChats(
	ctx context.Context,
	opts ...repo.AIChatOption,
) ([]entity.AIChat, error) {
	options := repo.AIChatOptions{}
	for _, opt := range opts {
		opt(&options)
	}

	query := goqu.
		From(tableAIChat.String()).
		Select(tableAIChat.ColumnAll()).
		Where(goqu.I(tableAIChat.ColumnDeletedAt()).IsNull())

	whereExps, orderedExps := qb.buildAIChatExpressions(options)

	query = qb.buildAIChatQuery(
		query,
		options,
		whereExps,
		orderedExps,
	)

	sql, args, err := query.ToSQL()
	if err != nil {
		return nil, errs.New(err)
	}

	var aiChats []entity.AIChat
	if err := pgxscan.Select(ctx, qb.db, &aiChats, sql, args...); err != nil {
		return nil, errs.New(err)
	}

	return aiChats, nil
}

func (qb *QueryBuilder) CountAIChats(
	ctx context.Context,
	opts ...repo.AIChatOption,
) (int64, error) {
	options := repo.AIChatOptions{}
	for _, opt := range opts {
		opt(&options)
	}

	query := goqu.
		From(tableAIChat.String()).
		Select(goqu.COUNT(tableAIChat.ColumnAll())).
		Where(goqu.I(tableAIChat.ColumnDeletedAt()).IsNull())

	whereExps, _ := qb.buildAIChatExpressions(options)

	query = qb.buildAIChatQuery(query, options, whereExps, nil)

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

func (qb *QueryBuilder) buildAIChatExpressions(
	options repo.AIChatOptions,
) (whereExps []goqu.Expression, orderedExps []exp.OrderedExpression) {
	options.Search = strings.TrimSpace(options.Search)
	if options.Search != "" {
		searchExp, distanceExp := qb.buildSearch(
			options.Search,
			tableAIChat.ColumnTitle(),
		)
		whereExps = append(whereExps, searchExp)
		orderedExps = append(orderedExps, distanceExp.Asc())
	}

	orderedExps = append(
		orderedExps,
		goqu.I(tableAIChat.ColumnUpdatedAt()).Desc(),
	)

	if options.UserID != uuid.Nil {
		whereExps = append(
			whereExps,
			goqu.I(tableAIChat.ColumnUserID()).Eq(options.UserID),
		)
	}

	return whereExps, orderedExps
}

func (qb *QueryBuilder) buildAIChatQuery(
	query *goqu.SelectDataset,
	options repo.AIChatOptions,
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
