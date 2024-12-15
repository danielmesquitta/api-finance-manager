package handler

import (
	"net/http"

	"github.com/danielmesquitta/api-finance-manager/internal/app/restapi/dto"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/usecase"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type BudgetHandler struct {
	ub *usecase.UpsertBudgetUseCase
	gb *usecase.GetBudgetUseCase
	db *usecase.DeleteBudgetUseCase
}

func NewBudgetHandler(
	ub *usecase.UpsertBudgetUseCase,
	gb *usecase.GetBudgetUseCase,
	db *usecase.DeleteBudgetUseCase,
) *BudgetHandler {
	return &BudgetHandler{
		ub: ub,
		gb: gb,
		db: db,
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

	claims := getUserClaims(c)
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

// @Summary Get budget
// @Description Get budget
// @Tags Budget
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param month query int false "Month" Format(1-12)
// @Success 200 {object} dto.GetBudgetResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /v1/budgets [get]
func (h BudgetHandler) Get(c echo.Context) error {
	claims := getUserClaims(c)
	userID := uuid.Must(uuid.Parse(claims.Issuer))

	ctx := c.Request().Context()
	out, err := h.gb.Execute(ctx, usecase.GetBudgetUseCaseInput{
		UserID: userID,
	})
	if err != nil {
		return errs.New(err)
	}

	return c.JSON(http.StatusOK, dto.GetBudgetResponse{
		GetBudgetUseCaseOutput: *out,
	})
}

// @Summary Delete budget
// @Description Delete budget
// @Tags Budget
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 204
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /v1/budgets [delete]
func (h BudgetHandler) Delete(c echo.Context) error {
	claims := getUserClaims(c)
	userID := uuid.Must(uuid.Parse(claims.Issuer))

	ctx := c.Request().Context()
	err := h.db.Execute(ctx, userID)
	if err != nil {
		return errs.New(err)
	}

	return c.NoContent(http.StatusNoContent)
}
