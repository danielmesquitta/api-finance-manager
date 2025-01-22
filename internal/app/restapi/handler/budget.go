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
	ub   *usecase.UpsertBudget
	gb   *usecase.GetBudget
	gbc  *usecase.GetBudgetCategory
	db   *usecase.DeleteBudget
	lbct *usecase.ListBudgetCategoryTransactions
}

func NewBudgetHandler(
	ub *usecase.UpsertBudget,
	gb *usecase.GetBudget,
	gbc *usecase.GetBudgetCategory,
	db *usecase.DeleteBudget,
	lbct *usecase.ListBudgetCategoryTransactions,
) *BudgetHandler {
	return &BudgetHandler{
		ub:   ub,
		gb:   gb,
		gbc:  gbc,
		db:   db,
		lbct: lbct,
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
	userID := uuid.MustParse(claims.Issuer)
	body.UserID = userID

	ctx := c.Request().Context()
	if err := h.ub.Execute(
		ctx,
		body.UpsertBudgetInput,
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
// @Param date query string true "Date" format(date)
// @Success 200 {object} dto.GetBudgetResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /v1/budgets [get]
func (h BudgetHandler) Get(c echo.Context) error {
	claims := getUserClaims(c)
	userID := uuid.MustParse(claims.Issuer)

	date, err := parseDateParam(c, queryParamDate)
	if err != nil {
		return errs.New(err)
	}

	ctx := c.Request().Context()
	out, err := h.gb.Execute(ctx, usecase.GetBudgetInput{
		UserID: userID,
		Date:   date,
	})
	if err != nil {
		return errs.New(err)
	}

	return c.JSON(http.StatusOK, dto.GetBudgetResponse{
		GetBudgetOutput: *out,
	})
}

// @Summary Get budget category
// @Description Get budget category
// @Tags Budget
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param category_id path string true "Category ID" format(uuid)
// @Param date query string true "Date" format(date)
// @Success 200 {object} dto.GetBudgetCategoryResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /v1/budgets/categories/{category_id} [get]
func (h BudgetHandler) GetCategory(c echo.Context) error {
	claims := getUserClaims(c)
	userID := uuid.MustParse(claims.Issuer)

	date := c.QueryParam(queryParamDate)
	categoryID := uuid.MustParse(c.Param(pathParamCategoryID))

	ctx := c.Request().Context()
	out, err := h.gbc.Execute(ctx, usecase.GetBudgetCategoryInput{
		UserID:     userID,
		Date:       date,
		CategoryID: categoryID,
	})
	if err != nil {
		return errs.New(err)
	}

	return c.JSON(http.StatusOK, dto.GetBudgetCategoryResponse{
		GetBudgetCategoryOutput: *out,
	})
}

// @Summary List budget category transactions
// @Description List budget category transactions
// @Tags Budget
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param category_id path string true "Category ID" format(uuid)
// @Param page query int false "Page"
// @Param page_size query int false "Page size"
// @Param date query string true "Date" format(date)
// @Success 200 {object} dto.ListTransactionsResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /v1/budgets/categories/{category_id}/transactions [get]
func (h *BudgetHandler) ListCategoryTransactions(c echo.Context) error {
	claims := getUserClaims(c)
	userID := uuid.MustParse(claims.Issuer)

	paginationIn := parsePaginationParams(c)

	date, err := parseDateParam(c, queryParamDate)
	if err != nil {
		return errs.New(err)
	}

	categoryID := uuid.MustParse(c.Param(pathParamCategoryID))

	in := usecase.ListBudgetCategoryTransactionsInput{
		PaginationInput: paginationIn,
		Date:            date,
		UserID:          userID,
		CategoryID:      categoryID,
	}

	ctx := c.Request().Context()
	res, err := h.lbct.Execute(ctx, in)
	if err != nil {
		return errs.New(err)
	}

	return c.JSON(http.StatusOK, res)
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
// @Failure 404 {object} dto.ErrorResponse
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
