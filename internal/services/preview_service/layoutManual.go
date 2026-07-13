package preview_service

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
)

// layoutManualPositions places zones exactly where the manual zone editor put
// them: canvas = normalized position × side. The mapping must stay trivially
// invertible (p = pos / side) so dragging in the editor is exact. The zone
// radius shrinks just enough to keep the closest pair of zones from
// overlapping.
func (this *PreviewLayoutService) layoutManualPositions(zones []entities.Zone, side float64) {
	metrics := newCanvasMetrics(side)
	px := make([]float64, len(zones))
	py := make([]float64, len(zones))
	for i, zone := range zones {
		p := *zone.ManualPosition
		px[i] = p[0] * side
		py[i] = p[1] * side
	}
	radius := radiusFromClosestPair(px, py, metrics.zoneRadiusMax, metrics.minGap)
	this.commitPositions(zones, px, py, radius)
}
