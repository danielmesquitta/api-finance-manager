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
	opts ...repo.ListTransactionsOption,
) ([]entity.Transaction, error) {
	options := repo.ListTransactionsOptions{}
	for _, opt := range opts {
		opt(&options)
	}

	query := goqu.
		From(TableTransaction).
		Select("*").
		Where(goqu.Ex{ColumnTransactionDeletedAt: nil})

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
	opts ...repo.ListTransactionsOption,
) ([]entity.TransactionWithCategoryAndInstitution, error) {
	options := repo.ListTransactionsOptions{}
	for _, opt := range opts {
		opt(&options)
	}

	query := goqu.
		From(TableTransaction).
		Select(
			fmt.Sprintf("%s.*", TableTransaction),
			goqu.I(ColumnCategoryName).As("category_name"),
			goqu.I(ColumnInstitutionName).As("institution_name"),
			goqu.I(ColumnInstitutionLogo).As("institution_logo"),
		).
		LeftJoin(
			goqu.I(TableCategory),
			goqu.
				On(
					goqu.I(ColumnTransactionCategoryID).
						Eq(goqu.I(ColumnCategoryID)),
				),
		).
		LeftJoin(
			goqu.I(TableInstitution),
			goqu.
				On(
					goqu.I(ColumnTransactionInstitutionID).
						Eq(goqu.I(ColumnInstitutionID)),
				),
		).
		Where(goqu.Ex{ColumnTransactionDeletedAt: nil})

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

	var transactions []entity.TransactionWithCategoryAndInstitution
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
	opts ...repo.ListTransactionsOption,
) (int64, error) {
	options := repo.ListTransactionsOptions{}
	for _, opt := range opts {
		opt(&options)
	}

	query := goqu.
		From(TableTransaction).
		Select(goqu.COUNT("*")).
		Where(goqu.Ex{ColumnTransactionDeletedAt: nil})

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

func (qb *QueryBuilder) buildTransactionExpressions(
	userID uuid.UUID,
	options repo.ListTransactionsOptions,
) (whereExps []goqu.Expression, orderedExps []exp.OrderedExpression) {
	whereExps = append(
		whereExps,
		goqu.I(ColumnTransactionUserID).Eq(userID),
	)

	options.Search = strings.TrimSpace(options.Search)
	if options.Search != "" {
		searchExp, distanceExp := qb.buildSearch(
			options.Search,
			ColumnTransactionName,
		)
		whereExps = append(whereExps, searchExp)
		orderedExps = append(orderedExps, distanceExp.Asc())
	}

	if options.CategoryID != uuid.Nil {
		whereExps = append(
			whereExps,
			goqu.I(ColumnTransactionCategoryID).
				Eq(options.CategoryID),
		)
	}

	if options.InstitutionID != uuid.Nil {
		whereExps = append(
			whereExps,
			goqu.I(ColumnTransactionInstitutionID).
				Eq(options.InstitutionID),
		)
	}

	if !options.StartDate.IsZero() {
		whereExps = append(
			whereExps,
			goqu.I(ColumnTransactionDate).Gte(options.StartDate),
		)
	}

	if !options.EndDate.IsZero() {
		whereExps = append(
			whereExps,
			goqu.I(ColumnTransactionDate).Lte(options.EndDate),
		)
	}

	if options.IsExpense {
		whereExps = append(
			whereExps,
			goqu.I(ColumnTransactionAmount).Lt(0),
		)
	}

	if options.IsIncome {
		whereExps = append(
			whereExps,
			goqu.I(ColumnTransactionAmount).Gt(0),
		)
	}

	if options.IsIgnored != nil {
		whereExps = append(
			whereExps,
			goqu.I(ColumnTransactionIsIgnored).Eq(*options.IsIgnored),
		)
	}

	if options.PaymentMethodID != uuid.Nil {
		whereExps = append(
			whereExps,
			goqu.I(ColumnTransactionPaymentMethodID).
				Eq(options.PaymentMethodID),
		)
	}

	orderedExps = append(
		orderedExps,
		goqu.I(ColumnTransactionName).Asc(),
	)

	return whereExps, orderedExps
}

func (qb *QueryBuilder) buildTransactionsQuery(
	query *goqu.SelectDataset,
	options repo.ListTransactionsOptions,
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
