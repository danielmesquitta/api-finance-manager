package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h HealthHandler) Health(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]bool{"ok": true})
}
