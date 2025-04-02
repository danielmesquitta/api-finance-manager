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
	"github.com/google/uuid"
)

func (qb *QueryBuilder) ListAIChats(
	ctx context.Context,
	opts ...repo.AIChatOptions,
) ([]entity.AIChat, error) {
	options := prepareOptions(opts...)

	query := goqu.
		From(schema.AIChat.String()).
		Select(schema.AIChat.All()).
		Distinct(schema.AIChat.ID()).
		Where(goqu.I(schema.AIChat.DeletedAt()).IsNull())

	joins := qb.buildAIChatJoins(options)

	whereExps, orderedExps := qb.buildAIChatExpressions(options)

	orderedExps = append(
		[]exp.OrderedExpression{goqu.I(schema.AIChat.ID()).Asc()},
		orderedExps...,
	)

	query = qb.buildAIChatQuery(
		query,
		options,
		whereExps,
		joins,
		orderedExps,
	)

	var aiChats []entity.AIChat
	if err := qb.Scan(ctx, query, &aiChats); err != nil {
		return nil, errs.New(err)
	}

	return aiChats, nil
}

func (qb *QueryBuilder) CountAIChats(
	ctx context.Context,
	opts ...repo.AIChatOptions,
) (int64, error) {
	options := prepareOptions(opts...)

	query := goqu.
		From(schema.AIChat.String()).
		Select(
			goqu.COUNT(
				goqu.DISTINCT(schema.AIChat.ID()),
			),
		).
		Where(goqu.I(schema.AIChat.DeletedAt()).IsNull())

	joins := qb.buildAIChatJoins(options)

	whereExps, _ := qb.buildAIChatExpressions(options)

	query = qb.buildAIChatQuery(query, options, whereExps, joins, nil)

	var count int64
	if err := qb.Scan(ctx, query, &count); err != nil {
		return 0, errs.New(err)
	}

	return count, nil
}

func (qb *QueryBuilder) buildAIChatExpressions(
	options repo.AIChatOptions,
) (whereExps []goqu.Expression, orderedExps []exp.OrderedExpression) {
	options.Search = strings.TrimSpace(options.Search)
	if options.Search != "" {
		searchExp, orderExp := qb.buildSearch(
			options.Search,
			schema.AIChat.Title(),
			schema.AIChatMessage.Message(),
			schema.AIChatAnswer.Message(),
		)
		whereExps = append(whereExps, searchExp)
		orderedExps = append(orderedExps, orderExp.Desc())
	}

	orderedExps = append(
		orderedExps,
		goqu.I(schema.AIChat.UpdatedAt()).Desc(),
	)

	if options.UserID != uuid.Nil {
		whereExps = append(
			whereExps,
			goqu.I(schema.AIChat.UserID()).Eq(options.UserID),
		)
	}

	return whereExps, orderedExps
}

func (qb *QueryBuilder) buildAIChatQuery(
	query *goqu.SelectDataset,
	options repo.AIChatOptions,
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

func (qb *QueryBuilder) buildAIChatJoins(
	options repo.AIChatOptions,
) []Join {
	aiChatMessage := Join{
		Table: goqu.I(schema.AIChatMessage.String()),
		Condition: goqu.
			On(
				goqu.I(schema.AIChat.ID()).
					Eq(goqu.I(schema.AIChatMessage.AiChatID())),
			),
	}

	aiChatAnswer := Join{
		Table: goqu.I(schema.AIChatAnswer.String()),
		Condition: goqu.
			On(
				goqu.I(schema.AIChatAnswer.AiChatMessageID()).
					Eq(goqu.I(schema.AIChatMessage.ID())),
			),
	}

	if options.Search != "" {
		return []Join{
			aiChatMessage,
			aiChatAnswer,
		}
	}

	return []Join{}
}
