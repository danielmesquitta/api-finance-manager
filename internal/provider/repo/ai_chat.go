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

type AIChatOption func(*AIChatOptions)

func WithAIChatPagination(
	limit uint,
	offset uint,
) AIChatOption {
	return func(o *AIChatOptions) {
		o.Limit = limit
		o.Offset = offset
	}
}

func WithAIChatUser(userID uuid.UUID) AIChatOption {
	return func(o *AIChatOptions) {
		o.UserID = userID
	}
}

func WithAIChatSearch(search string) AIChatOption {
	return func(o *AIChatOptions) {
		o.Search = search
	}
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
	GetAIChat(
		ctx context.Context,
		id uuid.UUID,
	) (*entity.AIChat, error)
	GetLatestAIChatByUserID(
		ctx context.Context,
		userID uuid.UUID,
	) (*entity.FullAIChat, error)
	ListAIChats(
		ctx context.Context,
		opts ...AIChatOption,
	) ([]entity.AIChat, error)
	CountAIChats(
		ctx context.Context,
		opts ...AIChatOption,
	) (int64, error)
	UpdateAIChat(
		ctx context.Context,
		params UpdateAIChatParams,
	) error
}
