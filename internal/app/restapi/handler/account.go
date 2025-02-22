package handler

import (
	"net/http"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/usecase"
	"github.com/gofiber/fiber/v2"
)

type AccountHandler struct {
	ca *usecase.CreateAccounts
}

func NewAccountHandler(
	ca *usecase.CreateAccounts,
) *AccountHandler {
	return &AccountHandler{
		ca: ca,
	}
}

// @Summary Sync accounts from open finance
// @Description Webhook to sync user accounts from open finance
// @Tags Account
// @Security BasicAuth
// @Accept json
// @Produce json
// @Param request body dto.CreateAccountsRequest true "Request body"
// @Success 204
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /v1/admin/accounts [post]
func (h AccountHandler) Create(c *fiber.Ctx) error {
	in := usecase.CreateAccountsInput{}
	if err := c.BodyParser(&in); err != nil {
		return errs.New(err)
	}

	ctx := c.UserContext()
	if err := h.ca.Execute(ctx, in); err != nil {
		return errs.New(err)
	}

	return c.SendStatus(http.StatusNoContent)
}
