package handler

import (
	"net/http"

	"github.com/danielmesquitta/api-finance-manager/internal/app/restapi/dto"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/usecase"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	si *usecase.SignIn
	rt *usecase.RefreshToken
}

func NewAuthHandler(
	si *usecase.SignIn,
	rt *usecase.RefreshToken,
) *AuthHandler {
	return &AuthHandler{
		si: si,
		rt: rt,
	}
}

// @Summary Sign in
// @Description Authenticate user through Google or Apple token
// @Tags Auth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.SignInRequest true "Request body"
// @Success 200 {object} dto.SignInResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /v1/auth/sign-in [post]
func (h AuthHandler) SignIn(c echo.Context) error {
	token := c.Request().Header.Get("Authorization")
	if token == "" {
		return errs.ErrUnauthorized
	}

	var body dto.SignInRequest
	if err := c.Bind(&body); err != nil {
		return errs.New(err)
	}

	ctx := c.Request().Context()
	out, err := h.si.Execute(
		ctx,
		usecase.SignInInput{Token: token, Provider: body.Provider},
	)
	if err != nil {
		return errs.New(err)
	}

	return c.JSON(http.StatusOK, out)
}

// @Summary Refresh token
// @Description Use refresh token to generate new access token
// @Tags Auth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} dto.SignInResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /v1/auth/refresh [post]
func (h AuthHandler) RefreshToken(c echo.Context) error {
	claims := getUserClaims(c)

	ctx := c.Request().Context()
	out, err := h.rt.Execute(
		ctx,
		uuid.MustParse(claims.Issuer),
	)
	if err != nil {
		return errs.New(err)
	}

	return c.JSON(http.StatusOK, out)
}
