package query

import (
	"context"
	"strings"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/db"
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
		From(db.TableTransaction.String()).
		Select(db.TableTransaction.ColumnAll()).
		Where(goqu.I(db.TableTransaction.ColumnDeletedAt()).IsNull())

	whereExps, orderedExps := qb.buildTransactionExpressions(userID, options)

	query = qb.buildTransactionsQuery(query, options, whereExps, orderedExps)

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
		From(db.TableTransaction.String()).
		Select(
			db.TableTransaction.ColumnAll(),
			goqu.I(db.TableTransactionCategory.ColumnName()).
				As("category_name"),
			goqu.I(db.TableInstitution.ColumnName()).As("institution_name"),
			goqu.I(db.TableInstitution.ColumnLogo()).As("institution_logo"),
			goqu.I(db.TablePaymentMethod.ColumnName()).
				As("payment_method_name"),
		).
		LeftJoin(
			goqu.I(db.TableTransactionCategory.String()),
			goqu.
				On(
					goqu.I(db.TableTransaction.ColumnCategoryID()).
						Eq(goqu.I(db.TableTransactionCategory.ColumnID())),
				),
		).
		LeftJoin(
			goqu.I(db.TableInstitution.String()),
			goqu.
				On(
					goqu.I(db.TableTransaction.ColumnInstitutionID()).
						Eq(goqu.I(db.TableInstitution.ColumnID())),
				),
		).
		LeftJoin(
			goqu.I(db.TablePaymentMethod.String()),
			goqu.
				On(
					goqu.I(db.TableTransaction.ColumnPaymentMethodID()).
						Eq(goqu.I(db.TablePaymentMethod.ColumnID())),
				),
		).
		Where(goqu.I(db.TableTransaction.ColumnDeletedAt()).IsNull())

	whereExps, orderedExps := qb.buildTransactionExpressions(userID, options)

	query = qb.buildTransactionsQuery(
		query,
		options,
		whereExps,
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
		From(db.TableTransaction.String()).
		Select(goqu.COUNT(db.TableTransaction.ColumnAll())).
		Where(goqu.I(db.TableTransaction.ColumnDeletedAt()).IsNull())

	whereExps, _ := qb.buildTransactionExpressions(userID, options)

	query = qb.buildTransactionsQuery(query, options, whereExps, nil)

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
		From(db.TableTransaction.String()).
		Select(goqu.SUM(db.TableTransaction.ColumnAmount())).
		Where(goqu.I(db.TableTransaction.ColumnDeletedAt()).IsNull())

	whereExps, _ := qb.buildTransactionExpressions(userID, options)

	query = qb.buildTransactionsQuery(query, options, whereExps, nil)

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
		From(db.TableTransaction.String()).
		Select(
			goqu.I(db.TableTransaction.ColumnCategoryID()),
			goqu.SUM(db.TableTransaction.ColumnAmount()).As("sum"),
		).
		Where(goqu.I(db.TableTransaction.ColumnDeletedAt()).IsNull()).
		GroupBy(goqu.I(db.TableTransaction.ColumnCategoryID()))

	whereExps, _ := qb.buildTransactionExpressions(userID, options)

	query = qb.buildTransactionsQuery(query, options, whereExps, nil)

	sql, args, err := query.ToSQL()
	if err != nil {
		return nil, errs.New(err)
	}

	rows, err := qb.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, errs.New(err)
	}
	defer rows.Close()

	out := map[uuid.UUID]int64{}
	for rows.Next() {
		row := struct {
			CategoryID uuid.UUID
			Sum        int64
		}{}

		if err := rows.Scan(&row.CategoryID, &row.Sum); err != nil {
			return nil, err
		}

		out[row.CategoryID] = row.Sum
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return out, nil
}

func (qb *QueryBuilder) buildTransactionExpressions(
	userID uuid.UUID,
	options repo.TransactionOptions,
) (whereExps []goqu.Expression, orderedExps []exp.OrderedExpression) {
	whereExps = append(
		whereExps,
		goqu.I(db.TableTransaction.ColumnUserID()).Eq(userID),
	)

	options.Search = strings.TrimSpace(options.Search)
	if options.Search != "" {
		searchExp, distanceExp := qb.buildSearch(
			options.Search,
			db.TableTransaction.ColumnName(),
		)
		whereExps = append(whereExps, searchExp)
		orderedExps = append(orderedExps, distanceExp.Asc())
	}

	if len(options.CategoryIDs) > 0 {
		whereExps = append(
			whereExps,
			goqu.I(db.TableTransaction.ColumnCategoryID()).
				In(options.CategoryIDs),
		)
	}

	if len(options.InstitutionIDs) > 0 {
		whereExps = append(
			whereExps,
			goqu.I(db.TableTransaction.ColumnInstitutionID()).
				In(options.InstitutionIDs),
		)
	}

	if len(options.PaymentMethodIDs) > 0 {
		whereExps = append(
			whereExps,
			goqu.I(db.TableTransaction.ColumnPaymentMethodID()).
				In(options.PaymentMethodIDs),
		)
	}

	if !options.StartDate.IsZero() {
		whereExps = append(
			whereExps,
			goqu.I(db.TableTransaction.ColumnDate()).Gte(options.StartDate),
		)
	}

	if !options.EndDate.IsZero() {
		whereExps = append(
			whereExps,
			goqu.I(db.TableTransaction.ColumnDate()).Lte(options.EndDate),
		)
	}

	if options.IsExpense {
		whereExps = append(
			whereExps,
			goqu.I(db.TableTransaction.ColumnAmount()).Lt(0),
		)
	}

	if options.IsIncome {
		whereExps = append(
			whereExps,
			goqu.I(db.TableTransaction.ColumnAmount()).Gt(0),
		)
	}

	if options.IsIgnored != nil {
		whereExps = append(
			whereExps,
			goqu.I(db.TableTransaction.ColumnIsIgnored()).
				Eq(*options.IsIgnored),
		)
	}

	orderedExps = append(
		orderedExps,
		goqu.I(db.TableTransaction.ColumnName()).Asc(),
	)

	return whereExps, orderedExps
}

func (qb *QueryBuilder) buildTransactionsQuery(
	query *goqu.SelectDataset,
	options repo.TransactionOptions,
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
