package utils

import "math"

func ComputeContentScale(mapSize, totalZones int) float64 {
	const referenceArea = 160.0 * 160.0 / 4.0
	zoneArea := float64(mapSize) * float64(mapSize) / math.Max(1, float64(totalZones))
	return math.Max(0.5, math.Min(2.5, math.Sqrt(zoneArea/referenceArea)))
}
