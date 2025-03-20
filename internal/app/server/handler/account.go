package handler

import (
	"net/http"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/usecase/account"
	"github.com/gofiber/fiber/v2"
)

type AccountHandler struct {
	ca  *account.CreateAccountsUseCase
	gab *account.GetAccountsBalanceUseCase
	sab *account.SyncAccountsBalancesUseCase
}

func NewAccountHandler(
	ca *account.CreateAccountsUseCase,
	gab *account.GetAccountsBalanceUseCase,
	sab *account.SyncAccountsBalancesUseCase,
) *AccountHandler {
	return &AccountHandler{
		ca:  ca,
		gab: gab,
		sab: sab,
	}
}

// @Summary Sync accounts from open finance
// @Description Webhook to sync user accounts from open finance
// @Tags Account
// @Security BasicAuth
// @Accept json
// @Produce json
// @Param request body dto.CreateAccountsRequest true "Request body"
// @Success 204
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /v1/admin/accounts [post]
func (h *AccountHandler) Create(c *fiber.Ctx) error {
	in := account.CreateAccountsUseCaseInput{}
	if err := c.BodyParser(&in); err != nil {
		return errs.New(err)
	}

	ctx := c.UserContext()
	if err := h.ca.Execute(ctx, in); err != nil {
		return errs.New(err)
	}

	return c.SendStatus(http.StatusNoContent)
}

// @Summary Get accounts balance
// @Description Gets user total balance and transactions monthly balance with previous month comparison
// @Tags Balance
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param search query string false "Search"
// @Param start_date query string false "Start date" format(date)
// @Param end_date query string false "End date" format(date)
// @Param institution_ids query []string false "Institution IDs"
// @Param category_ids query []string false "Category IDs"
// @Param payment_method_ids query []string false "Payment method IDs"
// @Param is_expense query bool false "Filter only expenses"
// @Param is_income query bool false "Filter only incomes"
// @Param is_ignored query bool false "Filter ignored or not ignored transactions"
// @Success 200 {object} dto.GetAccountsBalanceResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /v1/balances [get]
func (h *AccountHandler) GetBalance(c *fiber.Ctx) error {
	userID, _, err := GetUser(c)
	if err != nil {
		return errs.New(err)
	}

	transactionOptions, err := prepareTransactionOptions(c)
	if err != nil {
		return errs.New(err)
	}

	in := account.GetAccountsBalanceUseCaseInput{
		TransactionOptions: *transactionOptions,
		UserID:             userID,
	}

	ctx := c.UserContext()
	out, err := h.gab.Execute(ctx, in)
	if err != nil {
		return errs.New(err)
	}

	return c.JSON(out)
}

// @Summary Sync account balances from open finance
// @Description Sync account balances from open finance
// @Tags Balance
// @Security BasicAuth
// @Accept json
// @Produce json
// @Success 204
// @Failure 500 {object} dto.ErrorResponse
// @Router /v1/admin/balances/sync [post]
func (h *AccountHandler) Sync(c *fiber.Ctx) error {
	ctx := c.UserContext()
	if err := h.sab.Execute(ctx); err != nil {
		return errs.New(err)
	}

	return c.SendStatus(http.StatusNoContent)
}
