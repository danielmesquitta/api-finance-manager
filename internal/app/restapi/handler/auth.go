package handler

import (
	"net/http"

	"github.com/danielmesquitta/api-finance-manager/internal/app/restapi/dto"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/usecase"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/jwtutil"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	si *usecase.SignInUseCase
	rt *usecase.RefreshTokenUseCase
}

func NewAuthHandler(
	si *usecase.SignInUseCase,
	rt *usecase.RefreshTokenUseCase,
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
// @Param request body dto.SignInRequestDTO true "Request body"
// @Success 200 {object} dto.SignInResponseDTO
// @Failure 400 {object} dto.ErrorResponseDTO
// @Failure 401 {object} dto.ErrorResponseDTO
// @Failure 500 {object} dto.ErrorResponseDTO
// @Router /v1/auth/sign-in [post]
func (h AuthHandler) SignIn(c echo.Context) error {
	token := c.Request().Header.Get("Authorization")
	if token == "" {
		return errs.ErrUnauthorized
	}

	var body dto.SignInRequestDTO
	if err := c.Bind(&body); err != nil {
		return errs.New(err)
	}

	ctx := c.Request().Context()
	out, err := h.si.Execute(
		ctx,
		usecase.SignInUseCaseInput{Token: token, Provider: body.Provider},
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
// @Success 200 {object} dto.SignInResponseDTO
// @Failure 400 {object} dto.ErrorResponseDTO
// @Failure 401 {object} dto.ErrorResponseDTO
// @Failure 500 {object} dto.ErrorResponseDTO
// @Router /v1/auth/refresh [post]
func (h AuthHandler) RefreshToken(c echo.Context) error {
	claims := c.Get("claims").(*jwtutil.UserClaims)

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
