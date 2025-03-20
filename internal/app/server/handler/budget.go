package handler

import (
	"net/http"

	"github.com/danielmesquitta/api-finance-manager/internal/app/server/dto"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/usecase/budget"
	"github.com/gofiber/fiber/v2"
)

type BudgetHandler struct {
	ub   *budget.UpsertBudgetUseCase
	gb   *budget.GetBudgetUseCase
	gbc  *budget.GetBudgetCategoryUseCase
	db   *budget.DeleteBudgetUseCase
	lbct *budget.ListBudgetCategoryTransactionsUseCase
}

func NewBudgetHandler(
	ub *budget.UpsertBudgetUseCase,
	gb *budget.GetBudgetUseCase,
	gbc *budget.GetBudgetCategoryUseCase,
	db *budget.DeleteBudgetUseCase,
	lbct *budget.ListBudgetCategoryTransactionsUseCase,
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
func (h BudgetHandler) Upsert(c *fiber.Ctx) error {
	var in dto.UpsertBudgetRequest
	if err := c.BodyParser(&in); err != nil {
		return errs.New(err)
	}

	userID, _, err := GetUser(c)
	if err != nil {
		return errs.New(err)
	}
	in.UserID = userID

	ctx := c.UserContext()
	if err := h.ub.Execute(
		ctx,
		in.UpsertBudgetUseCaseInput,
	); err != nil {
		return errs.New(err)
	}

	return c.SendStatus(http.StatusNoContent)
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
func (h BudgetHandler) Get(c *fiber.Ctx) error {
	userID, _, err := GetUser(c)
	if err != nil {
		return errs.New(err)
	}

	date, err := parseDateQueryParam(c, QueryParamDate)
	if err != nil {
		return errs.New(err)
	}

	ctx := c.UserContext()
	out, err := h.gb.Execute(ctx, budget.GetBudgetUseCaseInput{
		UserID: userID,
		Date:   date,
	})
	if err != nil {
		return errs.New(err)
	}

	return c.JSON(dto.GetBudgetResponse{
		GetBudgetUseCaseOutput: *out,
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
func (h BudgetHandler) GetTransactionCategoryByID(c *fiber.Ctx) error {
	categoryID, err := parseUUIDPathParam(c, pathParamCategoryID)
	if err != nil {
		return errs.New(err)
	}

	userID, _, err := GetUser(c)
	if err != nil {
		return errs.New(err)
	}

	date := c.Query(QueryParamDate)

	ctx := c.UserContext()
	out, err := h.gbc.Execute(ctx, budget.GetBudgetCategoryUseCaseInput{
		UserID:     userID,
		Date:       date,
		CategoryID: categoryID,
	})
	if err != nil {
		return errs.New(err)
	}

	return c.JSON(dto.GetBudgetCategoryResponse{
		GetBudgetCategoryUseCaseOutput: *out,
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
func (h *BudgetHandler) ListCategoryTransactions(c *fiber.Ctx) error {
	userID, _, err := GetUser(c)
	if err != nil {
		return errs.New(err)
	}

	paginationIn := parsePaginationParams(c)

	date, err := parseDateQueryParam(c, QueryParamDate)
	if err != nil {
		return errs.New(err)
	}

	categoryID, err := parseUUIDPathParam(c, pathParamCategoryID)
	if err != nil {
		return errs.New(err)
	}

	in := budget.ListBudgetCategoryTransactionsUseCaseInput{
		PaginationInput: paginationIn,
		Date:            date,
		UserID:          userID,
		CategoryID:      categoryID,
	}

	ctx := c.UserContext()
	res, err := h.lbct.Execute(ctx, in)
	if err != nil {
		return errs.New(err)
	}

	return c.JSON(res)
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
func (h BudgetHandler) Delete(c *fiber.Ctx) error {
	userID, _, err := GetUser(c)
	if err != nil {
		return errs.New(err)
	}

	ctx := c.UserContext()
	if err := h.db.Execute(ctx, userID); err != nil {
		return errs.New(err)
	}

	return c.SendStatus(http.StatusNoContent)
}
