package handler

import (
	"net/http"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/usecase"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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
func (h InstitutionHandler) Sync(c *fiber.Ctx) error {
	if err := h.si.Execute(c.UserContext()); err != nil {
		return errs.New(err)
	}

	return c.SendStatus(http.StatusNoContent)
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
func (h InstitutionHandler) List(c *fiber.Ctx) error {
	search := c.Query(queryParamSearch)
	paginationIn := parsePaginationParams(c)

	in := usecase.ListInstitutionsInput{
		PaginationInput: paginationIn,
		InstitutionOptions: repo.InstitutionOptions{
			Search: search,
		},
	}

	ctx := c.UserContext()
	res, err := h.li.Execute(ctx, in)
	if err != nil {
		return errs.New(err)
	}

	return c.JSON(res)
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
func (h InstitutionHandler) ListUserInstitutions(c *fiber.Ctx) error {
	search := c.Query(queryParamSearch)
	paginationIn := parsePaginationParams(c)
	claims := GetUserClaims(c)
	userID := uuid.MustParse(claims.Issuer)

	in := usecase.ListInstitutionsInput{
		PaginationInput: paginationIn,
		InstitutionOptions: repo.InstitutionOptions{
			UserID: userID,
			Search: search,
		},
	}

	ctx := c.UserContext()
	res, err := h.li.Execute(ctx, in)
	if err != nil {
		return errs.New(err)
	}

	return c.JSON(res)
}
