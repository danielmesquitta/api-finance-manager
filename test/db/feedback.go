package db

import (
	"context"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/db/schema"
	"github.com/doug-martin/goqu/v9"
)

func (tdb *TestDB) GetLatestFeedbackByUserID(
	ctx context.Context,
	userID string,
) (*entity.Feedback, error) {
	query := goqu.
		Select(schema.Feedback.ColumnAll()).
		From(schema.Feedback.Table()).
		Where(goqu.Ex{schema.Feedback.ColumnUserID(): userID}).
		Order(goqu.I(schema.Feedback.ColumnCreatedAt()).Desc()).
		Limit(1)

	dest := &entity.Feedback{}
	if err := tdb.Scan(ctx, query, dest); err != nil {
		return nil, err
	}

	return dest, nil
}
