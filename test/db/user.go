package db

import (
	"context"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/db/schema"
	"github.com/doug-martin/goqu/v9"
)

func (tdb *TestDB) GetUserByAuthID(
	ctx context.Context,
	authID string,
) (*entity.User, error) {
	query := goqu.
		Select(schema.User.ColumnAll()).
		From(schema.User.Table()).
		Where(goqu.Ex{schema.User.ColumnAuthID(): authID})

	dest := &entity.User{}
	if err := tdb.Scan(ctx, query, dest); err != nil {
		return nil, err
	}

	return dest, nil
}
