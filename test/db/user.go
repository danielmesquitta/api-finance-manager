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
		Select(schema.User.All()).
		From(schema.User.String()).
		Where(
			goqu.I(schema.User.ID()).Eq(userID),
			goqu.I(schema.User.DeletedAt()).IsNotNull(),
		).
		Order(goqu.I(schema.User.DeletedAt()).Desc()).
		Limit(1)

	dest := &entity.User{}
	if err := tdb.Scan(ctx, query, dest); err != nil {
		return nil, err
	}

	return dest, nil
}
