package repo

import (
	"context"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/google/uuid"
)

type AIChatMessageOptions struct {
	Limit    uint      `json:"-"`
	Offset   uint      `json:"-"`
	AIChatID uuid.UUID `json:"ai_chat_id"`
}

type AIChatMessageOption func(*AIChatMessageOptions)

func WithAIChatMessagePagination(
	limit uint,
	offset uint,
) AIChatMessageOption {
	return func(o *AIChatMessageOptions) {
		o.Limit = limit
		o.Offset = offset
	}
}

func WithAIChatMessageAIChat(aiChatID uuid.UUID) AIChatMessageOption {
	return func(o *AIChatMessageOptions) {
		o.AIChatID = aiChatID
	}
}

type AIChatMessageRepo interface {
	ListAIChatMessages(
		ctx context.Context,
		opts ...AIChatMessageOption,
	) ([]entity.AIChatMessage, error)
	CountAIChatMessages(
		ctx context.Context,
		opts ...AIChatMessageOption,
	) (int64, error)
}
