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
func Clamp[T data.Numeric](value, lowest, highest T) T {
	if value < lowest {
		return lowest
	}

	if value > highest {
		return highest
	}

	return value
}

// Scale returns value scaled by multiplier, clamped to a minimum of 0, disregarding fractional parts.
func Scale(value, multiplier float64) int {
	return max(0, int(value*multiplier))
}

// ScaleRound returns value scaled by multiplier, clamped to a minimum of 0, rounding to the nearest integer.
func ScaleRound(value, multiplier float64) int {
	return max(0, int(math.Round(value*multiplier)))
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
func CalculatePointTowards(from, toward image.Point, distance float64) (movedPoint image.Point, ok bool) {
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

// GetVectorOnQuadraticBezierCurve evaluates the quadratic Bézier (start(P0)→ctrl(P1)→end(P2)) at
// the point along the curve t in [0,1] by applying this formula:
//
// B(t) = (1-t)²*P0 + 2*(1-t)*t*P1 + t²*P2.
func GetVectorOnQuadraticBezierCurve(start, ctrl, end data.Vec2[float64], t float64) data.Vec2[float64] {
	mt := 1 - t
	return start.MultiplyScalar(mt * mt).Add(ctrl.MultiplyScalar(2 * mt * t)).Add(end.MultiplyScalar(t * t))
}

// GetPointOnQuadraticBezierCurve evaluates the quadratic Bézier (start(P0)→ctrl(P1)→end(P2)) at
// the point along the curve t in [0,1] by applying this formula:
//
// B(t) = (1-t)²*P0 + 2*(1-t)*t*P1 + t²*P2.
func GetPointOnQuadraticBezierCurve(start, ctrl, end image.Point, t float64) image.Point {
	return GetVectorOnQuadraticBezierCurve(
		data.Vec2FromPoint[float64](start),
		data.Vec2FromPoint[float64](ctrl),
		data.Vec2FromPoint[float64](end),
		t).ToPointRounded()
}
