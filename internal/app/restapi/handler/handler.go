package handler

import (
	"strconv"

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
