package db

import (
	"context"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/db/schema"
	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
)

func (tdb *TestDB) GetLatestDeletedUser(
	ctx context.Context,
	userID uuid.UUID,
) (*entity.User, error) {
	query := goqu.
		Select(schema.User.ColumnAll()).
		From(schema.User.Table()).
		Where(
			goqu.I(schema.User.ColumnID()).Eq(userID),
			goqu.I(schema.User.ColumnDeletedAt()).IsNotNull(),
		).
		Order(goqu.I(schema.User.ColumnDeletedAt()).Desc()).
		Limit(1)

	dest := &entity.User{}
	if err := tdb.Scan(ctx, query, dest); err != nil {
		return nil, err
	}

	return dest, nil
}
