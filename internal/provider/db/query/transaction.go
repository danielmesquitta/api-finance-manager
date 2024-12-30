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
		From(string(TableTransaction)).
		Select("*")

	whereExps, orderedExps := qb.buildTransactionExpressions(userID, options)

	qb.setTransactionsQuery(query, options, whereExps, orderedExps)

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
		From(string(TableTransaction)).
		Select(
			"*",
			goqu.I(fmt.Sprintf("%s.%s", TableCategory, ColumnCategoryName)).
				As("category_name"),
			goqu.I(fmt.Sprintf("%s.%s", TableInstitution, ColumnInstitutionName)).
				As("institution_name"),
			goqu.I(fmt.Sprintf("%s.%s", TableInstitution, ColumnInstitutionLogo)).
				As("institution_logo"),
		).
		LeftJoin(
			goqu.I(string(TableCategory)),
			goqu.
				On(
					goqu.I(fmt.Sprintf("%s.%s", TableTransaction, ColumnTransactionCategoryID)).
						Eq(goqu.I(fmt.Sprintf("%s.%s", TableCategory, ColumnCategoryID))),
				),
		).
		LeftJoin(
			goqu.I(string(TableInstitution)),
			goqu.
				On(
					goqu.I(fmt.Sprintf("%s.%s", TableTransaction, ColumnTransactionInstitutionID)).
						Eq(goqu.I(fmt.Sprintf("%s.%s", TableInstitution, ColumnInstitutionID))),
				),
		)

	whereExps, orderedExps := qb.buildTransactionExpressions(userID, options)

	qb.setTransactionsQuery(query, options, whereExps, orderedExps)

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
		From(string(TableTransaction)).
		Select(goqu.COUNT("*"))

	whereExps, _ := qb.buildTransactionExpressions(userID, options)

	qb.setTransactionsQuery(query, options, whereExps, nil)

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
		goqu.I(string(ColumnTransactionUserID)).Eq(userID),
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
			goqu.I(string(ColumnTransactionCategoryID)).
				Eq(options.CategoryID),
		)
	}

	if options.InstitutionID != uuid.Nil {
		whereExps = append(
			whereExps,
			goqu.I(string(ColumnTransactionInstitutionID)).
				Eq(options.InstitutionID),
		)
	}

	if !options.StartDate.IsZero() {
		whereExps = append(
			whereExps,
			goqu.I(string(ColumnTransactionDate)).Gte(options.StartDate),
		)
	}

	if !options.EndDate.IsZero() {
		whereExps = append(
			whereExps,
			goqu.I(string(ColumnTransactionDate)).Lte(options.EndDate),
		)
	}

	orderedExps = append(
		orderedExps,
		goqu.I(string(ColumnTransactionName)).Asc(),
	)

	return whereExps, orderedExps
}

func (qb *QueryBuilder) setTransactionsQuery(
	query *goqu.SelectDataset,
	options repo.ListTransactionsOptions,
	whereExps []goqu.Expression,
	orderedExps []exp.OrderedExpression,
) {
	if len(whereExps) == 1 {
		query.
			Where(whereExps[0])
	} else if len(whereExps) > 0 {
		query.
			Where(goqu.And(whereExps...))
	}

	if len(orderedExps) > 0 {
		query.
			Order(orderedExps...)
	}

	if options.Limit > 0 {
		query.Limit(options.Limit)
	}

	if options.Offset > 0 {
		query.Offset(options.Offset)
	}
}
