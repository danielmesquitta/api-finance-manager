package handler

import (
	"net/http"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/usecase"
	"github.com/labstack/echo/v4"
)

type AccountHandler struct {
	sa *usecase.SyncAccounts
}

func NewAccountHandler(
	sa *usecase.SyncAccounts,
) *AccountHandler {
	return &AccountHandler{
		sa: sa,
	}
}

// @Summary Sync accounts from open finance
// @Description Webhook to sync user accounts from open finance
// @Tags Account
// @Security BasicAuth
// @Accept json
// @Produce json
// @Param request body dto.SyncAccountsRequest true "Request body"
// @Success 204
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /v1/admin/accounts/sync [post]
func (h AccountHandler) Sync(c echo.Context) error {
	in := usecase.SyncAccountsInput{}
	if err := c.Bind(&in); err != nil {
		return errs.New(err)
	}

	ctx := c.Request().Context()
	if err := h.sa.Execute(ctx, in); err != nil {
		return errs.New(err)
	}

	return c.NoContent(http.StatusNoContent)
}
