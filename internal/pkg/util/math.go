package util

import (
	"github.com/shopspring/decimal"
)

// Trunc truncates float64 number to given precision.
func Trunc(number float64, precision int) float64 {
	result, _ := decimal.NewFromFloat(number).Truncate(int32(precision)).Float64()

	return result
}
