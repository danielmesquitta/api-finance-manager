package handler

import (
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/usecase/paymentmethod"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
	"github.com/gofiber/fiber/v2"
)

type PaymentMethodHandler struct {
	lpm *paymentmethod.ListPaymentMethodsUseCase
}

func NewPaymentMethodHandler(
	lpm *paymentmethod.ListPaymentMethodsUseCase,
) *PaymentMethodHandler {
	return &PaymentMethodHandler{
		lpm: lpm,
	}
}

// @Summary List payment methods
// @Description List payment methods
// @Tags Payment Method
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param search query string false "Search"
// @Param page query int false "Page"
// @Param page_size query int false "Page size"
// @Success 200 {object} dto.ListPaymentMethodsResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /v1/payment-methods [get]
func (h PaymentMethodHandler) List(c *fiber.Ctx) error {
	search := c.Query(QueryParamSearch)
	paginationIn := parsePaginationParams(c)

	in := paymentmethod.ListPaymentMethodsUseCaseInput{
		PaginationInput: paginationIn,
		PaymentMethodOptions: repo.PaymentMethodOptions{
			Search: search,
		},
	}

	ctx := c.UserContext()
	res, err := h.lpm.Execute(ctx, in)
	if err != nil {
		return errs.New(err)
	}

	return c.JSON(res)
}
