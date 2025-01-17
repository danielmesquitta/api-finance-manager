package usecase

import (
	"math"
	"time"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
)

type PaginationInput struct {
	Page     uint `json:"page"`
	PageSize uint `json:"page_size"`
}

func preparePaginationInput(in *PaginationInput) (offset uint) {
	if in.Page < 1 {
		in.Page = 1
	}
	if in.PageSize < 1 {
		in.PageSize = 20
	}
	return (in.Page - 1) * in.PageSize
}

func preparePaginationOutput[T any](
	out *entity.PaginatedList[T],
	in PaginationInput,
	count int64,
) {
	out.Page = in.Page
	out.PageSize = in.PageSize
	out.TotalItems = uint(count)
	out.TotalPages = uint(math.Ceil(float64(count) / float64(in.PageSize)))
}

func toMonthStart(
	date time.Time,
) time.Time {
	monthStart := time.Date(
		date.Year(),
		date.Month(),
		1,
		0,
		0,
		0,
		0,
		time.Local,
	)
	return monthStart
}

func toMonthEnd(
	date time.Time,
) time.Time {
	monthEnd := time.Date(
		date.Year(),
		date.Month()+1,
		1,
		0,
		0,
		0,
		0,
		time.Local,
	)

	monthEnd = monthEnd.Add(-time.Nanosecond)

	return monthEnd
}

func toMonthDay(
	monthStart time.Time,
	day int,
) time.Time {
	monthSameDay := time.Date(
		monthStart.Year(),
		monthStart.Month(),
		day,
		0,
		0,
		0,
		0,
		time.Local,
	)
	// If the previous month has less days than the current month, we need to
	// go back until we reach the last day of the month
	for monthSameDay.Month() != monthStart.Month() {
		monthSameDay = monthSameDay.AddDate(0, 0, -1)
	}
	return monthSameDay
}

func getStartOfDay(date time.Time) time.Time {
	startOfDay := time.Date(
		date.Year(),
		date.Month(),
		date.Day(),
		0,
		0,
		0,
		0,
		time.Local,
	)

	return startOfDay
}
