package db

import (
	"context"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/db/schema"
	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
)

func (tdb *TestDB) GetLatestUserDeletedAIChat(
	ctx context.Context,
	userID uuid.UUID,
) (*entity.AIChat, error) {
	query := goqu.
		Select(schema.AIChat.All()).
		From(schema.AIChat.String()).
		Where(
			goqu.I(schema.AIChat.UserID()).Eq(userID),
			goqu.I(schema.AIChat.DeletedAt()).IsNotNull(),
		).
		Order(goqu.I(schema.AIChat.DeletedAt()).Desc()).
		Limit(1)

	dest := &entity.AIChat{}
	if err := tdb.Scan(ctx, query, dest); err != nil {
		return nil, err
	}

	return dest, nil
}
