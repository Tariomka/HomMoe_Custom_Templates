package helpers

import "math"

func RoundWithPrecision(value float64, decimalPrecision int) float64 {
	multiplier := math.Pow(10, float64(decimalPrecision))
	return math.Round(value*multiplier) / multiplier
}
