package dateutil

import (
	"reflect"
	"testing"
	"time"
)

func TestCalculateComparisonDates(t *testing.T) {
	type args struct {
		startDate time.Time
		endDate   time.Time
	}
	type Test struct {
		name string
		args args
		want *ComparisonDates
	}
	tests := []Test{
		func() Test {
			startDate := MustParseISOString("2024-11-01T00:00:00-03:00")
			endDate := MustParseISOString("2024-11-30T23:59:59.999999999-03:00")
			cmpStartDate := MustParseISOString("2024-10-01T00:00:00-03:00")
			cmpEndDate := MustParseISOString(
				"2024-10-31T23:59:59.999999999-03:00",
			)

			return Test{
				name: "Full month comparison",
				args: args{
					startDate: startDate,
					endDate:   endDate,
				},
				want: &ComparisonDates{
					StartDate:           startDate,
					EndDate:             endDate,
					ComparisonStartDate: cmpStartDate,
					ComparisonEndDate:   cmpEndDate,
				},
			}
		}(),
		func() Test {
			startDate := MustParseISOString("2024-11-01T00:00:00-03:00")
			endDate := MustParseISOString("2024-11-15T23:59:59.999999999-03:00")
			cmpStartDate := MustParseISOString("2024-10-01T00:00:00-03:00")
			cmpEndDate := MustParseISOString(
				"2024-10-15T23:59:59.999999999-03:00",
			)

			return Test{
				name: "Partial month comparison",
				args: args{
					startDate: startDate,
					endDate:   endDate,
				},
				want: &ComparisonDates{
					StartDate:           startDate,
					EndDate:             endDate,
					ComparisonStartDate: cmpStartDate,
					ComparisonEndDate:   cmpEndDate,
				},
			}
		}(),
		func() Test {
			startDate := MustParseISOString("2024-10-31T00:00:00-03:00")
			endDate := MustParseISOString("2024-11-11T23:59:59.999999999-03:00")
			cmpStartDate := MustParseISOString("2024-10-20T00:00:00-03:00")
			cmpEndDate := MustParseISOString(
				"2024-10-30T23:59:59.999999999-03:00",
			)

			return Test{
				name: "10 days period comparison",
				args: args{
					startDate: startDate,
					endDate:   endDate,
				},
				want: &ComparisonDates{
					StartDate:           startDate,
					EndDate:             endDate,
					ComparisonStartDate: cmpStartDate,
					ComparisonEndDate:   cmpEndDate,
				},
			}
		}(),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CalculateComparisonDates(tt.args.startDate, tt.args.endDate); !reflect.DeepEqual(
				got,
				tt.want,
			) {
				t.Errorf(
					"calculateComparisonDates() = %+v, want %+v",
					got,
					tt.want,
				)
			}
		})
	}
}
