package helpers

import (
	"image"
	"math"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
)

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

// CalculatePointTowards returns `from` moved toward `toward` by `distance` pixels.
// ok is false when the two points are coincident.
func CalculatePointTowards(
	from, toward image.Point,
	distance float64) (movedPoint image.Point, ok bool) {
	delta := data.Vec2FromPoint[float64](toward.Sub(from))
	deltaDistance := math.Hypot(delta.X, delta.Y)
	if deltaDistance < 1 {
		return movedPoint, false
	}

	movedPoint = data.Vec2FromPoint[float64](from).
		Add(delta.DivideScalar(deltaDistance).MultiplyScalar(distance)).
		ToPoint()
	return movedPoint, true
}
