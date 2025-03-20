package handler

import (
	"github.com/danielmesquitta/api-finance-manager/internal/app/server/dto"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/usecase/calc"
	"github.com/gofiber/fiber/v2"
)

type CalculatorHandler struct {
	cci *calc.CalculateCompoundInterestUseCase
	cer *calc.CalculateEmergencyReserveUseCase
	cr  *calc.CalculateRetirementUseCase
	csi *calc.CalculateSimpleInterestUseCase
	cvi *calc.CalculateCashVsInstallmentsUseCase
}

func NewCalculatorHandler(
	cci *calc.CalculateCompoundInterestUseCase,
	cer *calc.CalculateEmergencyReserveUseCase,
	cr *calc.CalculateRetirementUseCase,
	csi *calc.CalculateSimpleInterestUseCase,
	cvi *calc.CalculateCashVsInstallmentsUseCase,
) *CalculatorHandler {
	return &CalculatorHandler{
		cci: cci,
		cer: cer,
		cr:  cr,
		csi: csi,
		cvi: cvi,
	}
}

// @Summary Calculate compound interest
// @Description Calculate compound interest
// @Tags Calculator
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.CompoundInterestRequest true "Request body"
// @Success 200 {object} dto.CompoundInterestResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /v1/calculator/compound-interest [post]
func (h CalculatorHandler) CompoundInterest(c *fiber.Ctx) error {
	body := &dto.CompoundInterestRequest{}
	if err := c.BodyParser(body); err != nil {
		return errs.New(err)
	}

	ctx := c.UserContext()
	out, err := h.cci.Execute(
		ctx,
		body.CalculateCompoundInterestUseCaseInput,
	)
	if err != nil {
		return errs.New(err)
	}

	return c.JSON(dto.CompoundInterestResponse{
		CalculateCompoundInterestUseCaseOutput: *out,
	})
}

// @Summary Calculate emergency reserve
// @Description Calculate emergency reserve
// @Tags Calculator
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.EmergencyReserveRequest true "Request body"
// @Success 200 {object} dto.EmergencyReserveResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /v1/calculator/emergency-reserve [post]
func (h CalculatorHandler) EmergencyReserve(c *fiber.Ctx) error {
	body := &dto.EmergencyReserveRequest{}
	if err := c.BodyParser(body); err != nil {
		return errs.New(err)
	}

	ctx := c.UserContext()
	out, err := h.cer.Execute(
		ctx,
		body.CalculateEmergencyReserveUseCaseInput,
	)
	if err != nil {
		return errs.New(err)
	}

	return c.JSON(dto.EmergencyReserveResponse{
		CalculateEmergencyReserveUseCaseOutput: *out,
	})
}

// @Summary Calculate retirement
// @Description Calculate investments needed for retirement
// @Tags Calculator
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.RetirementRequest true "Request body"
// @Success 200 {object} dto.RetirementResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /v1/calculator/retirement [post]
func (h CalculatorHandler) Retirement(c *fiber.Ctx) error {
	body := &dto.RetirementRequest{}
	if err := c.BodyParser(body); err != nil {
		return errs.New(err)
	}

	ctx := c.UserContext()
	out, err := h.cr.Execute(
		ctx,
		body.CalculateRetirementUseCaseInput,
	)
	if err != nil {
		return errs.New(err)
	}

	return c.JSON(dto.RetirementResponse{
		CalculateRetirementUseCaseOutput: *out,
	})
}

// @Summary Calculate simple interest
// @Description Calculate simple interest
// @Tags Calculator
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.SimpleInterestRequest true "Request body"
// @Success 200 {object} dto.SimpleInterestResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /v1/calculator/simple-interest [post]
func (h CalculatorHandler) SimpleInterest(c *fiber.Ctx) error {
	body := &dto.SimpleInterestRequest{}
	if err := c.BodyParser(body); err != nil {
		return errs.New(err)
	}

	ctx := c.UserContext()
	out, err := h.csi.Execute(
		ctx,
		body.CalculateSimpleInterestUseCaseInput,
	)
	if err != nil {
		return errs.New(err)
	}

	return c.JSON(dto.SimpleInterestResponse{
		CalculateSimpleInterestUseCaseOutput: *out,
	})
}

// @Summary Calculate cash vs installments
// @Description Calculate cash vs installments
// @Tags Calculator
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.CashVsInstallmentsRequest true "Request body"
// @Success 200 {object} dto.CashVsInstallmentsResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /v1/calculator/cash-vs-installments [post]
func (h CalculatorHandler) CashVsInstallments(c *fiber.Ctx) error {
	body := &dto.CashVsInstallmentsRequest{}
	if err := c.BodyParser(body); err != nil {
		return errs.New(err)
	}

	ctx := c.UserContext()
	out, err := h.cvi.Execute(
		ctx,
		body.CalculateCashVsInstallmentsUseCaseInput,
	)
	if err != nil {
		return errs.New(err)
	}

	return c.JSON(dto.CashVsInstallmentsResponse{
		CalculateCashVsInstallmentsUseCaseOutput: *out,
	})
}
