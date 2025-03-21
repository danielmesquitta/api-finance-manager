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
	opts ...repo.TransactionOption,
) ([]entity.Transaction, error) {
	options := repo.TransactionOptions{}
	for _, opt := range opts {
		opt(&options)
	}

	query := goqu.
		From(schema.Transaction.Table()).
		Select(schema.Transaction.ColumnAll()).
		Where(goqu.I(schema.Transaction.ColumnDeletedAt()).IsNull())

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
	opts ...repo.TransactionOption,
) ([]entity.FullTransaction, error) {
	options := repo.TransactionOptions{}
	for _, opt := range opts {
		opt(&options)
	}

	query := goqu.
		From(schema.Transaction.Table()).
		Select(
			schema.Transaction.ColumnAll(),
			goqu.I(schema.TransactionCategory.ColumnName()).
				As("category_name"),
			goqu.I(schema.Institution.ColumnName()).As("institution_name"),
			goqu.I(schema.Institution.ColumnLogo()).As("institution_logo"),
			goqu.I(schema.PaymentMethod.ColumnName()).
				As("payment_method_name"),
		).
		Where(goqu.I(schema.Transaction.ColumnDeletedAt()).IsNull())

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
	opts ...repo.TransactionOption,
) (int64, error) {
	options := repo.TransactionOptions{}
	for _, opt := range opts {
		opt(&options)
	}

	query := goqu.
		From(schema.Transaction.Table()).
		Select(goqu.COUNT(schema.Transaction.ColumnAll())).
		Where(goqu.I(schema.Transaction.ColumnDeletedAt()).IsNull())

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
	opts ...repo.TransactionOption,
) (int64, error) {
	options := repo.TransactionOptions{}
	for _, opt := range opts {
		opt(&options)
	}

	query := goqu.
		From(schema.Transaction.Table()).
		Select(goqu.SUM(schema.Transaction.ColumnAmount())).
		Where(goqu.I(schema.Transaction.ColumnDeletedAt()).IsNull())

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
	opts ...repo.TransactionOption,
) (map[uuid.UUID]int64, error) {
	options := repo.TransactionOptions{}
	for _, opt := range opts {
		opt(&options)
	}

	query := goqu.
		From(schema.Transaction.Table()).
		Select(
			goqu.I(schema.Transaction.ColumnCategoryID()),
			goqu.SUM(schema.Transaction.ColumnAmount()).As("sum"),
		).
		Where(goqu.I(schema.Transaction.ColumnDeletedAt()).IsNull()).
		GroupBy(goqu.I(schema.Transaction.ColumnCategoryID()))

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
		goqu.I(schema.Transaction.ColumnUserID()).Eq(userID),
	)

	options.Search = strings.TrimSpace(options.Search)
	if options.Search != "" {
		searchExp, orderExp := qb.buildSearch(
			options.Search,
			schema.Transaction.ColumnName(),
			schema.TransactionCategory.ColumnName(),
			schema.Institution.ColumnName(),
			schema.PaymentMethod.ColumnName(),
		)
		whereExps = append(whereExps, searchExp)
		orderedExps = append(orderedExps, orderExp.Desc())
	}

	if len(options.CategoryIDs) > 0 {
		whereExps = append(
			whereExps,
			goqu.I(schema.Transaction.ColumnCategoryID()).
				In(options.CategoryIDs),
		)
	}

	if len(options.InstitutionIDs) > 0 {
		whereExps = append(
			whereExps,
			goqu.I(schema.Transaction.ColumnInstitutionID()).
				In(options.InstitutionIDs),
		)
	}

	if len(options.PaymentMethodIDs) > 0 {
		whereExps = append(
			whereExps,
			goqu.I(schema.Transaction.ColumnPaymentMethodID()).
				In(options.PaymentMethodIDs),
		)
	}

	if !options.StartDate.IsZero() {
		whereExps = append(
			whereExps,
			goqu.I(schema.Transaction.ColumnDate()).Gte(options.StartDate),
		)
	}

	if !options.EndDate.IsZero() {
		whereExps = append(
			whereExps,
			goqu.I(schema.Transaction.ColumnDate()).Lte(options.EndDate),
		)
	}

	if options.IsExpense {
		whereExps = append(
			whereExps,
			goqu.I(schema.Transaction.ColumnAmount()).Lt(0),
		)
	}

	if options.IsIncome {
		whereExps = append(
			whereExps,
			goqu.I(schema.Transaction.ColumnAmount()).Gt(0),
		)
	}

	if options.IsIgnored != nil {
		whereExps = append(
			whereExps,
			goqu.I(schema.Transaction.ColumnIsIgnored()).
				Eq(*options.IsIgnored),
		)
	}

	orderedExps = append(
		orderedExps,
		goqu.I(schema.Transaction.ColumnName()).Asc(),
	)

	return whereExps, orderedExps
}

func (qb *QueryBuilder) buildTransactionJoins(
	options repo.TransactionOptions,
	shouldJoinAll ...bool,
) []Join {
	institutionJoin := Join{
		Table: goqu.I(schema.Institution.Table()),
		Condition: goqu.
			On(
				goqu.I(schema.Transaction.ColumnInstitutionID()).
					Eq(goqu.I(schema.Institution.ColumnID())),
			),
	}

	transactionCategoryJoin := Join{
		Table: goqu.I(schema.TransactionCategory.Table()),
		Condition: goqu.
			On(
				goqu.I(schema.Transaction.ColumnCategoryID()).
					Eq(goqu.I(schema.TransactionCategory.ColumnID())),
			),
	}

	paymentMethodJoin := Join{
		Table: goqu.I(schema.PaymentMethod.Table()),
		Condition: goqu.
			On(
				goqu.I(schema.Transaction.ColumnPaymentMethodID()).
					Eq(goqu.I(schema.PaymentMethod.ColumnID())),
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
