package handler

import (
	"net/http"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/usecase"
	"github.com/labstack/echo/v4"
)

type CategoryHandler struct {
	sc *usecase.SyncCategoriesUseCase
	lc *usecase.ListCategoriesUseCase
}

func NewCategoryHandler(
	sc *usecase.SyncCategoriesUseCase,
	lc *usecase.ListCategoriesUseCase,
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
// @Router /v1/admin/categories/sync [post]
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
// @Failure 500 {object} dto.ErrorResponse
// @Router /v1/categories [get]
func (h CategoryHandler) List(c echo.Context) error {
	search, page, pageSize := getPaginationParams(c)
	in := usecase.ListCategoriesUseCaseInput{
		PaginationInput: usecase.PaginationInput{
			Search:   search,
			Page:     page,
			PageSize: pageSize,
		},
	}

	res, err := h.lc.Execute(c.Request().Context(), in)
	if err != nil {
		return errs.New(err)
	}

	return c.JSON(http.StatusOK, res)
}
