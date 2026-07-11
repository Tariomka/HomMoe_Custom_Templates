package misc

import "math"

// GetShortestAngleDistance returns the signed shortest distance from angle 'from' to 'to' in radians.
// The result is always between ±[math.Pi].
func GetShortestAngleDistance(from, to float64) float64 {
	const maxTurn = math.Pi * 2
	// Handle wrapping using math.Mod
	delta := math.Mod(math.Abs(to-from), maxTurn)
	if delta > math.Pi {
		delta = maxTurn - delta
	}
	return delta
}
