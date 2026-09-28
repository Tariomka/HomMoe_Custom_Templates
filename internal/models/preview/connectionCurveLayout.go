package preview

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
)

const (
	parallelCurveGapPx     = 21.0 // perpendicular distance between two curves that connect the same pair of zones.
	obstacleClearancePx    = 8.0  // how far beyond a zone's radius a curve must pass before zone stops pushing it aside.
	obstacleBulgePaddingPx = 6.0  // extra push applied on top of the clearance so a deflected curve does not merely graze the obstacle.
	obstacleChordMargin    = 0.08 // ignores obstacles sitting near either end of the chord.
)

type ConnectionCurveLayout struct {
	Positions  map[string]data.Vec2[float64]
	ZoneRadius float64
}

// Build returns one curve per connection whose two endpoints are positioned, grouped by zone pair in first-seen order.
func (this ConnectionCurveLayout) Build(connections []template_model.Connection) []ConnectionCurve {
	order, groups := groupConnectionsByPair(connections)
	curves := make([]ConnectionCurve, 0, len(connections))
	for _, key := range order {
		start, hasStart := this.Positions[key[0]]
		end, hasEnd := this.Positions[key[1]]
		if !hasStart || !hasEnd {
			continue
		}

		delta := end.Subtract(start)
		normal := delta.RotateClockwise().DivideScalar(max(delta.Distance(), 1))
		midPoint := start.Add(end).MultiplyScalar(0.5)
		positiveNeed, negativeNeed, obstructed := this.obstacleDetours(start, end, normal)
		group := groups[key]
		for slot, connectionIndex := range group {
			bulge := slotBulge(slot, len(group), positiveNeed, negativeNeed, obstructed)
			// A quadratic Bézier's midpoint sits halfway between the chord midpoint
			// and the control point, hence the doubled offset.
			curves = append(curves, ConnectionCurve{
				ConnectionIndex: connectionIndex,
				Start:           start,
				End:             end,
				Control:         midPoint.Add(normal.MultiplyScalar(2.0 * bulge)),
			})
		}
	}
	return curves
}

// obstacleDetours returns the smallest midpoint offset along +normal and along
// -normal that passes every zone lying close to the chord, and whether any zone
// does. Max/min are order-independent, so map iteration cannot change them.
func (this ConnectionCurveLayout) obstacleDetours(
	chordStart, chordEnd data.Vec2[float64],
	normal data.Vec2[float64]) (positiveNeed, negativeNeed float64, obstructed bool) {
	segment := chordEnd.Subtract(chordStart)
	segmentLengthSquared := segment.SquaredLength()
	if segmentLengthSquared < 1 {
		return 0, 0, false
	}

	clearance := this.ZoneRadius + obstacleClearancePx
	reach := clearance + obstacleBulgePaddingPx
	for _, center := range this.Positions {
		ratio := center.Subtract(chordStart).DotProduct(segment) / segmentLengthSquared
		if ratio <= obstacleChordMargin || ratio >= 1-obstacleChordMargin {
			continue
		}

		offset := center.Subtract(chordStart.Add(segment.MultiplyScalar(ratio)))
		if offset.Distance() >= clearance {
			continue
		}

		signed := offset.DotProduct(normal)
		if !obstructed {
			positiveNeed, negativeNeed, obstructed = signed+reach, signed-reach, true
			continue
		}
		positiveNeed = max(positiveNeed, signed+reach)
		negativeNeed = min(negativeNeed, signed-reach)
	}
	return positiveNeed, negativeNeed, obstructed
}

// slotBulge is the midpoint offset of the slot-th of count curves sharing a pair.
func slotBulge(slot, count int, positiveNeed, negativeNeed float64, obstructed bool) float64 {
	if !obstructed {
		return (float64(slot) - float64(count-1)/2.0) * parallelCurveGapPx
	}

	negativeCount := count / 2
	if count%2 == 1 && -negativeNeed < positiveNeed {
		negativeCount++
	}
	if slot < negativeCount {
		return negativeNeed - float64(negativeCount-1-slot)*parallelCurveGapPx
	}
	return positiveNeed + float64(slot-negativeCount)*parallelCurveGapPx
}

// groupConnectionsByPair buckets connection indices by unordered endpoint pair,
// keyed in canonical (lexicographic) order and preserving first-seen order.
func groupConnectionsByPair(
	connections []template_model.Connection) ([][2]string, map[[2]string][]int) {
	groups := make(map[[2]string][]int)
	order := make([][2]string, 0)
	for index, connection := range connections {
		key := [2]string{connection.From, connection.To}
		if key[0] > key[1] {
			key[0], key[1] = key[1], key[0]
		}
		if _, seen := groups[key]; !seen {
			order = append(order, key)
		}
		groups[key] = append(groups[key], index)
	}
	return order, groups
}
