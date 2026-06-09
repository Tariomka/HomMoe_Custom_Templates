package utils

import (
	"math"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
)

func NormalizeZoneSize(zoneSize float64) float64 {
	if math.IsNaN(zoneSize) || math.IsInf(zoneSize, 0) {
		return 1.0
	}

	return helpers.RoundWithPrecision(math.Max(0.1, math.Min(zoneSize, 2.0)), 2)
}
