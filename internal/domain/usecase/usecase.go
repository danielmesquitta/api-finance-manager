package usecase

import (
	"time"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
)

type PaginationInput struct {
	Search   string `json:"search"`
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
}

func preparePaginationInput(in *PaginationInput) (offset int) {
	if in.Page == 0 {
		in.Page = 1
	}
	if in.PageSize == 0 {
		in.PageSize = 20
	}
	return (in.Page - 1) * in.PageSize
}

func parseDateToMonthStart(dateStr string) (monthStart time.Time, err error) {
	date, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return time.Time{}, errs.ErrInvalidDate
	}
	monthStart = time.Date(
		date.Year(),
		date.Month(),
		1,
		0,
		0,
		0,
		0,
		time.Local,
	)
	return monthStart, nil
}
