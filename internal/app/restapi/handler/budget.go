package handler

import (
	"net/http"

	"github.com/danielmesquitta/api-finance-manager/internal/app/restapi/dto"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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
func (h BudgetHandler) Upsert(c *fiber.Ctx) error {
	var in dto.UpsertBudgetRequest
	if err := c.BodyParser(&in); err != nil {
		return errs.New(err)
	}

	claims := GetUserClaims(c)
	userID := uuid.MustParse(claims.Issuer)
	in.UserID = userID

	ctx := c.UserContext()
	if err := h.ub.Execute(
		ctx,
		in.UpsertBudgetInput,
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
	claims := GetUserClaims(c)
	userID := uuid.MustParse(claims.Issuer)

	date, err := parseDateParam(c, QueryParamDate)
	if err != nil {
		return errs.New(err)
	}

	ctx := c.UserContext()
	out, err := h.gb.Execute(ctx, usecase.GetBudgetInput{
		UserID: userID,
		Date:   date,
	})
	if err != nil {
		return errs.New(err)
	}

	return c.JSON(dto.GetBudgetResponse{
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
func (h BudgetHandler) GetTransactionCategory(c *fiber.Ctx) error {
	claims := GetUserClaims(c)
	userID := uuid.MustParse(claims.Issuer)

	date := c.Query(QueryParamDate)
	categoryID := uuid.MustParse(c.Params(pathParamCategoryID))

	ctx := c.UserContext()
	out, err := h.gbc.Execute(ctx, usecase.GetBudgetCategoryInput{
		UserID:     userID,
		Date:       date,
		CategoryID: categoryID,
	})
	if err != nil {
		return errs.New(err)
	}

	return c.JSON(dto.GetBudgetCategoryResponse{
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
func (h *BudgetHandler) ListCategoryTransactions(c *fiber.Ctx) error {
	claims := GetUserClaims(c)
	userID := uuid.MustParse(claims.Issuer)

	paginationIn := parsePaginationParams(c)

	date, err := parseDateParam(c, QueryParamDate)
	if err != nil {
		return errs.New(err)
	}

	categoryID := uuid.MustParse(c.Params(pathParamCategoryID))

	in := usecase.ListBudgetCategoryTransactionsInput{
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
	claims := GetUserClaims(c)
	userID := uuid.Must(uuid.Parse(claims.Issuer))

	ctx := c.UserContext()
	err := h.db.Execute(ctx, userID)
	if err != nil {
		return errs.New(err)
	}

	return c.SendStatus(http.StatusNoContent)
}
