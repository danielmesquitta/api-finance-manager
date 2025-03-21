package handler

import (
	"github.com/danielmesquitta/api-finance-manager/internal/app/server/dto"
	"github.com/gofiber/fiber/v2"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// @Summary Health check
// @Description Health check
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} dto.HealthResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /health [post]
func (h *HealthHandler) Health(c *fiber.Ctx) error {
	return c.JSON(dto.HealthResponse{Status: "ok"})
}
