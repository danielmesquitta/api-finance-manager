package handler

import (
	"strconv"

	"github.com/danielmesquitta/api-finance-manager/internal/app/restapi/middleware"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/jwtutil"
	"github.com/labstack/echo/v4"
)

const (
	queryParamSearch   = "search"
	queryParamPage     = "page"
	queryParamPageSize = "page_size"
)

func getPaginationParams(
	c echo.Context,
) (search string, page int, pageSize int) {
	search = c.QueryParam(queryParamSearch)
	page, _ = strconv.Atoi(c.QueryParam(queryParamPage))
	pageSize, _ = strconv.Atoi(c.QueryParam(queryParamPageSize))
	return search, page, pageSize
}

func getUserClaims(
	c echo.Context,
) *jwtutil.UserClaims {
	return c.Get(middleware.ClaimsKey).(*jwtutil.UserClaims)
}
