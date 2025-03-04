package handler

import (
	"net/http"

	"github.com/danielmesquitta/api-finance-manager/internal/app/restapi/dto"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type FeedbackHandler struct {
	cf *usecase.CreateFeedback
}

func NewFeedbackHandler(
	cf *usecase.CreateFeedback,
) *FeedbackHandler {
	return &FeedbackHandler{
		cf: cf,
	}
}

// @Summary Create feedback
// @Description Create feedback
// @Tags Feedback
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.CreateFeedbackRequest true "Request body"
// @Success 204
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /v1/feedbacks [post]
func (h *FeedbackHandler) Create(c *fiber.Ctx) error {
	var body dto.CreateFeedbackRequest
	if err := c.BodyParser(&body); err != nil {
		return errs.New(err)
	}

	claims := GetUserClaims(c)
	userID := uuid.MustParse(claims.Issuer)
	body.UserID = userID

	ctx := c.UserContext()
	if err := h.cf.Execute(
		ctx,
		body.CreateFeedbackInput,
	); err != nil {
		return errs.New(err)
	}

	return c.SendStatus(http.StatusNoContent)
}
