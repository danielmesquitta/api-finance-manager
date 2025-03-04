package handler

import (
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/usecase"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AIChatMessageHandler struct {
	lacm *usecase.ListAIChatMessages
}

func NewAIChatMessageHandler(
	lacm *usecase.ListAIChatMessages,
) *AIChatMessageHandler {
	return &AIChatMessageHandler{
		lacm: lacm,
	}
}

// @Summary List ai chat messages
// @Description List ai chat messages
// @Tags AI Chat
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param ai_chat_id path string true "AI Chat ID" format(uuid)
// @Param search query string false "Search"
// @Param page query int false "Page"
// @Param page_size query int false "Page size"
// @Success 200 {object} dto.ListAIChatMessagesResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /v1/ai-chats/{ai_chat_id}/messages [get]
func (h AIChatMessageHandler) List(c *fiber.Ctx) error {
	paginationIn := parsePaginationParams(c)

	aiChatID := uuid.MustParse(c.Params(pathParamAIChatID))

	in := usecase.ListAIChatMessagesInput{
		PaginationInput: paginationIn,
		AIChatMessageOptions: repo.AIChatMessageOptions{
			AIChatID: aiChatID,
		},
	}

	ctx := c.UserContext()
	res, err := h.lacm.Execute(ctx, in)
	if err != nil {
		return errs.New(err)
	}

	return c.JSON(res)
}
