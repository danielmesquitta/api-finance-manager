package handler

import (
	"net/http"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/usecase"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
	"github.com/labstack/echo/v4"
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
func (h CategoryHandler) Sync(c echo.Context) error {
	if err := h.sc.Execute(c.Request().Context()); err != nil {
		return errs.New(err)
	}

	return c.NoContent(http.StatusNoContent)
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
func (h CategoryHandler) List(c echo.Context) error {
	search := c.QueryParam(queryParamSearch)
	paginationIn := parsePaginationParams(c)

	in := usecase.ListCategoriesInput{
		PaginationInput: paginationIn,
		TransactionCategoryOptions: repo.TransactionCategoryOptions{
			Search: search,
		},
	}

	ctx := c.Request().Context()
	res, err := h.lc.Execute(ctx, in)
	if err != nil {
		return errs.New(err)
	}

	return c.JSON(http.StatusOK, res)
}
