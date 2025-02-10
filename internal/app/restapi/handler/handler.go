package handler

import (
	"strconv"
	"time"

	"github.com/danielmesquitta/api-finance-manager/internal/app/restapi/middleware"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/usecase"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/jwtutil"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type QueryParam = string

const (
	queryParamSearch          QueryParam = "search"
	queryParamPage            QueryParam = "page"
	queryParamPageSize        QueryParam = "page_size"
	queryParamDate            QueryParam = "date"
	queryParamStartDate       QueryParam = "start_date"
	queryParamEndDate         QueryParam = "end_date"
	queryParamInstitutionID   QueryParam = "institution_id"
	queryParamCategoryID      QueryParam = "category_id"
	queryParamIsExpense       QueryParam = "is_expense"
	queryParamIsIncome        QueryParam = "is_income"
	queryParamIsIgnored       QueryParam = "is_ignored"
	queryParamPaymentMethodID QueryParam = "payment_method_id"
)

type PathParam = string

const (
	pathParamCategoryID    PathParam = "category_id"
	pathParamTransactionID PathParam = "transaction_id"
)

func parsePaginationParams(
	c echo.Context,
) usecase.PaginationInput {
	page, _ := strconv.Atoi(c.QueryParam(queryParamPage))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(c.QueryParam(queryParamPageSize))
	if pageSize < 1 {
		pageSize = 20
	}

	return usecase.PaginationInput{
		Page:     uint(page),
		PageSize: uint(pageSize),
	}
}

func parseDateParam(
	c echo.Context,
	param QueryParam,
) (time.Time, error) {
	paramValue := c.QueryParam(param)
	if paramValue == "" {
		return time.Time{}, nil
	}

	date, err := time.Parse(time.RFC3339, paramValue)
	if err != nil {
		return time.Time{}, errs.ErrInvalidDate
	}

	return date, nil
}

func parseUUIDsParam(
	c echo.Context,
	param QueryParam,
) ([]uuid.UUID, error) {
	values := c.QueryParams()[param]
	if len(values) == 0 {
		return nil, nil
	}

	var uuids []uuid.UUID
	for _, value := range values {
		if value == "" {
			continue
		}
		id, err := uuid.Parse(value)
		if err != nil {
			return nil, errs.ErrInvalidUUID
		}
		uuids = append(uuids, id)
	}

	return uuids, nil
}

func parseBoolParam(
	c echo.Context,
	param QueryParam,
) (bool, error) {
	paramValue := c.QueryParam(param)
	if paramValue == "" {
		return false, nil
	}

	b, err := strconv.ParseBool(paramValue)
	if err != nil {
		return false, errs.ErrInvalidBool
	}

	return b, nil
}

func parseNillableBoolParam(
	c echo.Context,
	param QueryParam,
) (*bool, error) {
	paramValue := c.QueryParam(param)
	if paramValue == "" {
		return nil, nil
	}

	b, err := strconv.ParseBool(paramValue)
	if err != nil {
		return nil, errs.ErrInvalidBool
	}

	return &b, nil
}

func getUserClaims(
	c echo.Context,
) *jwtutil.UserClaims {
	return c.Get(middleware.ClaimsKey).(*jwtutil.UserClaims)
}

func prepareTransactionOptions(
	c echo.Context,
) (*repo.TransactionOptions, error) {
	search := c.QueryParam(queryParamSearch)

	paymentMethodIDs, err := parseUUIDsParam(c, queryParamPaymentMethodID)
	if err != nil {
		return nil, errs.New(err)
	}

	institutionIDs, err := parseUUIDsParam(c, queryParamInstitutionID)
	if err != nil {
		return nil, errs.New(err)
	}

	categoryIDs, err := parseUUIDsParam(c, queryParamCategoryID)
	if err != nil {
		return nil, errs.New(err)
	}

	isExpense, err := parseBoolParam(c, queryParamIsExpense)
	if err != nil {
		return nil, errs.New(err)
	}

	isIncome, err := parseBoolParam(c, queryParamIsExpense)
	if err != nil {
		return nil, errs.New(err)
	}

	isIgnored, err := parseNillableBoolParam(c, queryParamIsIgnored)
	if err != nil {
		return nil, errs.New(err)
	}

	return &repo.TransactionOptions{
		Search:           search,
		CategoryIDs:      categoryIDs,
		InstitutionIDs:   institutionIDs,
		PaymentMethodIDs: paymentMethodIDs,
		IsExpense:        isExpense,
		IsIncome:         isIncome,
		IsIgnored:        isIgnored,
	}, nil
}
