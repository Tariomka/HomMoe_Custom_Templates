package topology

import (
	"math"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
)

// circlePoint returns the point at the given angle (radians) on a circle of the
// given radius centred on (centreX, centreY), in normalized [0,1] layout space.
func circlePoint(angle, centreX, centreY, radius float64) models.Vector2 {
	return models.NewPosition(
		centreX+math.Cos(angle)*radius,
		centreY+math.Sin(angle)*radius)
}

// squarePerimeterPoint maps a parameter t in [0,1) to a point travelling
// clockwise around the perimeter of a square centred on (centreX, centreY) with
// the given half-side, starting from the top-left corner.
func squarePerimeterPoint(t, centreX, centreY, half float64) models.Vector2 {
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
		return models.NewPosition(left+frac*2.0*half, top)
	case 1: // right edge, top to bottom
		return models.NewPosition(right, top+frac*2.0*half)
	case 2: // bottom edge, right to left
		return models.NewPosition(right-frac*2.0*half, bottom)
	default: // left edge, bottom to top
		return models.NewPosition(left, bottom-frac*2.0*half)
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
	pairs [][2]int
	seen  map[[2]int]bool
}

func newPairBuilder() *pairBuilder {
	return &pairBuilder{seen: map[[2]int]bool{}}
}

func (this *pairBuilder) add(indexA, indexB int) {
	if indexA == indexB {
		return
	}
	if indexA > indexB {
		indexA, indexB = indexB, indexA
	}
	key := [2]int{indexA, indexB}
	if this.seen[key] {
		return
	}
	this.seen[key] = true
	this.pairs = append(this.pairs, key)
}
