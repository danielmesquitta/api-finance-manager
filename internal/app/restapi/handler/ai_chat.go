package handler

import (
	"net/http"

	"github.com/danielmesquitta/api-finance-manager/internal/app/restapi/dto"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/usecase"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
)

type AIChatHandler struct {
	cac   *usecase.CreateAIChat
	dac   *usecase.DeleteAIChat
	uac   *usecase.UpdateAIChat
	lac   *usecase.ListAIChats
	lacmr *usecase.ListAIChatMessagesAndAnswers
}

func NewAIChatHandler(
	cac *usecase.CreateAIChat,
	dac *usecase.DeleteAIChat,
	uac *usecase.UpdateAIChat,
	lac *usecase.ListAIChats,
) *AIChatHandler {
	return &AIChatHandler{
		cac: cac,
		dac: dac,
		uac: uac,
		lac: lac,
	}
}

// @Summary Create ai chat
// @Description Create ai chat
// @Tags AI Chat
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 201 {object} dto.CreateAIChatResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /v1/ai-chats [post]
func (h *AIChatHandler) Create(c *fiber.Ctx) error {
	userID, tier, err := GetUser(c)
	if err != nil {
		return errs.New(err)
	}

	in := usecase.CreateAIChatInput{
		UserID: userID,
		Tier:   tier,
	}

	ctx := c.UserContext()
	out, err := h.cac.Execute(ctx, in)
	if err != nil {
		return errs.New(err)
	}

	return c.Status(http.StatusCreated).JSON(out)
}

// @Summary Delete ai chat
// @Description Delete ai chat
// @Tags AI Chat
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param ai_chat_id path string true "AI Chat ID" format(uuid)
// @Success 204
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /v1/ai-chats/{ai_chat_id} [delete]
func (h *AIChatHandler) Delete(c *fiber.Ctx) error {
	aiChatID, err := parseUUIDPathParam(c, pathParamAIChatID)
	if err != nil {
		return errs.New(err)
	}

	ctx := c.UserContext()
	if err := h.dac.Execute(ctx, aiChatID); err != nil {
		return errs.New(err)
	}

	return c.SendStatus(http.StatusNoContent)
}

// @Summary Update ai chat
// @Description Update ai chat
// @Tags AI Chat
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param ai_chat_id path string true "AI Chat ID" format(uuid)
// @Param request body dto.UpdateAIChatRequest true "Request body"
// @Success 204
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /v1/ai-chats/{ai_chat_id} [put]
func (h *AIChatHandler) Update(c *fiber.Ctx) error {
	var body dto.UpdateAIChatRequest
	if err := c.BodyParser(&body); err != nil {
		return errs.New(err)
	}

	aiChatID, err := parseUUIDPathParam(c, pathParamAIChatID)
	if err != nil {
		return errs.New(err)
	}

	_, tier, err := GetUser(c)
	if err != nil {
		return errs.New(err)
	}

	var in usecase.UpdateAIChatInput
	if err := copier.Copy(&in, body); err != nil {
		return errs.New(err)
	}
	in.ID = aiChatID
	in.Tier = tier

	ctx := c.UserContext()
	if err := h.uac.Execute(ctx, in); err != nil {
		return errs.New(err)
	}

	return c.SendStatus(http.StatusNoContent)
}

// @Summary List ai chats
// @Description List ai chats
// @Tags AI Chat
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param search query string false "Search"
// @Param page query int false "Page"
// @Param page_size query int false "Page size"
// @Success 200 {object} dto.ListAIChatsResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /v1/ai-chats [get]
func (h AIChatHandler) List(c *fiber.Ctx) error {
	search := c.Query(QueryParamSearch)
	paginationIn := parsePaginationParams(c)

	userID, _, err := GetUser(c)
	if err != nil {
		return errs.New(err)
	}

	in := usecase.ListAIChatsInput{
		PaginationInput: paginationIn,
		AIChatOptions: repo.AIChatOptions{
			Search: search,
			UserID: userID,
		},
	}

	ctx := c.UserContext()
	res, err := h.lac.Execute(ctx, in)
	if err != nil {
		return errs.New(err)
	}

	return c.JSON(res)
}

// @Summary List ai chats messages
// @Description List ai chats messages
// @Tags AI Chat
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param ai_chat_id path string true "AI Chat ID" format(uuid)
// @Param page query int false "Page"
// @Param page_size query int false "Page size"
// @Success 200 {object} dto.ListAIChatMessagesAndAnswersResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /v1/ai-chats/{ai_chat_id}/messages [get]
func (h AIChatHandler) ListMessages(c *fiber.Ctx) error {
	aiChatID, err := parseUUIDPathParam(c, pathParamAIChatID)
	if err != nil {
		return errs.New(err)
	}

	paginationIn := parsePaginationParams(c)

	in := usecase.ListAIChatMessagesAndAnswersInput{
		PaginationInput: paginationIn,
		AIChatID:        aiChatID,
	}

	ctx := c.UserContext()
	res, err := h.lacmr.Execute(ctx, in)
	if err != nil {
		return errs.New(err)
	}

	return c.JSON(res)
}
