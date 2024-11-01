package handler

import (
	"net/http"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/usecase"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	uc *usecase.SignInUseCase
}

func NewAuthHandler(
	uc *usecase.SignInUseCase,
) *AuthHandler {
	return &AuthHandler{
		uc: uc,
	}
}

func (h AuthHandler) SignIn(c echo.Context) error {
	token := c.Request().Header.Get("Authorization")
	if token == "" {
		return errs.New("missing authorization header")
	}

	in := usecase.SignInUseCaseInput{Token: token}

	ctx := c.Request().Context()
	user, err := h.uc.Execute(ctx, in)
	if err != nil {
		return errs.New(err)
	}

	return c.JSON(http.StatusOK, user)
}
