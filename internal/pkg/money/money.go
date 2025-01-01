package money

import (
	"github.com/shopspring/decimal"
)

// FromCents converts an int64 value (representing cents) to a
// decimal.Decimal value (representing the major unit from a currency).
func FromCents(value int64) decimal.Decimal {
	return decimal.New(value, -2)
}

type FloatOrDecimal interface {
	~float64 | decimal.Decimal
}

// ToCents converts a float64 or decimal.Decimal value to an int64 value
// representing cents. It is equivalent to math.Round(value * 100), but using decimal.Decimal
func ToCents[T FloatOrDecimal](value T) int64 {
	oneHundred := decimal.New(100, 0)

	switch v := any(value).(type) {
	case decimal.Decimal:
		return v.
			Mul(oneHundred).
			Round(0).
			IntPart()

	case float64:
		return decimal.
			NewFromFloat(v).
			Mul(oneHundred).
			Round(0).
			IntPart()

	default:
		panic("unsupported type")
	}
}

// ToMonthlyInterestRate converts an annual to a monthly interest rate,
// it is equivalent to math.Pow(1+annualInterestRate, 1.0/12) - 1, but using decimal.Decimal
func ToMonthlyInterestRate(annualInterestRate decimal.Decimal) decimal.Decimal {
	// Add 1 to the annual interest rate
	one := decimal.New(1, 0)
	base := one.Add(annualInterestRate)

	// Compute the 1/12th power
	twelve := decimal.New(12, 0)
	exponent := one.Div(twelve)
	monthlyRate := base.Pow(exponent)

	// Subtract 1 to get the monthly interest rate
	monthlyRate = monthlyRate.Sub(one)

	return monthlyRate
}

// ToAnnualInterestRate converts a monthly to an annual interest rate,
// it is equivalent to math.Pow(1+monthlyInterestRate, 12) - 1, but using decimal.Decimal
func ToAnnualInterestRate(monthlyInterestRate decimal.Decimal) decimal.Decimal {
	// Add 1 to the monthly interest rate
	one := decimal.New(1, 0)
	base := one.Add(monthlyInterestRate)

	// Compute the 12th power
	exponent := decimal.New(12, 0)
	annualRate := base.Pow(exponent)

	// Subtract 1 to get the annual interest rate
	annualRate = annualRate.Sub(one)

	return annualRate
}
