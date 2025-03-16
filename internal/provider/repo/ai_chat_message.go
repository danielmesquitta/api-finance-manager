package repo

import (
	"context"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/google/uuid"
)

type AIChatMessageRepo interface {
	CreateAIChatMessage(
		ctx context.Context,
		params CreateAIChatMessageParams,
	) (*entity.AIChatMessage, error)
	DeleteAIChatMessages(ctx context.Context, aiChatID uuid.UUID) error
}
