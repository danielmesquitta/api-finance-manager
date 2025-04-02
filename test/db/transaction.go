package db

import (
	"context"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/db/schema"
	"github.com/doug-martin/goqu/v9"
)

func (tdb *TestDB) GetLatestTransactionByUserID(
	ctx context.Context,
	userID string,
) (*entity.Transaction, error) {
	query := goqu.
		Select(schema.Transaction.All()).
		From(schema.Transaction.String()).
		Where(goqu.Ex{schema.Transaction.UserID(): userID}).
		Order(goqu.I(schema.Transaction.CreatedAt()).Desc()).
		Limit(1)

	dest := &entity.Transaction{}
	if err := tdb.Scan(ctx, query, dest); err != nil {
		return nil, err
	}

	return dest, nil
}
