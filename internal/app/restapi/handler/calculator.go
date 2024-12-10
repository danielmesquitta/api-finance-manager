package handler

import (
	"net/http"

	"github.com/danielmesquitta/api-finance-manager/internal/app/restapi/dto"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/usecase"
	"github.com/labstack/echo/v4"
)

type CalculatorHandler struct {
	cci *usecase.CalculateCompoundInterestUseCase
	cr  *usecase.CalculateRetirementUseCase
	csi *usecase.CalculateSimpleInterestUseCase
}

func NewCalculatorHandler(
	cci *usecase.CalculateCompoundInterestUseCase,
	cr *usecase.CalculateRetirementUseCase,
	csi *usecase.CalculateSimpleInterestUseCase,
) *CalculatorHandler {
	return &CalculatorHandler{
		cci: cci,
		cr:  cr,
		csi: csi,
	}
}

// @Summary Calculate compound interest
// @Description Calculate compound interest
// @Tags Calculator
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.CompoundInterestRequestDTO true "Request body"
// @Success 200 {object} dto.CompoundInterestResponseDTO
// @Failure 400 {object} dto.ErrorResponseDTO
// @Failure 401 {object} dto.ErrorResponseDTO
// @Failure 500 {object} dto.ErrorResponseDTO
// @Router /v1/calculator/compound-interest [post]
func (h CalculatorHandler) CompoundInterest(c echo.Context) error {
	body := &dto.CompoundInterestRequestDTO{}
	if err := c.Bind(body); err != nil {
		return errs.New(err)
	}

	output, err := h.cci.Execute(
		c.Request().Context(),
		body.CalculateCompoundInterestUseCaseInput,
	)
	if err != nil {
		return errs.New(err)
	}

	return c.JSON(http.StatusOK, dto.CompoundInterestResponseDTO{
		CalculateCompoundInterestUseCaseOutput: *output,
	})
}

// @Summary Calculate retirement
// @Description Calculate investments needed for retirement
// @Tags Calculator
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.RetirementRequestDTO true "Request body"
// @Success 200 {object} dto.RetirementResponseDTO
// @Failure 400 {object} dto.ErrorResponseDTO
// @Failure 401 {object} dto.ErrorResponseDTO
// @Failure 500 {object} dto.ErrorResponseDTO
// @Router /v1/calculator/retirement [post]
func (h CalculatorHandler) Retirement(c echo.Context) error {
	body := &dto.RetirementRequestDTO{}
	if err := c.Bind(body); err != nil {
		return errs.New(err)
	}

	output, err := h.cr.Execute(
		c.Request().Context(),
		body.CalculateRetirementUseCaseInput,
	)
	if err != nil {
		return errs.New(err)
	}

	return c.JSON(http.StatusOK, dto.RetirementResponseDTO{
		CalculateRetirementUseCaseOutput: *output,
	})
}

// @Summary Calculate simple interest
// @Description Calculate simple interest
// @Tags Calculator
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.SimpleInterestRequestDTO true "Request body"
// @Success 200 {object} dto.SimpleInterestResponseDTO
// @Failure 400 {object} dto.ErrorResponseDTO
// @Failure 401 {object} dto.ErrorResponseDTO
// @Failure 500 {object} dto.ErrorResponseDTO
// @Router /v1/calculator/simple-interest [post]
func (h CalculatorHandler) SimpleInterest(c echo.Context) error {
	body := &dto.SimpleInterestRequestDTO{}
	if err := c.Bind(body); err != nil {
		return errs.New(err)
	}

	output, err := h.csi.Execute(
		c.Request().Context(),
		body.CalculateSimpleInterestUseCaseInput,
	)
	if err != nil {
		return errs.New(err)
	}

	return c.JSON(http.StatusOK, dto.SimpleInterestResponseDTO{
		CalculateSimpleInterestUseCaseOutput: *output,
	})
}
