package handler

import (
	"net/http"

	"github.com/danielmesquitta/api-finance-manager/internal/app/restapi/dto"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/usecase"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/jwtutil"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type BudgetHandler struct {
	ub *usecase.UpsertBudgetUseCase
}

func NewBudgetHandler(
	ub *usecase.UpsertBudgetUseCase,
) *BudgetHandler {
	return &BudgetHandler{
		ub: ub,
	}
}

// @Summary Create or update budget
// @Description Create or update budget
// @Tags Budget
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.UpsertBudgetRequest true "Request body"
// @Success 204
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /v1/budgets [post]
func (h BudgetHandler) Upsert(c echo.Context) error {
	var body dto.UpsertBudgetRequest
	if err := c.Bind(&body); err != nil {
		return errs.New(err)
	}

	claims := c.Get("claims").(*jwtutil.UserClaims)
	body.UserID = uuid.Must(uuid.Parse(claims.Issuer))

	ctx := c.Request().Context()
	if err := h.ub.Execute(
		ctx,
		body.UpsertBudgetUseCaseInput,
	); err != nil {
		return errs.New(err)
	}

	return c.NoContent(http.StatusNoContent)
}
