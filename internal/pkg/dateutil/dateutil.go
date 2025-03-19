package dateutil

import "time"

func MustParseISOString(date string) time.Time {
	parsedDate, err := time.Parse(time.RFC3339, date)
	if err != nil {
		panic(err)
	}
	return parsedDate
}

func ToDayStart(
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
		date.Location(),
	)
	return dayStart
}

func ToDayEnd(
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
		date.Location(),
	)

	monthEnd = monthEnd.Add(-time.Nanosecond)

	return monthEnd
}

func ToMonthStart(
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
		date.Location(),
	)
	return monthStart
}

func ToMonthEnd(
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
		date.Location(),
	)

	monthEnd = monthEnd.Add(-time.Nanosecond)

	return monthEnd
}

func ToMonthDay(
	date time.Time,
	day int,
) time.Time {
	monthSameDay := time.Date(
		date.Year(),
		date.Month(),
		day,
		0,
		0,
		0,
		0,
		date.Location(),
	)
	// If day is greater than the last day of the month, adjust to the last day
	for monthSameDay.Month() != date.Month() {
		monthSameDay = monthSameDay.AddDate(0, 0, -1)
	}
	return monthSameDay
}

type ComparisonDates struct {
	StartDate           time.Time `json:"start_date,omitzero"`
	EndDate             time.Time `json:"end_date,omitzero"`
	ComparisonStartDate time.Time `json:"comparison_start_date,omitzero"`
	ComparisonEndDate   time.Time `json:"comparison_end_date,omitzero"`
}

func CalculateComparisonDates(startDate, endDate time.Time) *ComparisonDates {
	if endDate.After(time.Now()) {
		endDate = time.Now()
	}

	if startDate.After(endDate) {
		startDate = endDate
	}

	startDate = ToDayStart(startDate)
	endDate = ToDayEnd(endDate)

	out := &ComparisonDates{
		StartDate: startDate,
		EndDate:   endDate,
	}

	isMonthComparison := startDate.Month() == endDate.Month() &&
		startDate.Year() == endDate.Year()
	isFullMonthComparison := startDate.Equal(ToMonthStart(startDate)) &&
		endDate.Equal(ToMonthEnd(endDate))

	if isMonthComparison {
		comparisonStartDate := startDate.AddDate(0, -1, 0)
		if comparisonStartDate.Month() == startDate.Month() {
			days := comparisonStartDate.Day()
			comparisonStartDate = comparisonStartDate.AddDate(0, 0, -days)
		}

		comparisonEndDate := endDate.AddDate(0, -1, 0)
		if comparisonEndDate.Month() == endDate.Month() {
			days := comparisonEndDate.Day()
			comparisonEndDate = comparisonEndDate.AddDate(0, 0, -days)
		}

		if isFullMonthComparison {
			comparisonStartDate = ToMonthStart(comparisonStartDate)
			comparisonEndDate = ToMonthEnd(comparisonEndDate)
		}

		out.ComparisonStartDate = comparisonStartDate
		out.ComparisonEndDate = comparisonEndDate
	} else {
		duration := endDate.Sub(startDate)
		durationDays := int(duration.Hours() / 24)
		out.ComparisonEndDate = startDate.Add(-time.Nanosecond)
		out.ComparisonStartDate = ToDayStart(startDate.AddDate(0, 0, -durationDays))
	}

	return out
}
