package handler

import (
	"net/http"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/usecase"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type TransactionHandler struct {
	sa *usecase.SyncTransactions
	lt *usecase.ListTransactions
}

func NewTransactionHandler(
	sa *usecase.SyncTransactions,
	lt *usecase.ListTransactions,
) *TransactionHandler {
	return &TransactionHandler{
		sa: sa,
		lt: lt,
	}
}

// @Summary Sync transactions from open finance
// @Description Webhook to sync transactions from open finance
// @Tags Transaction
// @Security BasicAuth
// @Accept json
// @Produce json
// @Success 204
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /v1/admin/transactions/sync [post]
func (h *TransactionHandler) Sync(c echo.Context) error {
	ctx := c.Request().Context()
	if err := h.sa.Execute(ctx); err != nil {
		return errs.New(err)
	}

	return c.NoContent(http.StatusNoContent)
}

// @Summary List transactions
// @Description List transactions
// @Tags Transaction
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param search query string false "Search"
// @Param page query int false "Page"
// @Param page_size query int false "Page size"
// @Param start_date query string false "Start date" format(date)
// @Param end_date query string false "End date" format(date)
// @Param institution_id query string false "Institution ID" format(uuid)
// @Param category_id query string false "Category ID" format(uuid)
// @Param payment_method_id query string false "Payment method ID" format(uuid)
// @Param is_expense query bool false "Filter only expenses"
// @Param is_income query bool false "Filter only incomes"
// @Param is_ignored query bool false "Filter ignored or not ignored transactions"
// @Success 200 {object} dto.ListTransactionsResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /v1/transactions [get]
func (h *TransactionHandler) List(c echo.Context) error {
	claims := getUserClaims(c)
	userID := uuid.MustParse(claims.Issuer)

	search := c.QueryParam(queryParamSearch)

	paginationIn := parsePaginationParams(c)

	startDate, endDate, err := parseDateFilterParams(c)
	if err != nil {
		return errs.New(err)
	}

	paymentMethodID, err := parseUUIDParam(c, queryParamPaymentMethodID)
	if err != nil {
		return errs.New(err)
	}

	institutionID, err := parseUUIDParam(c, queryParamInstitutionID)
	if err != nil {
		return errs.New(err)
	}

	categoryID, err := parseUUIDParam(c, queryParamCategoryID)
	if err != nil {
		return errs.New(err)
	}

	isExpense, err := parseBoolParam(c, queryParamIsExpense)
	if err != nil {
		return errs.New(err)
	}

	isIncome, err := parseBoolParam(c, queryParamIsExpense)
	if err != nil {
		return errs.New(err)
	}

	isIgnored, err := parseNillableBoolParam(c, queryParamIsIgnored)
	if err != nil {
		return errs.New(err)
	}

	in := usecase.ListTransactionsInput{
		PaginationInput: paginationIn,
		ListTransactionsOptions: repo.ListTransactionsOptions{
			Limit:           0,
			Offset:          0,
			Search:          search,
			StartDate:       startDate,
			EndDate:         endDate,
			CategoryID:      categoryID,
			InstitutionID:   institutionID,
			IsExpense:       isExpense,
			IsIncome:        isIncome,
			IsIgnored:       isIgnored,
			PaymentMethodID: paymentMethodID,
		},
		UserID: userID,
	}

	ctx := c.Request().Context()
	res, err := h.lt.Execute(ctx, in)
	if err != nil {
		return errs.New(err)
	}

	return c.JSON(http.StatusOK, res)
}
