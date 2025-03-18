package repo

import (
	"context"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/google/uuid"
)

type AIChatMessageRepo interface {
	GenerateAIChatMessage(
		ctx context.Context,
		params GenerateAIChatMessageParams,
	) (*entity.AIChatMessage, error)
	DeleteAIChatMessages(ctx context.Context, aiChatID uuid.UUID) error
}
