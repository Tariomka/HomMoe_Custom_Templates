package utils

import (
	"math"

	"github.com/Tariomka/hommoe_custom_templates/internal/constants"
)

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

func SliderToMapSize(value float32, includeExp bool) int {
	all := constants.MapSizes
	if includeExp {
		all = append(append([]int{}, constants.MapSizes...), constants.ExpMapSizes...)
	}
	if len(all) == 1 {
		return all[0]
	}
	idx := max(int(math.Round(float64(value)*float64(len(all)-1))), 0)
	if idx >= len(all) {
		idx = len(all) - 1
	}
	return all[idx]
}

// MapSizeToSlider returns the [0,1] slider position for a map size value.
func MapSizeToSlider(size int, includeExp bool) float32 {
	all := constants.MapSizes
	if includeExp {
		all = append(append([]int{}, constants.MapSizes...), constants.ExpMapSizes...)
	}
	for i, value := range all {
		if value == size {
			if len(all) <= 1 {
				return 0
			}
			return float32(i) / float32(len(all)-1)
		}
	}
	// Closest match.
	closest := 0
	best := 1 << 31
	for i, value := range all {
		diff := value - size
		if diff < 0 {
			diff = -diff
		}
		if diff < best {
			best = diff
			closest = i
		}
	}
	if len(all) <= 1 {
		return 0
	}
	return float32(closest) / float32(len(all)-1)
}

// RoundedRange snaps a [0,1] slider value to the nearest integer in [low, high].
func RoundedRange(value float32, low, high int) int {
	return min(max(int(RoundHalfAway(float64(Denormalize(value, float32(low), float32(high))))), low), high)
}
