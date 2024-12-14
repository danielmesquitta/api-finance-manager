package handler

import (
	"net/http"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/usecase"
	"github.com/labstack/echo/v4"
)

type CategoryHandler struct {
	sc *usecase.SyncCategoriesUseCase
}

func NewCategoryHandler(
	sc *usecase.SyncCategoriesUseCase,
) *CategoryHandler {
	return &CategoryHandler{
		sc: sc,
	}
}

// @Summary Sync categories from open finance
// @Description Sync categories from open finance
// @Tags Category
// @Security BasicAuth
// @Accept json
// @Produce json
// @Success 204
// @Failure 500 {object} dto.ErrorResponse
// @Router /v1/admin/categories/sync [post]
func (h CategoryHandler) Sync(c echo.Context) error {
	if err := h.sc.Execute(c.Request().Context()); err != nil {
		return errs.New(err)
	}

	return c.NoContent(http.StatusNoContent)
}
