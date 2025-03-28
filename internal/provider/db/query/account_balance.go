package query

import (
	"context"
	"time"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/db/schema"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
)

func (qb *QueryBuilder) GetUserBalanceOnDate(
	ctx context.Context,
	userID uuid.UUID,
	date time.Time,
	options ...repo.AccountBalanceOptions,
) (int64, error) {
	opts := prepareOptions(options...)

	subQuery := goqu.
		From(schema.AccountBalance.Table()).
		Select(schema.AccountBalance.ColumnAmount()).
		Where(
			goqu.I(schema.AccountBalance.ColumnAccountID()).Eq(goqu.I("a.id")),
			goqu.I(schema.AccountBalance.ColumnCreatedAt()).Lte(date),
			goqu.I(schema.AccountBalance.ColumnDeletedAt()).IsNull(),
		).
		Order(goqu.I(schema.AccountBalance.ColumnCreatedAt()).Desc()).
		Limit(1)

	query := goqu.
		From(goqu.I(schema.Account.Table())).
		Select(goqu.L("COALESCE(SUM(ab.amount), 0)::bigint AS total_balance")).
		LeftJoin(
			goqu.Lateral(subQuery).As("ab"),
			goqu.On(goqu.L("TRUE")),
		).
		Join(
			goqu.I(schema.UserInstitution.Table()),
			goqu.On(
				goqu.I(schema.Account.ColumnUserInstitutionID()).
					Eq(goqu.I(schema.UserInstitution.ColumnID())),
			)).
		Where(
			goqu.Ex{
				schema.Account.ColumnType():           entity.AccountTypeBank,
				schema.UserInstitution.ColumnUserID(): userID,
			},
			goqu.I(schema.Account.ColumnDeletedAt()).IsNull(),
			goqu.I(schema.UserInstitution.ColumnDeletedAt()).IsNull(),
		)

	if len(opts.InstitutionIDs) > 0 {
		query = query.Where(
			goqu.I(schema.UserInstitution.ColumnInstitutionID()).
				In(opts.InstitutionIDs),
		)
	}

	var totalBalance int64
	if err := qb.Scan(ctx, query, &totalBalance); err != nil {
		return 0, errs.New(err)
	}

	return totalBalance, nil
}
