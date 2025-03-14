package handler

import (
	"net/http"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type BalanceHandler struct {
	gb *usecase.GetBalance
	sb *usecase.SyncBalances
}

func NewBalanceHandler(
	gb *usecase.GetBalance,
	sb *usecase.SyncBalances,
) *BalanceHandler {
	return &BalanceHandler{
		gb: gb,
		sb: sb,
	}
}

// @Summary Get balance
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
// @Success 200 {object} dto.GetBalanceResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /v1/balances [get]
func (h *BalanceHandler) Get(c *fiber.Ctx) error {
	claims := GetUserClaims(c)
	userID := uuid.MustParse(claims.Issuer)

	transactionOptions, err := prepareTransactionOptions(c)
	if err != nil {
		return errs.New(err)
	}

	in := usecase.GetBalanceInput{
		TransactionOptions: *transactionOptions,
		UserID:             userID,
	}

	ctx := c.UserContext()
	out, err := h.gb.Execute(ctx, in)
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
func (h BalanceHandler) Sync(c *fiber.Ctx) error {
	ctx := c.UserContext()
	if err := h.sb.Execute(ctx); err != nil {
		return errs.New(err)
	}

	return c.SendStatus(http.StatusNoContent)
}
