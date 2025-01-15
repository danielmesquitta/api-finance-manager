package handler

import (
	"strconv"
	"time"

	"github.com/danielmesquitta/api-finance-manager/internal/app/restapi/middleware"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/usecase"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/jwtutil"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type QueryParam = string

const (
	queryParamSearch        QueryParam = "search"
	queryParamPage          QueryParam = "page"
	queryParamPageSize      QueryParam = "page_size"
	queryParamStartDate     QueryParam = "start_date"
	queryParamEndDate       QueryParam = "end_date"
	queryParamInstitutionID QueryParam = "institution_id"
	queryParamCategoryID    QueryParam = "category_id"
	queryParamIsExpense     QueryParam = "is_expense"
	queryParamIsIncome      QueryParam = "is_income"
	queryParamPaymentMethod QueryParam = "payment_method"
)

func parsePaginationParams(
	c echo.Context,
) usecase.PaginationInput {
	page, _ := strconv.Atoi(c.QueryParam(queryParamPage))
	pageSize, _ := strconv.Atoi(c.QueryParam(queryParamPageSize))

	return usecase.PaginationInput{
		Page:     uint(page),
		PageSize: uint(pageSize),
	}
}

func parseDateFilterParams(
	c echo.Context,
) (startDate time.Time, endDate time.Time, err error) {
	startDateStr := c.QueryParam(queryParamStartDate)

	if startDateStr != "" {
		startDate, err = time.Parse(time.RFC3339, startDateStr)
		if err != nil {
			return time.Time{}, time.Time{}, errs.ErrInvalidDate
		}
	}

	endDateStr := c.QueryParam(queryParamEndDate)
	if endDateStr != "" {
		endDate, err = time.Parse(time.RFC3339, endDateStr)
		if err != nil {
			return time.Time{}, time.Time{}, errs.ErrInvalidDate
		}
	}

	return startDate, endDate, nil
}

func parseUUIDParam(
	c echo.Context,
	param QueryParam,
) (uuid.UUID, error) {
	paramValue := c.QueryParam(param)
	if paramValue == "" {
		return uuid.Nil, nil
	}

	id, err := uuid.Parse(paramValue)
	if err != nil {
		return uuid.Nil, errs.ErrInvalidUUID
	}

	return id, nil
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

func getUserClaims(
	c echo.Context,
) *jwtutil.UserClaims {
	return c.Get(middleware.ClaimsKey).(*jwtutil.UserClaims)
}
