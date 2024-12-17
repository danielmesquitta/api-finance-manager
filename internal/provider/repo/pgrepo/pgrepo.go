package pgrepo

import (
	"math"
)

func calculateMaxLevenshteinDistance(search string, percentage float64) int {
	searchLength := float64(len(search))
	maxLevenshteinDistance := int(math.Floor(percentage * searchLength))
	return maxLevenshteinDistance
}
