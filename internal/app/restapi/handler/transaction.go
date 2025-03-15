package handler

import (
	"net/http"

	"github.com/danielmesquitta/api-finance-manager/internal/app/restapi/dto"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type TransactionHandler struct {
	sa *usecase.SyncTransactions
	lt *usecase.ListTransactions
	gt *usecase.GetTransaction
	ut *usecase.UpdateTransaction
	ct *usecase.CreateTransaction
}

func NewTransactionHandler(
	sa *usecase.SyncTransactions,
	lt *usecase.ListTransactions,
	gt *usecase.GetTransaction,
	ut *usecase.UpdateTransaction,
	ct *usecase.CreateTransaction,
) *TransactionHandler {
	return &TransactionHandler{
		sa: sa,
		lt: lt,
		gt: gt,
		ut: ut,
		ct: ct,
	}
}

// @Summary Create transactions
// @Description Create transactions
// @Tags Transaction
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.CreateTransactionRequest true "Request body"
// @Success 201
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /v1/transactions [post]
func (h *TransactionHandler) Create(c *fiber.Ctx) error {
	in := usecase.CreateTransactionInput{}
	if err := c.BodyParser(&in); err != nil {
		return errs.New(err)
	}

	in.UserID = uuid.MustParse(GetUserClaims(c).Issuer)

	ctx := c.UserContext()
	if err := h.ct.Execute(ctx, in); err != nil {
		return errs.New(err)
	}

	return c.SendStatus(http.StatusCreated)
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
func (h TransactionHandler) Get(c *fiber.Ctx) error {
	claims := GetUserClaims(c)
	userID := uuid.MustParse(claims.Issuer)

	transactionID := uuid.MustParse(c.Params(pathParamTransactionID))

	in := usecase.GetTransactionInput{
		ID:     transactionID,
		UserID: userID,
	}

	ctx := c.UserContext()
	transaction, err := h.gt.Execute(ctx, in)
	if err != nil {
		return errs.New(err)
	}

	return c.JSON(dto.GetTransactionResponse{
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
func (h *TransactionHandler) Sync(c *fiber.Ctx) error {
	userIDs, err := parseUUIDsParam(c, QueryParamUserIDs)
	if err != nil {
		return errs.New(err)
	}

	in := usecase.SyncTransactionsInput{
		UserIDs: userIDs,
	}

	ctx := c.UserContext()
	if err := h.sa.Execute(ctx, in); err != nil {
		return errs.New(err)
	}

	return c.SendStatus(http.StatusNoContent)
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
func (h *TransactionHandler) List(c *fiber.Ctx) error {
	claims := GetUserClaims(c)
	userID := uuid.MustParse(claims.Issuer)

	paginationIn := parsePaginationParams(c)

	transactionOptions, err := prepareTransactionOptions(c)
	if err != nil {
		return errs.New(err)
	}

	in := usecase.ListTransactionsInput{
		PaginationInput:    paginationIn,
		TransactionOptions: *transactionOptions,
		UserID:             userID,
	}

	ctx := c.UserContext()
	res, err := h.lt.Execute(ctx, in)
	if err != nil {
		return errs.New(err)
	}

	return c.JSON(res)
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
func (h TransactionHandler) Update(c *fiber.Ctx) error {
	in := usecase.UpdateTransactionInput{}
	if err := c.BodyParser(&in); err != nil {
		return errs.New(err)
	}

	claims := GetUserClaims(c)
	userID := uuid.MustParse(claims.Issuer)

	transactionID := uuid.MustParse(c.Params(pathParamTransactionID))

	in.ID = transactionID
	in.UserID = userID

	ctx := c.UserContext()
	if err := h.ut.Execute(ctx, in); err != nil {
		return errs.New(err)
	}

	return c.SendStatus(http.StatusNoContent)
}
