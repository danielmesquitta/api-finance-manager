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
	"github.com/jackc/pgx/v5"
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
		From(tableTransaction).
		Select("*").
		Where(goqu.Ex{tableTransaction.ColumnDeletedAt(): nil})

	whereExps, orderedExps := qb.buildTransactionExpressions(userID, options)

	query = qb.buildTransactionsQuery(query, options, whereExps, orderedExps)

	sql, args, err := query.ToSQL()
	if err != nil {
		return nil, errs.New(err)
	}

	var transactions []entity.Transaction
	if err := pgxscan.Select(ctx, qb.db, &transactions, sql, args...); err != nil {
		return nil, errs.New(err)
	}

	return transactions, nil
}

func (qb *QueryBuilder) ListTransactionsWithCategoriesAndInstitutions(
	ctx context.Context,
	userID uuid.UUID,
	opts ...repo.TransactionOption,
) ([]entity.FullTransaction, error) {
	options := repo.TransactionOptions{}
	for _, opt := range opts {
		opt(&options)
	}

	query := goqu.
		From(tableTransaction).
		Select(
			fmt.Sprintf("%s.*", tableTransaction),
			goqu.I(tableTransactionCategory.ColumnName()).As("category_name"),
			goqu.I(tableInstitution.ColumnName()).As("institution_name"),
			goqu.I(tableInstitution.ColumnLogo()).As("institution_logo"),
			goqu.I(tablePaymentMethod.ColumnName()).As("payment_method_name"),
		).
		LeftJoin(
			goqu.I(tableTransactionCategory.String()),
			goqu.
				On(
					goqu.I(tableTransaction.ColumnCategoryID()).
						Eq(goqu.I(tableTransactionCategory.ColumnID())),
				),
		).
		LeftJoin(
			goqu.I(tableInstitution.String()),
			goqu.
				On(
					goqu.I(tableTransaction.ColumnInstitutionID()).
						Eq(goqu.I(tableInstitution.ColumnID())),
				),
		).
		LeftJoin(
			goqu.I(tablePaymentMethod.String()),
			goqu.
				On(
					goqu.I(tableTransaction.ColumnPaymentMethodID()).
						Eq(goqu.I(tablePaymentMethod.ColumnID())),
				),
		).
		Where(goqu.Ex{tableTransaction.ColumnDeletedAt(): nil})

	whereExps, orderedExps := qb.buildTransactionExpressions(userID, options)

	query = qb.buildTransactionsQuery(
		query,
		options,
		whereExps,
		orderedExps,
	)

	sql, args, err := query.ToSQL()
	if err != nil {
		return nil, errs.New(err)
	}

	var transactions []entity.FullTransaction
	if err := pgxscan.Select(ctx, qb.db, &transactions, sql, args...); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
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
		From(tableTransaction).
		Select(goqu.COUNT("*")).
		Where(goqu.Ex{tableTransaction.ColumnDeletedAt(): nil})

	whereExps, _ := qb.buildTransactionExpressions(userID, options)

	query = qb.buildTransactionsQuery(query, options, whereExps, nil)

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
		From(tableTransaction).
		Select(goqu.SUM(tableTransaction.ColumnAmount())).
		Where(goqu.Ex{tableTransaction.ColumnDeletedAt(): nil})

	whereExps, _ := qb.buildTransactionExpressions(userID, options)

	query = qb.buildTransactionsQuery(query, options, whereExps, nil)

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
		From(tableTransaction).
		Select(
			goqu.I(tableTransaction.ColumnCategoryID()),
			goqu.SUM(tableTransaction.ColumnAmount()).As("sum"),
		).
		Where(goqu.Ex{tableTransaction.ColumnDeletedAt(): nil}).
		GroupBy(goqu.I(tableTransaction.ColumnCategoryID()))

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
		goqu.I(tableTransaction.ColumnUserID()).Eq(userID),
	)

	options.Search = strings.TrimSpace(options.Search)
	if options.Search != "" {
		searchExp, distanceExp := qb.buildSearch(
			options.Search,
			tableTransaction.ColumnName(),
		)
		whereExps = append(whereExps, searchExp)
		orderedExps = append(orderedExps, distanceExp.Asc())
	}

	if len(options.CategoryIDs) > 0 {
		whereExps = append(
			whereExps,
			goqu.I(tableTransaction.ColumnCategoryID()).
				In(options.CategoryIDs),
		)
	}

	if len(options.InstitutionIDs) > 0 {
		whereExps = append(
			whereExps,
			goqu.I(tableTransaction.ColumnInstitutionID()).
				In(options.InstitutionIDs),
		)
	}

	if len(options.PaymentMethodIDs) > 0 {
		whereExps = append(
			whereExps,
			goqu.I(tableTransaction.ColumnPaymentMethodID()).
				In(options.PaymentMethodIDs),
		)
	}

	if !options.StartDate.IsZero() {
		whereExps = append(
			whereExps,
			goqu.I(tableTransaction.ColumnDate()).Gte(options.StartDate),
		)
	}

	if !options.EndDate.IsZero() {
		whereExps = append(
			whereExps,
			goqu.I(tableTransaction.ColumnDate()).Lte(options.EndDate),
		)
	}

	if options.IsExpense {
		whereExps = append(
			whereExps,
			goqu.I(tableTransaction.ColumnAmount()).Lt(0),
		)
	}

	if options.IsIncome {
		whereExps = append(
			whereExps,
			goqu.I(tableTransaction.ColumnAmount()).Gt(0),
		)
	}

	if options.IsIgnored != nil {
		whereExps = append(
			whereExps,
			goqu.I(tableTransaction.ColumnIsIgnored()).Eq(*options.IsIgnored),
		)
	}

	orderedExps = append(
		orderedExps,
		goqu.I(tableTransaction.ColumnName()).Asc(),
	)

	return whereExps, orderedExps
}

func (qb *QueryBuilder) buildTransactionsQuery(
	query *goqu.SelectDataset,
	options repo.TransactionOptions,
	whereExps []goqu.Expression,
	orderedExps []exp.OrderedExpression,
) *goqu.SelectDataset {
	if len(whereExps) == 1 {
		query = query.
			Where(whereExps[0])
	} else if len(whereExps) > 0 {
		query = query.
			Where(goqu.And(whereExps...))
	}

	if len(orderedExps) > 0 {
		query = query.
			Order(orderedExps...)
	}

	if options.Limit > 0 {
		query = query.Limit(options.Limit)
	}

	if options.Offset > 0 {
		query = query.Offset(options.Offset)
	}

	return query
}
