package handler

import (
	"net/http"

	"github.com/danielmesquitta/api-finance-manager/internal/app/restapi/dto"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/usecase"
	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	uc *usecase.CreateUserUseCase
}

func NewUserHandler(
	uc *usecase.CreateUserUseCase,
) *UserHandler {
	return &UserHandler{
		uc: uc,
	}
}

func (h UserHandler) Create(c echo.Context) error {
	reqDTO := dto.CreateUserRequestDTO{}
	if err := c.Bind(&reqDTO); err != nil {
		return errs.New(err)
	}

	in := usecase.CreateUserUseCaseInput{}
	if err := copier.Copy(&in, reqDTO); err != nil {
		return errs.New(err)
	}

	ctx := c.Request().Context()
	user, err := h.uc.Execute(ctx, in)
	if err != nil {
		return errs.New(err)
	}

	return c.JSON(http.StatusOK, user)
}
