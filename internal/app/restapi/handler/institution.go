package handler

import (
	"net/http"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/usecase"
	"github.com/labstack/echo/v4"
)

type InstitutionHandler struct {
	si *usecase.SyncInstitutions
}

func NewInstitutionHandler(
	si *usecase.SyncInstitutions,
) *InstitutionHandler {
	return &InstitutionHandler{
		si: si,
	}
}

// @Summary Sync institutions from open finance
// @Description Sync institutions from open finance
// @Tags Institution
// @Security BasicAuth
// @Accept json
// @Produce json
// @Success 204
// @Failure 500 {object} dto.ErrorResponse
// @Router /v1/admin/institutions/sync [post]
func (h InstitutionHandler) Sync(c echo.Context) error {
	if err := h.si.Execute(c.Request().Context()); err != nil {
		return errs.New(err)
	}

	return c.NoContent(http.StatusNoContent)
}
