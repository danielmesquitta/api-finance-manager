package handler

import (
	"net/http"

	"github.com/danielmesquitta/api-finance-manager/internal/app/restapi/dto"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/usecase"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type TransactionHandler struct {
	sa *usecase.SyncTransactions
	lt *usecase.ListTransactions
	gt *usecase.GetTransaction
	ut *usecase.UpdateTransaction
}

func NewTransactionHandler(
	sa *usecase.SyncTransactions,
	lt *usecase.ListTransactions,
	gt *usecase.GetTransaction,
	ut *usecase.UpdateTransaction,
) *TransactionHandler {
	return &TransactionHandler{
		sa: sa,
		lt: lt,
		gt: gt,
		ut: ut,
	}
}

// @Summary Get transaction
// @Description Get transaction
// @Tags Transaction
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param transaction_id path string true "Transaction ID" format(uuid)
// @Success 200 {object} dto.GetTransactionResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /v1/transactions/{transaction_id} [get]
func (h TransactionHandler) Get(c echo.Context) error {
	claims := getUserClaims(c)
	userID := uuid.MustParse(claims.Issuer)

	transactionID := uuid.MustParse(c.Param(pathParamTransactionID))

	in := usecase.GetTransactionInput{
		TransactionID: transactionID,
		UserID:        userID,
	}

	ctx := c.Request().Context()
	transaction, err := h.gt.Execute(ctx, in)
	if err != nil {
		return errs.New(err)
	}

	return c.JSON(http.StatusOK, dto.GetTransactionResponse{
		FullTransaction: *transaction,
	})
}

// @Summary Sync transactions from open finance
// @Description Webhook to sync transactions from open finance
// @Tags Transaction
// @Security BasicAuth
// @Accept json
// @Produce json
// @Param user_ids query []string false "User IDs"
// @Success 204
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /v1/admin/transactions/sync [post]
func (h *TransactionHandler) Sync(c echo.Context) error {
	userIDs, err := parseUUIDsParam(c, queryParamUserIDs)
	if err != nil {
		return errs.New(err)
	}

	in := usecase.SyncTransactionsInput{
		UserIDs: userIDs,
	}

	ctx := c.Request().Context()
	if err := h.sa.Execute(ctx, in); err != nil {
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
// @Param date query string false "Date" format(date)
// @Param institution_ids query []string false "Institution IDs"
// @Param category_ids query []string false "Category IDs"
// @Param payment_method_ids query []string false "Payment method IDs"
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

	paginationIn := parsePaginationParams(c)

	date, err := parseDateParam(c, queryParamDate)
	if err != nil {
		return errs.New(err)
	}

	transactionOptions, err := prepareTransactionOptions(c)
	if err != nil {
		return errs.New(err)
	}

	in := usecase.ListTransactionsInput{
		PaginationInput:    paginationIn,
		TransactionOptions: *transactionOptions,
		Date:               date,
		UserID:             userID,
	}

	ctx := c.Request().Context()
	res, err := h.lt.Execute(ctx, in)
	if err != nil {
		return errs.New(err)
	}

	return c.JSON(http.StatusOK, res)
}

// @Summary Update transaction
// @Description Update transaction
// @Tags Transaction
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.UpdateTransactionRequest true "Request body"
// @Param transaction_id path string true "Transaction ID" format(uuid)
// @Success 204
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /v1/transactions/{transaction_id} [put]
func (h TransactionHandler) Update(c echo.Context) error {
	claims := getUserClaims(c)
	userID := uuid.MustParse(claims.Issuer)

	transactionID := uuid.MustParse(c.Param(pathParamTransactionID))

	in := usecase.UpdateTransactionInput{
		ID:     transactionID,
		UserID: userID,
	}

	ctx := c.Request().Context()
	if err := h.ut.Execute(ctx, in); err != nil {
		return errs.New(err)
	}

	return c.NoContent(http.StatusNoContent)
}
