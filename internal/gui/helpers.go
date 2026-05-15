package gui

import (
	"math"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
)

// mapRange linearly maps a [0,1] slider value into [low, high].
func mapRange(value, low, high float32) float32 {
	if value < 0 {
		value = 0
	}
	if value > 1 {
		value = 1
	}
	return low + value*(high-low)
}

// mapRangeInv is the inverse of mapRange: maps a value in [low, high]
// back to its [0,1] slider position.
func mapRangeInv(value, low, high float32) float32 {
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

// indexOf returns the index of value in items, or -1 when not present.
func indexOf[T comparable](items []T, value T) int {
	for i, candidate := range items {
		if candidate == value {
			return i
		}
	}
	return -1
}

// mapSizeLabelInt returns the short S/M/L/XL/H/G/C label for an integer size.
func mapSizeLabelInt(size int) string {
	switch {
	case size == 64:
		return "S"
	case size == 80 || size == 96:
		return "M"
	case size == 112 || size == 128:
		return "L"
	case size == 144 || size == 160:
		return "XL"
	case size == 176 || size == 192:
		return "H"
	case size >= 208 && size <= 256:
		return "G"
	default:
		return "C"
	}
}

// roundedRange snaps a [0,1] slider value to the nearest integer in [low, high].
func roundedRange(value float32, low, high int) int {
	return min(max(int(roundHalfAway(float64(mapRange(value, float32(low), float32(high))))), low), high)
}

// mapSizeToSlider returns the [0,1] slider position for a map size value.
func mapSizeToSlider(size int, includeExp bool) float32 {
	all := mapSizes
	if includeExp {
		all = append(append([]int{}, mapSizes...), expMapSizes...)
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

func sliderToMapSize(value float32, includeExp bool) int {
	all := mapSizes
	if includeExp {
		all = append(append([]int{}, mapSizes...), expMapSizes...)
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

func topologyLabelFor(topology models.MapTopology) string {
	for i, value := range topologyValues {
		if value == topology {
			return topologyLabels[i]
		}
	}
	return topologyLabels[0]
}

func victoryIndex(id string) int {
	for i, value := range victoryIDs {
		if value == id {
			return i
		}
	}
	return 0
}
