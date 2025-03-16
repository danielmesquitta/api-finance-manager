package money

import (
	"math"
)

const percentageFactor = 10000
const centsFactor = 100

func decimalRound(value float64) float64 {
	return math.Round(value*centsFactor) / centsFactor
}

// ToPercentage converts a int64 to a float64 percentage.
func ToPercentage(value int64) float64 {
	return float64(value) / percentageFactor
}

func FromPercentage(value float64) int64 {
	return int64(value * percentageFactor)
}

// ToCents converts a float64 to cents.
func ToCents(value float64) int64 {
	return int64(decimalRound(value) * centsFactor)
}

// FromCents converts cents to a float64.
func FromCents(value int64) float64 {
	return float64(value) / centsFactor
}

// ToMonthlyInterestRate converts an annual interest rate to a monthly interest rate.
func ToMonthlyInterestRate(annualInterestRate float64) float64 {
	return math.Pow(1+annualInterestRate, 1.0/12) - 1
}

// ToAnnualInterestRate converts a monthly interest rate to an annual interest rate.
func ToAnnualInterestRate(monthlyInterestRate float64) float64 {
	return math.Pow(1+monthlyInterestRate, 12) - 1
}

// CalculatePercentageVariation calculates the percentage variation between two values.
func CalculatePercentageVariation(
	curr, prev int64,
) int64 {
	if prev == 0 {
		return 0
	}
	return FromPercentage((float64(curr) / float64(prev)) - 1)
}
