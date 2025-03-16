package repo

import (
	"context"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/google/uuid"
)

type AIChatAnswerRepo interface {
	CreateAIChatAnswer(
		ctx context.Context,
		params CreateAIChatAnswerParams,
	) (*entity.AIChatAnswer, error)
	DeleteAIChatAnswers(ctx context.Context, aiChatID uuid.UUID) error
	UpdateAIChatAnswer(
		ctx context.Context,
		params UpdateAIChatAnswerParams,
	) error
}
