package utils

// Denormalize returns the inverse of Normalize(float32, float32, float32), i.e. converts value from range [0,1] -> [low, high]
func Denormalize(value, low, high float32) float32 {
	if value < 0 {
		value = 0
	}
	if value > 1 {
		value = 1
	}
	return low + value*(high-low)
}

// Normalize returns a normalized value, i.e. converts value from range [low, high] to [0, 1]
func Normalize(value, low, high float32) float32 {
	if high == low {
		return 0
	}
	ratio := (value - low) / (high - low)
	if ratio < 0 {
		ratio = 0
	}
	if ratio > 1 {
		ratio = 1
	}
	return ratio
}

func RoundHalfAway(value float64) float64 {
	if value < 0 {
		return -RoundHalfAway(-value)
	}
	return float64(int(value + 0.5))
}

// RoundedRange snaps a [0,1] slider value to the nearest integer in [low, high].
func RoundedRange(value float32, low, high int) int {
	return min(max(int(RoundHalfAway(float64(Denormalize(value, float32(low), float32(high))))), low), high)
}
