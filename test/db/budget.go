package db

import (
	"context"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/db/schema"
	"github.com/doug-martin/goqu/v9"
)

func (tdb *TestDB) ListBudgets(
	ctx context.Context,
) ([]entity.Budget, error) {
	query := goqu.
		Select(schema.Budget.ColumnAll()).
		From(schema.Budget.Table())

	dest := []entity.Budget{}
	if err := tdb.Scan(ctx, query, &dest); err != nil {
		return nil, err
	}

	return dest, nil
}
