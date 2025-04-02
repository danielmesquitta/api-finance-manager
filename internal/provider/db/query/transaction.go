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

func (qb *QueryBuilder) ListTransactions(
	ctx context.Context,
	userID uuid.UUID,
	opts ...repo.TransactionOptions,
) ([]entity.Transaction, error) {
	options := prepareOptions(opts...)

	query := goqu.
		From(schema.Transaction.String()).
		Select(schema.Transaction.All()).
		Where(goqu.I(schema.Transaction.DeletedAt()).IsNull())

	joins := qb.buildTransactionJoins(options)

	whereExps, orderedExps := qb.buildTransactionExpressions(userID, options)

	query = qb.buildTransactionsQuery(
		query,
		options,
		whereExps,
		joins,
		orderedExps,
	)

	var transactions []entity.Transaction
	if err := qb.Scan(ctx, query, &transactions); err != nil {
		return nil, errs.New(err)
	}

	return transactions, nil
}

func (qb *QueryBuilder) ListFullTransactions(
	ctx context.Context,
	userID uuid.UUID,
	opts ...repo.TransactionOptions,
) ([]entity.FullTransaction, error) {
	options := prepareOptions(opts...)

	query := goqu.
		From(schema.Transaction.String()).
		Select(
			schema.Transaction.All(),
			goqu.I(schema.TransactionCategory.Name()).
				As("category_name"),
			goqu.I(schema.Institution.Name()).As("institution_name"),
			goqu.I(schema.Institution.Logo()).As("institution_logo"),
			goqu.I(schema.PaymentMethod.Name()).
				As("payment_method_name"),
		).
		Where(goqu.I(schema.Transaction.DeletedAt()).IsNull())

	joins := qb.buildTransactionJoins(options, true)

	whereExps, orderedExps := qb.buildTransactionExpressions(userID, options)

	query = qb.buildTransactionsQuery(
		query,
		options,
		whereExps,
		joins,
		orderedExps,
	)

	var transactions []entity.FullTransaction
	if err := qb.Scan(ctx, query, &transactions); err != nil {
		return nil, errs.New(err)
	}

	return transactions, nil
}

func (qb *QueryBuilder) CountTransactions(
	ctx context.Context,
	userID uuid.UUID,
	opts ...repo.TransactionOptions,
) (int64, error) {
	options := prepareOptions(opts...)

	query := goqu.
		From(schema.Transaction.String()).
		Select(goqu.COUNT(schema.Transaction.All())).
		Where(goqu.I(schema.Transaction.DeletedAt()).IsNull())

	joins := qb.buildTransactionJoins(options)

	whereExps, _ := qb.buildTransactionExpressions(userID, options)

	query = qb.buildTransactionsQuery(query, options, whereExps, joins, nil)

	var count int64
	if err := qb.Scan(ctx, query, &count); err != nil {
		return 0, errs.New(err)
	}

	return count, nil
}

func (qb *QueryBuilder) SumTransactions(
	ctx context.Context,
	userID uuid.UUID,
	opts ...repo.TransactionOptions,
) (int64, error) {
	options := prepareOptions(opts...)

	query := goqu.
		From(schema.Transaction.String()).
		Select(goqu.SUM(schema.Transaction.Amount())).
		Where(goqu.I(schema.Transaction.DeletedAt()).IsNull())

	joins := qb.buildTransactionJoins(options)

	whereExps, _ := qb.buildTransactionExpressions(userID, options)

	query = qb.buildTransactionsQuery(query, options, whereExps, joins, nil)

	var count int64
	if err := qb.Scan(ctx, query, &count); err != nil {
		return 0, errs.New(err)
	}

	return count, nil
}

func (qb *QueryBuilder) SumTransactionsByCategory(
	ctx context.Context,
	userID uuid.UUID,
	opts ...repo.TransactionOptions,
) (map[uuid.UUID]int64, error) {
	options := prepareOptions(opts...)

	query := goqu.
		From(schema.Transaction.String()).
		Select(
			goqu.I(schema.Transaction.CategoryID()),
			goqu.SUM(schema.Transaction.Amount()).As("sum"),
		).
		Where(goqu.I(schema.Transaction.DeletedAt()).IsNull()).
		GroupBy(goqu.I(schema.Transaction.CategoryID()))

	joins := qb.buildTransactionJoins(options)

	whereExps, _ := qb.buildTransactionExpressions(userID, options)

	query = qb.buildTransactionsQuery(query, options, whereExps, joins, nil)

	rows := []struct {
		CategoryID uuid.UUID `db:"category_id"`
		Sum        int64     `db:"sum"`
	}{}
	if err := qb.Scan(ctx, query, &rows); err != nil {
		return nil, errs.New(err)
	}

	out := map[uuid.UUID]int64{}
	for _, row := range rows {
		out[row.CategoryID] = row.Sum
	}

	return out, nil
}

func (qb *QueryBuilder) buildTransactionExpressions(
	userID uuid.UUID,
	options repo.TransactionOptions,
) (whereExps []goqu.Expression, orderedExps []exp.OrderedExpression) {
	whereExps = append(
		whereExps,
		goqu.I(schema.Transaction.UserID()).Eq(userID),
	)

	options.Search = strings.TrimSpace(options.Search)
	if options.Search != "" {
		searchExp, orderExp := qb.buildSearch(
			options.Search,
			schema.Transaction.Name(),
			schema.TransactionCategory.Name(),
			schema.Institution.Name(),
			schema.PaymentMethod.Name(),
		)
		whereExps = append(whereExps, searchExp)
		orderedExps = append(orderedExps, orderExp.Desc())
	}

	if len(options.CategoryIDs) > 0 {
		whereExps = append(
			whereExps,
			goqu.I(schema.Transaction.CategoryID()).
				In(options.CategoryIDs),
		)
	}

	if len(options.InstitutionIDs) > 0 {
		whereExps = append(
			whereExps,
			goqu.I(schema.Transaction.InstitutionID()).
				In(options.InstitutionIDs),
		)
	}

	if len(options.PaymentMethodIDs) > 0 {
		whereExps = append(
			whereExps,
			goqu.I(schema.Transaction.PaymentMethodID()).
				In(options.PaymentMethodIDs),
		)
	}

	if !options.StartDate.IsZero() {
		whereExps = append(
			whereExps,
			goqu.I(schema.Transaction.Date()).Gte(options.StartDate),
		)
	}

	if !options.EndDate.IsZero() {
		whereExps = append(
			whereExps,
			goqu.I(schema.Transaction.Date()).Lte(options.EndDate),
		)
	}

	if options.IsExpense {
		whereExps = append(
			whereExps,
			goqu.I(schema.Transaction.Amount()).Lt(0),
		)
	}

	if options.IsIncome {
		whereExps = append(
			whereExps,
			goqu.I(schema.Transaction.Amount()).Gt(0),
		)
	}

	if options.IsIgnored != nil {
		whereExps = append(
			whereExps,
			goqu.I(schema.Transaction.IsIgnored()).
				Eq(*options.IsIgnored),
		)
	}

	orderedExps = append(
		orderedExps,
		goqu.I(schema.Transaction.Name()).Asc(),
	)

	return whereExps, orderedExps
}

func (qb *QueryBuilder) buildTransactionJoins(
	options repo.TransactionOptions,
	shouldJoinAll ...bool,
) []Join {
	institutionJoin := Join{
		Table: goqu.I(schema.Institution.String()),
		Condition: goqu.
			On(
				goqu.I(schema.Transaction.InstitutionID()).
					Eq(goqu.I(schema.Institution.ID())),
			),
	}

	transactionCategoryJoin := Join{
		Table: goqu.I(schema.TransactionCategory.String()),
		Condition: goqu.
			On(
				goqu.I(schema.Transaction.CategoryID()).
					Eq(goqu.I(schema.TransactionCategory.ID())),
			),
	}

	paymentMethodJoin := Join{
		Table: goqu.I(schema.PaymentMethod.String()),
		Condition: goqu.
			On(
				goqu.I(schema.Transaction.PaymentMethodID()).
					Eq(goqu.I(schema.PaymentMethod.ID())),
			),
	}

	if (len(shouldJoinAll) > 0 && shouldJoinAll[0]) || options.Search != "" {
		return []Join{
			institutionJoin,
			transactionCategoryJoin,
			paymentMethodJoin,
		}
	}

	joins := []Join{}
	if len(options.CategoryIDs) > 0 {
		joins = append(joins, transactionCategoryJoin)
	}

	if len(options.InstitutionIDs) > 0 {
		joins = append(joins, institutionJoin)
	}

	if len(options.PaymentMethodIDs) > 0 {
		joins = append(joins, paymentMethodJoin)
	}

	return joins
}

func (qb *QueryBuilder) buildTransactionsQuery(
	query *goqu.SelectDataset,
	options repo.TransactionOptions,
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
