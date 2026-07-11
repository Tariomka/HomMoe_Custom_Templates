package preview_service

import (
	"image"
	"math"
	"sort"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
)

// layoutBalancedRings renders the Circles topology as concentric rings keyed
// off the zones' GeneratorRing stamps. The zone radius is binary-searched so
// the outermost ring still fits the canvas, and each ring keeps the zones'
// original angular neighbour ordering.
func (this *PreviewLayoutService) layoutBalancedRings(zones []entities.Zone, side float64) {
	metrics := newCanvasMetrics(side)
	if this.placeTrivial(zones, metrics) {
		return
	}

	rings := groupZonesByRing(zones)
	if len(rings) < 2 {
		// All zones in a single ring - degenerate; fall back to the ring path.
		this.layoutRingOrHub(zones, nil, side)
		return
	}

	drawRadius := side/2.0 - metrics.margin - metrics.zoneRadiusMax
	zoneRadius := fitRingZoneRadius(rings, drawRadius, metrics.minGap, metrics.zoneRadiusMax)
	ringRadii := assignRingRadii(rings, zoneRadius, drawRadius, metrics.minGap)
	this.layout.ZoneRadius = int(math.Round(zoneRadius))
	this.placeRings(zones, rings, ringRadii, metrics)
}

// groupZonesByRing buckets zone indices by their GeneratorRing tier. Ring
// index 0 holds the largest tier value present; higher indices move outwards
// through decreasing tiers. Callers must have verified allHaveRing.
func groupZonesByRing(zones []entities.Zone) [][]int {
	presentSet := map[int]bool{}
	for _, zone := range zones {
		presentSet[*zone.GeneratorRing] = true
	}
	present := make([]int, 0, len(presentSet))
	for tier := range presentSet {
		present = append(present, tier)
	}
	sort.Ints(present)

	ringCount := len(present)
	tierToRing := make(map[int]int, ringCount)
	for i, tier := range present {
		tierToRing[tier] = ringCount - 1 - i
	}

	rings := make([][]int, ringCount)
	for i, zone := range zones {
		ring := tierToRing[*zone.GeneratorRing]
		rings[ring] = append(rings[ring], i)
	}
	return rings
}

// assignRingRadii computes each ring's canvas radius for the given zone
// radius: at least its natural share of the draw radius, wide enough that the
// zones within the ring keep minGap of clearance, and clear of the previous
// ring.
func assignRingRadii(rings [][]int, zoneRadius, drawRadius, minGap float64) []float64 {
	minClearance := 2.0*zoneRadius + minGap
	radii := make([]float64, len(rings))
	for ringIndex := range rings {
		count := len(rings[ringIndex])
		natural := drawRadius * float64(ringIndex+1) / float64(len(rings))
		withinRing := 0.0
		if count >= 2 {
			withinRing = minClearance / (2.0 * math.Sin(math.Pi/float64(count)))
		} else if count == 1 && ringIndex > 0 {
			withinRing = minClearance
		}
		afterPrev := 0.0
		if ringIndex > 0 {
			afterPrev = radii[ringIndex-1] + minClearance
		}
		radii[ringIndex] = math.Max(natural, math.Max(withinRing, afterPrev))
	}
	return radii
}

// fitRingZoneRadius binary-searches the largest zone radius that keeps the
// outermost ring inside the available draw radius.
func fitRingZoneRadius(rings [][]int, drawRadius, minGap, maxRadius float64) float64 {
	lo, hi := csZoneRadiusFloor, maxRadius
	for range 32 {
		mid := (lo + hi) / 2.0
		radii := assignRingRadii(rings, mid, drawRadius, minGap)
		if radii[len(radii)-1] <= drawRadius {
			lo = mid
		} else {
			hi = mid
		}
	}
	return math.Max(lo, csZoneRadiusFloor)
}

// placeRings distributes each ring's zones evenly around its canvas radius,
// anchored at the first zone's raw angle so the rendered ring preserves the
// original neighbour ordering.
func (this *PreviewLayoutService) placeRings(
	zones []entities.Zone,
	rings [][]int,
	ringRadii []float64,
	metrics canvasMetrics) {
	rawCx, rawCy := positionCentroid(zones)
	for ringIndex, group := range rings {
		count := len(group)
		if count == 0 {
			continue
		}
		if count == 1 && ringIndex == 0 {
			this.layout.Positions[zones[group[0]].Name] = metrics.centre()
			continue
		}
		sorted := sortIndicesByAngle(zones, group, rawCx, rawCy)
		firstAngle := positionAngle(zones[sorted[0]], rawCx, rawCy)
		canvasRadius := ringRadii[ringIndex]
		for j, zoneIndex := range sorted {
			angle := firstAngle + 2.0*math.Pi*float64(j)/float64(count)
			x := metrics.cx + math.Cos(angle)*canvasRadius
			y := metrics.cy + math.Sin(angle)*canvasRadius
			this.layout.Positions[zones[zoneIndex].Name] = image.Pt(int(math.Round(x)), int(math.Round(y)))
		}
	}
}
