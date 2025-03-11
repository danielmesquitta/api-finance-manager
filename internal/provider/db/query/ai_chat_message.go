package query

import (
	"context"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/db/schema"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/google/uuid"
)

func (qb *QueryBuilder) ListAIChatMessages(
	ctx context.Context,
	opts ...repo.AIChatMessageOption,
) ([]entity.AIChatMessage, error) {
	options := repo.AIChatMessageOptions{}
	for _, opt := range opts {
		opt(&options)
	}

	query := goqu.
		From(schema.AIChatMessage.Table()).
		Select(schema.AIChatMessage.ColumnAll()).
		Where(goqu.I(schema.AIChatMessage.ColumnDeletedAt()).IsNull())

	whereExps, orderedExps := qb.buildAIChatMessageExpressions(options)

	query = qb.buildAIChatMessageQuery(
		query,
		options,
		whereExps,
		orderedExps,
	)

	var aiChatMessages []entity.AIChatMessage
	if err := qb.Scan(ctx, query, &aiChatMessages); err != nil {
		return nil, errs.New(err)
	}

	return aiChatMessages, nil
}

func (qb *QueryBuilder) CountAIChatMessages(
	ctx context.Context,
	opts ...repo.AIChatMessageOption,
) (int64, error) {
	options := repo.AIChatMessageOptions{}
	for _, opt := range opts {
		opt(&options)
	}

	query := goqu.
		From(schema.AIChatMessage.Table()).
		Select(goqu.COUNT(schema.AIChatMessage.ColumnAll())).
		Where(goqu.I(schema.AIChatMessage.ColumnDeletedAt()).IsNull())

	whereExps, _ := qb.buildAIChatMessageExpressions(options)

	query = qb.buildAIChatMessageQuery(query, options, whereExps, nil)

	var count int64
	if err := qb.Scan(ctx, query, &count); err != nil {
		return 0, errs.New(err)
	}

	return count, nil
}

func (qb *QueryBuilder) buildAIChatMessageExpressions(
	options repo.AIChatMessageOptions,
) (whereExps []goqu.Expression, orderedExps []exp.OrderedExpression) {
	orderedExps = append(
		orderedExps,
		goqu.I(schema.AIChatMessage.ColumnUpdatedAt()).Desc(),
	)

	if options.AIChatID != uuid.Nil {
		whereExps = append(
			whereExps,
			goqu.I(schema.AIChatMessage.ColumnAiChatID()).Eq(options.AIChatID),
		)
	}

	return whereExps, orderedExps
}

func (qb *QueryBuilder) buildAIChatMessageQuery(
	query *goqu.SelectDataset,
	options repo.AIChatMessageOptions,
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
