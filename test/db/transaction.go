package db

import (
	"context"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/db/schema"
	"github.com/doug-martin/goqu/v9"
)

func (tdb *TestDB) GetTransactionByID(
	ctx context.Context,
	id string,
) (*entity.Transaction, error) {
	query := goqu.
		Select(schema.Transaction.ColumnAll()).
		From(schema.Transaction.Table()).
		Where(goqu.Ex{schema.Transaction.ColumnID(): id})

	dest := &entity.Transaction{}
	if err := tdb.Scan(ctx, query, dest); err != nil {
		return nil, err
	}

	return dest, nil
}

func (tdb *TestDB) GetLatestTransactionByUserID(
	ctx context.Context,
	userID string,
) (*entity.Transaction, error) {
	query := goqu.
		Select(schema.Transaction.ColumnAll()).
		From(schema.Transaction.Table()).
		Where(goqu.Ex{schema.Transaction.ColumnUserID(): userID}).
		Order(goqu.I(schema.Transaction.ColumnCreatedAt()).Desc()).
		Limit(1)

	dest := &entity.Transaction{}
	if err := tdb.Scan(ctx, query, dest); err != nil {
		return nil, err
	}

	return dest, nil
}
