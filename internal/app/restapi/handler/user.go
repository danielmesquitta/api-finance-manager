package handler

import (
	"net/http"

	"github.com/danielmesquitta/api-finance-manager/internal/app/restapi/dto"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/usecase"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	gu *usecase.GetUser
}

func NewUserHandler(
	gu *usecase.GetUser,
) *UserHandler {
	return &UserHandler{
		gu: gu,
	}
}

// @Summary Get logged-in user profile
// @Description Get logged-in user profile
// @Tags User
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} dto.UserProfileResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /v1/users/profile [get]
func (h UserHandler) Profile(c echo.Context) error {
	claims := getUserClaims(c)
	userID := uuid.MustParse(claims.Issuer)

	user, err := h.gu.Execute(c.Request().Context(), userID)
	if err != nil {
		return errs.New(err)
	}

	return c.JSON(http.StatusOK, dto.UserProfileResponse{
		User: *user,
	})
}
