package preview_service

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
)

// layoutFixedPositions places zones at their exact GeneratorPosition stamps,
// preserving the deterministic geometric figure built by the Square, Geometric,
// Cross and Fractal topologies. The normalized positions are centered and
// uniformly scaled to fill the padded canvas (never relaxed), then the zone
// radius is shrunk just enough to keep the closest pair from overlapping.
func (this *PreviewLayoutService) layoutFixedPositions(zones []entities.Zone, side float64) {
	metrics := newCanvasMetrics(side)
	if this.placeTrivial(zones, metrics) {
		return
	}

	px, py := generatorCoords(zones)
	// The margin reserves room for the largest possible zone radius, so the
	// eventual radius (never larger) always fits.
	fitToCanvas(px, py, metrics, metrics.margin+metrics.zoneRadiusMax, true)
	radius := radiusFromClosestPair(px, py, metrics.zoneRadiusMax, metrics.minGap)
	this.commitPositions(zones, px, py, radius)
}
