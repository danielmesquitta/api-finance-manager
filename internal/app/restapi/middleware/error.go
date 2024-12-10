package middleware

import (
	"log/slog"
	"net/http"

	"github.com/danielmesquitta/api-finance-manager/internal/app/restapi/dto"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/labstack/echo/v4"
)

var mapErrTypeToStatusCode = map[errs.Type]int{
	errs.ErrTypeForbidden:    http.StatusForbidden,
	errs.ErrTypeUnauthorized: http.StatusUnauthorized,
	errs.ErrTypeValidation:   http.StatusBadRequest,
	errs.ErrTypeUnknown:      http.StatusInternalServerError,
	errs.ErrTypeNotFound:     http.StatusNotFound,
}

func (m *Middleware) ErrorHandler(
	defaultErrorHandler echo.HTTPErrorHandler,
) echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		if c.Response().Committed {
			return
		}

		if appErr, ok := err.(*errs.Err); ok {
			statusCode := mapErrTypeToStatusCode[appErr.Type]
			isInternalServerError := statusCode >= 500 || statusCode == 0
			if isInternalServerError {
				if err := m.handleInternalServerError(c, appErr); err != nil {
					slog.Error(
						"failed to handle internal server error",
						"err",
						err,
					)
				}
				return
			}

			if err := c.JSON(
				statusCode,
				dto.ErrorResponseDTO{Message: appErr.Error()},
			); err != nil {
				slog.Error("failed to handle app err", "err", err)
			}
			return
		}

		defaultErrorHandler(err, c)
	}
}

func (m *Middleware) handleInternalServerError(
	c echo.Context,
	appErr *errs.Err,
) error {
	statusCode := mapErrTypeToStatusCode[appErr.Type]
	req := c.Request()

	requestData := map[string]any{}
	_ = c.Bind(&requestData)

	slog.Error(
		appErr.Error(),
		"url", req.URL.Path,
		"body", requestData,
		"query", c.QueryParams(),
		"params", c.ParamValues(),
		"stacktrace", appErr.StackTrace,
	)

	return c.JSON(
		statusCode,
		dto.ErrorResponseDTO{Message: "internal server error"},
	)
}
