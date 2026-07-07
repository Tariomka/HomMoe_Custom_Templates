package helpers

import "math"

func RoundWithPrecision(value float64, decimalPrecision int) float64 {
	multiplier := math.Pow(10, float64(decimalPrecision))
	return math.Round(value*multiplier) / multiplier
}

// Clamp returns value clamped to [lowest, highest].
func Clamp(value, lowest, highest int) int {
	if value < lowest {
		return lowest
	}

	if value > highest {
		return highest
	}

	return value
}

// ClampFloat returns value clamped to [lowest, highest].
func ClampFloat(value, lowest, highest float64) float64 {
	if value < lowest {
		return lowest
	}

	if value > highest {
		return highest
	}

	return value
}

// Scale returns value scaled by multiplier, clamped to a minimum of 0.
func Scale(value, multiplier float64) int {
	return max(0, int(value*multiplier))
}

// BoolToInt returns 1 if boolean is true, otherwise 0. This is a compiler optimized function.
func BoolToInt(boolean bool) int {
	if boolean {
		return 1
	}

	return 0
}
