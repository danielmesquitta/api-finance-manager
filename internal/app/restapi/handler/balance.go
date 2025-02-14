package handler

import (
	"net/http"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/usecase"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
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
// @Param date query string false "Date" format(date)
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
func (h *BalanceHandler) Get(c echo.Context) error {
	claims := getUserClaims(c)
	userID := uuid.MustParse(claims.Issuer)

	date, err := parseDateParam(c, queryParamDate)
	if err != nil {
		return errs.New(err)
	}

	transactionOptions, err := prepareTransactionOptions(c)
	if err != nil {
		return errs.New(err)
	}

	in := usecase.GetBalanceInput{
		TransactionOptions: *transactionOptions,
		Date:               date,
		UserID:             userID,
	}

	ctx := c.Request().Context()
	res, err := h.gb.Execute(ctx, in)
	if err != nil {
		return errs.New(err)
	}

	return c.JSON(http.StatusOK, res)
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
func (h BalanceHandler) Sync(c echo.Context) error {
	ctx := c.Request().Context()
	if err := h.sb.Execute(ctx); err != nil {
		return errs.New(err)
	}

	return c.NoContent(http.StatusNoContent)
}
