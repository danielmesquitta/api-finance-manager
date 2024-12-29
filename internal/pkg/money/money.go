package money

import "math"

func decimalRound(value float64) float64 {
	return math.Round(value*100) / 100
}

func ToCents(value float64) int64 {
	return int64(decimalRound(value) * 100)
}

func FromCents(value int64) float64 {
	return float64(value) / 100
}

func ToMonthlyInterestRate(annualInterestRate float64) float64 {
	return math.Pow(1+annualInterestRate, 1.0/12) - 1
}

func ToAnnualInterestRate(monthlyInterestRate float64) float64 {
	return math.Pow(1+monthlyInterestRate, 12) - 1
}
