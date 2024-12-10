package handler

import (
	"net/http"

	"github.com/danielmesquitta/api-finance-manager/internal/app/restapi/dto"
	"github.com/labstack/echo/v4"
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
// @Success 200 {object} dto.HealthResponseDTO
// @Failure 500 {object} dto.ErrorResponseDTO
// @Router /health [get]
func (h HealthHandler) Health(c echo.Context) error {
	return c.JSON(http.StatusOK, dto.HealthResponseDTO{Ok: true})
}
