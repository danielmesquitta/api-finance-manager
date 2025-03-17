package handler

import (
	"github.com/danielmesquitta/api-finance-manager/internal/app/restapi/dto"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/usecase"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	gu *usecase.GetUser
	uu *usecase.UpdateUser
	du *usecase.DeleteUser
}

func NewUserHandler(
	gu *usecase.GetUser,
	uu *usecase.UpdateUser,
	du *usecase.DeleteUser,
) *UserHandler {
	return &UserHandler{
		gu: gu,
		uu: uu,
		du: du,
	}
}

// @Summary Get logged-in user profile
// @Description Get logged-in user profile
// @Tags User
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} dto.GetUserProfileResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /v1/users/profile [get]
func (h UserHandler) GetProfile(c *fiber.Ctx) error {
	userID, _, err := GetUser(c)
	if err != nil {
		return errs.New(err)
	}

	ctx := c.UserContext()
	user, err := h.gu.Execute(ctx, userID)
	if err != nil {
		return errs.New(err)
	}

	return c.JSON(dto.GetUserProfileResponse{
		User: *user,
	})
}

// @Summary Update logged-in user
// @Description Update logged-in user
// @Tags User
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.UpdateProfileRequest true "Request body"
// @Success 204
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /v1/users/profile [put]
func (h UserHandler) UpdateProfile(c *fiber.Ctx) error {
	var in usecase.UpdateUserInput
	if err := c.BodyParser(&in); err != nil {
		return errs.New(err)
	}

	userID, _, err := GetUser(c)
	if err != nil {
		return errs.New(err)
	}
	in.ID = userID

	ctx := c.UserContext()
	if err := h.uu.Execute(ctx, in); err != nil {
		return errs.New(err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// @Summary Delete logged-in user
// @Description Delete logged-in user
// @Tags User
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 204
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /v1/users/profile [delete]
func (h UserHandler) DeleteProfile(c *fiber.Ctx) error {
	userID, _, err := GetUser(c)
	if err != nil {
		return errs.New(err)
	}

	ctx := c.UserContext()
	if err := h.du.Execute(ctx, userID); err != nil {
		return errs.New(err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}
