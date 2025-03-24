package dto

import (
	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/usecase/aichat"
)

type ListAIChatsResponse struct {
	entity.PaginatedList[entity.AIChat]
}

type ListAIChatMessagesAndAnswersResponse struct {
	entity.PaginatedList[entity.AIChatMessageAndAnswer]
}

type CreateAIChatResponse struct {
	entity.AIChat
}

type UpdateAIChatRequest struct {
	aichat.UpdateAIChatUseCaseInput
}

type GenerateAIChatMessageRequest struct {
	aichat.GenerateAIChatMessageUseCaseInput
}

type GenerateAIChatMessageResponse struct {
	entity.AIChatAnswer
}
