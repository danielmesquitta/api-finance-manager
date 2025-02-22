package handler

import (
	"net/http"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/usecase"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
	"github.com/gofiber/fiber/v2"
)

type CategoryHandler struct {
	sc *usecase.SyncCategories
	lc *usecase.ListTransactionCategories
}

func NewCategoryHandler(
	sc *usecase.SyncCategories,
	lc *usecase.ListTransactionCategories,
) *CategoryHandler {
	return &CategoryHandler{
		sc: sc,
		lc: lc,
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
// @Router /v1/admin/transactions/categories/sync [post]
func (h CategoryHandler) Sync(c *fiber.Ctx) error {
	if err := h.sc.Execute(c.UserContext()); err != nil {
		return errs.New(err)
	}

	return c.SendStatus(http.StatusNoContent)
}

// @Summary List categories
// @Description List categories
// @Tags Category
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param search query string false "Search"
// @Param page query int false "Page"
// @Param page_size query int false "Page size"
// @Success 200 {object} dto.ListCategoriesResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /v1/transactions/categories [get]
func (h CategoryHandler) List(c *fiber.Ctx) error {
	search := c.Query(queryParamSearch)
	paginationIn := parsePaginationParams(c)

	in := usecase.ListCategoriesInput{
		PaginationInput: paginationIn,
		TransactionCategoryOptions: repo.TransactionCategoryOptions{
			Search: search,
		},
	}

	ctx := c.UserContext()
	res, err := h.lc.Execute(ctx, in)
	if err != nil {
		return errs.New(err)
	}

	return c.JSON(res)
}
