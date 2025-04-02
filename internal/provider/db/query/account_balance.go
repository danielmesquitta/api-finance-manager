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
		From(schema.AccountBalance.String()).
		Select(schema.AccountBalance.Amount()).
		Where(
			goqu.I(schema.AccountBalance.AccountID()).
				Eq(goqu.I(schema.Account.ID())),
			goqu.I(schema.AccountBalance.CreatedAt()).Lte(date),
			goqu.I(schema.AccountBalance.DeletedAt()).IsNull(),
		).
		Order(goqu.I(schema.AccountBalance.CreatedAt()).Desc()).
		Limit(1)

	query := goqu.
		From(schema.Account.String()).
		Select(goqu.L("COALESCE(SUM(ab.amount), 0)::bigint AS total_balance")).
		LeftJoin(
			goqu.Lateral(subQuery).As("ab"),
			goqu.On(goqu.L("TRUE")),
		).
		Join(
			goqu.I(schema.UserInstitution.String()),
			goqu.On(
				goqu.I(schema.Account.UserInstitutionID()).
					Eq(goqu.I(schema.UserInstitution.ID())),
			)).
		Where(
			goqu.Ex{
				schema.Account.Type():           entity.AccountTypeBank,
				schema.UserInstitution.UserID(): userID,
			},
			goqu.I(schema.Account.DeletedAt()).IsNull(),
			goqu.I(schema.UserInstitution.DeletedAt()).IsNull(),
		)

	if len(opts.InstitutionIDs) > 0 {
		query = query.Where(
			goqu.I(schema.UserInstitution.InstitutionID()).
				In(opts.InstitutionIDs),
		)
	}

	var totalBalance int64
	if err := qb.Scan(ctx, query, &totalBalance); err != nil {
		return 0, errs.New(err)
	}

	return totalBalance, nil
}
