package usecase

import (
	"math"
	"time"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/money"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
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

func toDayStart(
	date time.Time,
) time.Time {
	dayStart := time.Date(
		date.Year(),
		date.Month(),
		date.Day(),
		0,
		0,
		0,
		0,
		time.Local,
	)
	return dayStart
}

func toDayEnd(
	date time.Time,
) time.Time {
	monthEnd := time.Date(
		date.Year(),
		date.Month(),
		date.Day()+1,
		0,
		0,
		0,
		0,
		time.Local,
	)

	monthEnd = monthEnd.Add(-time.Nanosecond)

	return monthEnd
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

type ComparisonDates struct {
	StartDate           time.Time `json:"start_date,omitzero"`
	EndDate             time.Time `json:"end_date,omitzero"`
	ComparisonStartDate time.Time `json:"comparison_start_date,omitzero"`
	ComparisonEndDate   time.Time `json:"comparison_end_date,omitzero"`
}

func calculateComparisonDates(startDate, endDate time.Time) *ComparisonDates {
	if endDate.After(time.Now()) {
		endDate = time.Now()
	}

	if startDate.After(endDate) {
		startDate = endDate
	}

	startDate = toDayStart(startDate)
	endDate = toDayEnd(endDate)

	out := &ComparisonDates{
		StartDate: startDate,
		EndDate:   endDate,
	}

	isMonthComparison := startDate.Month() == endDate.Month() &&
		startDate.Year() == endDate.Year()
	isFullMonthComparison := startDate.Equal(toMonthStart(startDate)) &&
		endDate.Equal(toMonthEnd(endDate))

	if isMonthComparison {
		comparisonStartDate := startDate.AddDate(0, -1, 0)
		if comparisonStartDate.Month() == startDate.Month() {
			days := comparisonStartDate.Day()
			comparisonStartDate = comparisonStartDate.AddDate(0, 0, -days)
		}

		comparisonEndDate := endDate.AddDate(0, -1, 0)
		if comparisonEndDate.Month() == endDate.Month() {
			days := comparisonEndDate.Day()
			comparisonStartDate = comparisonStartDate.AddDate(0, 0, -days)
		}

		if isFullMonthComparison {
			comparisonStartDate = toMonthStart(comparisonStartDate)
			comparisonEndDate = toMonthEnd(comparisonEndDate)
		}

		out.ComparisonStartDate = comparisonStartDate
		out.ComparisonEndDate = comparisonEndDate
	} else {
		duration := endDate.Sub(startDate)
		durationDays := int(duration.Hours() / 24)
		out.ComparisonEndDate = startDate.AddDate(0, 0, -1)
		out.ComparisonStartDate = out.ComparisonEndDate.AddDate(0, 0, -durationDays)
	}

	return out
}

func prepareTransactionOptions(
	in repo.TransactionOptions,
) []repo.TransactionOption {
	opts := []repo.TransactionOption{}

	if in.Search != "" {
		opts = append(opts, repo.WithTransactionSearch(in.Search))
	}

	if len(in.CategoryIDs) > 0 {
		opts = append(
			opts,
			repo.WithTransactionCategories(in.CategoryIDs...),
		)
	}

	if len(in.InstitutionIDs) > 0 {
		opts = append(
			opts,
			repo.WithTransactionInstitutions(in.InstitutionIDs...),
		)
	}

	if len(in.PaymentMethodIDs) > 0 {
		opts = append(
			opts,
			repo.WithTransactionPaymentMethods(in.PaymentMethodIDs...),
		)
	}

	if !in.StartDate.IsZero() {
		opts = append(
			opts,
			repo.WithTransactionDateAfter(in.StartDate),
		)
	}

	if !in.EndDate.IsZero() {
		opts = append(
			opts,
			repo.WithTransactionDateBefore(in.EndDate),
		)
	}

	if in.IsExpense {
		opts = append(
			opts,
			repo.WithTransactionIsExpense(in.IsExpense),
		)
	}

	if in.IsIncome {
		opts = append(
			opts,
			repo.WithTransactionIsIncome(in.IsIncome),
		)
	}

	if in.IsIgnored != nil {
		opts = append(
			opts,
			repo.WithTransactionIsIgnored(*in.IsIgnored),
		)
	}

	return opts
}

func calculatePercentageVariation(
	curr, prev int64,
) int64 {
	if prev == 0 {
		return 0
	}
	return money.FromPercentage(1 - (float64(curr) / float64(prev)))
}
