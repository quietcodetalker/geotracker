package util

import "math"

// Trunc truncates float64 number to given precision.
func Trunc(number float64, precision int) float64 {
	pow10 := math.Pow10(precision)
	return math.Trunc(number*pow10) / pow10
}
