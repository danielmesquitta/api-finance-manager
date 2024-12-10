package handler

import (
	"net/http"

	"github.com/danielmesquitta/api-finance-manager/internal/app/restapi/dto"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/usecase"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	signInUseCase *usecase.SignInUseCase
}

func NewAuthHandler(
	signInUseCase *usecase.SignInUseCase,
) *AuthHandler {
	return &AuthHandler{
		signInUseCase: signInUseCase,
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
// @Router /auth/sign-in [post]
func (h AuthHandler) SignIn(c echo.Context) error {
	token := c.Request().Header.Get("Authorization")
	if token == "" {
		return errs.New("missing authorization header")
	}

	var body dto.SignInRequestDTO
	if err := c.Bind(&body); err != nil {
		return errs.New(err)
	}

	in := usecase.SignInUseCaseInput{Token: token, Provider: body.Provider}

	ctx := c.Request().Context()
	out, err := h.signInUseCase.Execute(ctx, in)
	if err != nil {
		return errs.New(err)
	}

	return c.JSON(http.StatusOK, out)
}
