package handler

import (
	"net/http"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/usecase"
	"github.com/labstack/echo/v4"
)

type TransactionHandler struct {
	sa *usecase.SyncTransactions
}

func NewTransactionHandler(
	sa *usecase.SyncTransactions,
) *TransactionHandler {
	return &TransactionHandler{
		sa: sa,
	}
}

// @Summary Sync transactions from open finance
// @Description Webhook to sync transactions from open finance
// @Tags Transaction
// @Security BasicAuth
// @Accept json
// @Produce json
// @Success 204
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /v1/admin/transactions/sync [post]
func (h TransactionHandler) Sync(c echo.Context) error {
	ctx := c.Request().Context()
	if err := h.sa.Execute(ctx); err != nil {
		return errs.New(err)
	}

	return c.NoContent(http.StatusNoContent)
}
