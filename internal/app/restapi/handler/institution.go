package handler

import (
	"net/http"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/usecase"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type InstitutionHandler struct {
	si *usecase.SyncInstitutions
	li *usecase.ListInstitutions
}

func NewInstitutionHandler(
	si *usecase.SyncInstitutions,
	li *usecase.ListInstitutions,
) *InstitutionHandler {
	return &InstitutionHandler{
		si: si,
		li: li,
	}
}

// @Summary Sync institutions from open finance
// @Description Sync institutions from open finance
// @Tags Institution
// @Security BasicAuth
// @Accept json
// @Produce json
// @Success 204
// @Failure 500 {object} dto.ErrorResponse
// @Router /v1/admin/institutions/sync [post]
func (h InstitutionHandler) Sync(c echo.Context) error {
	if err := h.si.Execute(c.Request().Context()); err != nil {
		return errs.New(err)
	}

	return c.NoContent(http.StatusNoContent)
}

// @Summary List institutions
// @Description List all institutions
// @Tags Institution
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param search query string false "Search"
// @Param page query int false "Page"
// @Param page_size query int false "Page size"
// @Success 200 {object} dto.ListInstitutionsResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /v1/institutions [get]
func (h InstitutionHandler) List(c echo.Context) error {
	search := c.QueryParam(queryParamSearch)
	paginationIn := parsePaginationParams(c)

	in := usecase.ListInstitutionsInput{
		PaginationInput: paginationIn,
		InstitutionOptions: repo.InstitutionOptions{
			Search: search,
		},
	}

	ctx := c.Request().Context()
	res, err := h.li.Execute(ctx, in)
	if err != nil {
		return errs.New(err)
	}

	return c.JSON(http.StatusOK, res)
}

// @Summary List user institutions
// @Description List user institutions
// @Tags Institution
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param search query string false "Search"
// @Param page query int false "Page"
// @Param page_size query int false "Page size"
// @Success 200 {object} dto.ListInstitutionsResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /v1/users/institutions [get]
func (h InstitutionHandler) ListUserInstitutions(c echo.Context) error {
	search := c.QueryParam(queryParamSearch)
	paginationIn := parsePaginationParams(c)
	claims := getUserClaims(c)
	userID := uuid.MustParse(claims.Issuer)

	in := usecase.ListInstitutionsInput{
		PaginationInput: paginationIn,
		InstitutionOptions: repo.InstitutionOptions{
			UserID: userID,
			Search: search,
		},
	}

	ctx := c.Request().Context()
	res, err := h.li.Execute(ctx, in)
	if err != nil {
		return errs.New(err)
	}

	return c.JSON(http.StatusOK, res)
}
