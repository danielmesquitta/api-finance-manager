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
// @Success 200 {object} dto.HealthResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /health [get]
func (h HealthHandler) Health(c echo.Context) error {
	return c.JSON(http.StatusOK, dto.HealthResponse{Ok: true})
}
