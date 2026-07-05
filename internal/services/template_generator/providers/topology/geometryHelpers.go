package topology

import (
	"math"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
)

// circlePoint returns the point at the given angle (radians) on a circle of the
// given radius centred on (centreX, centreY), in normalized [0,1] layout space.
func circlePoint(angle, centreX, centreY, radius float64) models.Position {
	return data.NewVec2(
		centreX+math.Cos(angle)*radius,
		centreY+math.Sin(angle)*radius)
}

// squarePerimeterPoint maps a parameter t in [0,1) to a point traveling
// clockwise around the perimeter of a square centred on (centreX, centreY) with
// the given half-side, starting from the top-left corner.
func squarePerimeterPoint(t, centreX, centreY, half float64) models.Position {
	t -= math.Floor(t)
	scaled := t * 4.0
	side := int(scaled)
	frac := scaled - float64(side)
	left := centreX - half
	right := centreX + half
	top := centreY - half
	bottom := centreY + half
	switch side {
	case 0: // top edge, left to right
		return data.NewVec2(left+frac*2.0*half, top)
	case 1: // right edge, top to bottom
		return data.NewVec2(right, top+frac*2.0*half)
	case 2: // bottom edge, right to left
		return data.NewVec2(right-frac*2.0*half, bottom)
	default: // left edge, bottom to top
		return data.NewVec2(left, bottom-frac*2.0*half)
	}
}

// nearestIndexInRange returns the index in [lo, hi) whose position is closest to
// the position at fromIndex, or -1 when the range is empty.
func nearestIndexInRange(positions models.Positions, fromIndex, lo, hi int) int {
	best := -1
	bestDistance := math.MaxFloat64
	for i := lo; i < hi; i++ {
		if i == fromIndex {
			continue
		}
		deltaX := positions[fromIndex].X - positions[i].X
		deltaY := positions[fromIndex].Y - positions[i].Y
		distance := deltaX*deltaX + deltaY*deltaY
		if distance < bestDistance {
			bestDistance = distance
			best = i
		}
	}
	return best
}

// pairBuilder accumulates unique, unordered index pairs that drive a topology's
// connections. Duplicate and self pairs are ignored so the reused random
// topology connection helpers never emit duplicate connection names.
type pairBuilder struct {
	pairs []models.ConnectionIndexes
	seen  map[models.ConnectionIndexes]bool
}

func newPairBuilder() *pairBuilder {
	return &pairBuilder{seen: map[models.ConnectionIndexes]bool{}}
}

func (this *pairBuilder) add(indexA, indexB int) {
	if indexA == indexB {
		return
	}
	if indexA > indexB {
		indexA, indexB = indexB, indexA
	}
	key := data.NewVec2(indexA, indexB)
	if this.seen[key] {
		return
	}
	this.seen[key] = true
	this.pairs = append(this.pairs, key)
}
