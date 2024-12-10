package moneyutil

import "math"

func Round2Decimal(value float64) float64 {
	return math.Round(value*100) / 100
}
