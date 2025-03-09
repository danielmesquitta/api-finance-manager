package query

import (
	"context"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/db"
	"github.com/doug-martin/goqu/v9"
)

func (tqb *TestQueryBuilder) GetUserByAuthID(
	ctx context.Context,
	authID string,
) (*entity.User, error) {
	query := goqu.
		Select(db.TableUser.ColumnAll()).
		From(db.TableUser.String()).
		Where(goqu.Ex{db.TableUser.ColumnAuthID(): authID})

	dest := &entity.User{}
	if err := tqb.Scan(ctx, query, dest); err != nil {
		return nil, err
	}

	return dest, nil
}
