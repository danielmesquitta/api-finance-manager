package dto

import (
	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/usecase"
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
	usecase.UpdateAIChatInput
}

type GenerateAIChatMessageRequest struct {
	usecase.GenerateAIChatMessageInput
}
