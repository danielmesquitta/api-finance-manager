package middleware

import (
	"errors"
	"net/http"

	"github.com/danielmesquitta/api-finance-manager/internal/app/restapi/dto"
	"github.com/danielmesquitta/api-finance-manager/internal/app/restapi/handler"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/gofiber/fiber/v2"
)

var mapAppErrToHTTPError = map[errs.Code]int{
	errs.ErrCodeForbidden:    http.StatusForbidden,
	errs.ErrCodeUnauthorized: http.StatusUnauthorized,
	errs.ErrCodeValidation:   http.StatusBadRequest,
	errs.ErrCodeUnknown:      http.StatusInternalServerError,
	errs.ErrCodeNotFound:     http.StatusNotFound,
}

func (m *Middleware) ErrorHandler(ctx *fiber.Ctx, err error) error {
	var appErr *errs.Err
	if errors.As(err, &appErr) {
		code, ok := mapAppErrToHTTPError[appErr.Code]
		if !ok {
			code = http.StatusInternalServerError
		}

		if code >= 500 {
			return m.handleInternalServerError(ctx, appErr)
		}

		return ctx.Status(code).JSON(
			dto.ErrorResponse{Message: appErr.Message},
		)
	}

	var fiberErr *fiber.Error
	if errors.As(err, &fiberErr) {
		code := fiberErr.Code

		if code >= 500 {
			return m.handleInternalServerError(ctx, errs.New(fiberErr))
		}

		return ctx.Status(code).JSON(
			dto.ErrorResponse{Message: fiberErr.Message},
		)
	}

	return nil
}

func (m *Middleware) handleInternalServerError(
	c *fiber.Ctx,
	appErr *errs.Err,
) error {
	statusCode := mapAppErrToHTTPError[appErr.Code]

	args := []any{
		"method", c.Method(),
		"url", c.BaseURL(),
	}

	queries := c.Queries()
	if len(queries) > 0 {
		args = append(args, "query", queries)
	}

	requestId := c.Get(fiber.HeaderXRequestID)
	if requestId != "" {
		args = append(args, "request_id", requestId)
	}

	userId := ""
	claims := handler.GetUserClaims(c)
	if claims != nil {
		userId = claims.Issuer
	}
	if userId != "" {
		args = append(args, "user_id", userId)
	}

	requestData := map[string]any{}
	_ = c.BodyParser(&requestData)
	if len(requestData) > 0 {
		args = append(args, "body", requestData)
	}

	args = append(args, "stacktrace", appErr.StackTrace)

	m.l.Error(
		appErr.Error(),
		args...,
	)

	return c.Status(statusCode).JSON(
		dto.ErrorResponse{Message: "internal server error"},
	)
}
