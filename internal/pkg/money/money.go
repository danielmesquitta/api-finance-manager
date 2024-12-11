package money

import "math"

func Round(value float64) float64 {
	return math.Round(value*100) / 100
}

func CompoundInterestAnnualToMonthlyRate(annualRate float64) float64 {
	return math.Pow(1+annualRate, 1.0/12) - 1
}

func CompoundInterestMonthlyToAnnualRate(monthlyRate float64) float64 {
	return math.Pow(1+monthlyRate, 12) - 1
}
