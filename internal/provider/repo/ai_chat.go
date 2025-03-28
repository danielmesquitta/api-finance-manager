package repo

import (
	"context"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/google/uuid"
)

type AIChatOptions struct {
	Limit  uint      `json:"-"`
	Offset uint      `json:"-"`
	UserID uuid.UUID `json:"-"`
	Search string    `json:"search"`
}

type AIChatRepo interface {
	CreateAIChat(
		ctx context.Context,
		userID uuid.UUID,
	) (*entity.AIChat, error)
	DeleteAIChat(
		ctx context.Context,
		id uuid.UUID,
	) error
	GetAIChatByID(
		ctx context.Context,
		id uuid.UUID,
	) (*entity.AIChat, error)
	GetLatestAIChatByUserID(
		ctx context.Context,
		userID uuid.UUID,
	) (*entity.FullAIChat, error)
	ListAIChats(
		ctx context.Context,
		opts ...AIChatOptions,
	) ([]entity.AIChat, error)
	CountAIChats(
		ctx context.Context,
		opts ...AIChatOptions,
	) (int64, error)
	UpdateAIChat(
		ctx context.Context,
		params UpdateAIChatParams,
	) error
	ListAIChatMessagesAndAnswers(
		ctx context.Context,
		params ListAIChatMessagesAndAnswersParams,
	) ([]entity.AIChatMessageAndAnswer, error)
	CountAIChatMessagesAndAnswers(
		ctx context.Context,
		aiChatID uuid.UUID,
	) (int64, error)
}
