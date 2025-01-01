package money

import (
	"math"
)

func decimalRound(value float64) float64 {
	return math.Round(value*100) / 100
}

// ToCents converts a float64 to cents.
func ToPercentage(value int64) float64 {
	return float64(value) / 10000
}

// ToCents converts a float64 to cents.
func ToCents(value float64) int64 {
	return int64(decimalRound(value) * 100)
}

// FromCents converts cents to a float64.
func FromCents(value int64) float64 {
	return float64(value) / 100
}

// ToMonthlyInterestRate converts an annual interest rate to a monthly interest rate.
func ToMonthlyInterestRate(annualInterestRate float64) float64 {
	return math.Pow(1+annualInterestRate, 1.0/12) - 1
}

func ToAnnualInterestRate(monthlyInterestRate float64) float64 {
	return math.Pow(1+monthlyInterestRate, 12) - 1
}
