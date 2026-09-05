package connection_editor

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
)

// IZoneEditorGeometryService is the contract of the manual zone editor's pure
// canvas geometry: node placement, connection curves, hit testing and drag
// snapping. It computes coordinates only and knows nothing about rendering.
type IZoneEditorGeometryService interface {
	// BuildGeometry lays the zones out into a square canvas of the given side
	// and curves every connection between them, spreading connections that
	// share a zone pair and bending clear of intermediate nodes.
	BuildGeometry(
		zones []template_model.Zone,
		connections []template_model.Connection,
		topology config.MapTopology,
		canvasSide int) models.ZoneEditorGeometry

	// HitTestNode returns the name of the zone whose node covers position, or
	// "" when no node does.
	HitTestNode(position models.Position, positions map[string]models.Position, zoneRadius float64) string

	// HitTestEdge returns the index of the edge whose curve passes closest to
	// position, or -1 when no curve is within reach.
	HitTestEdge(position models.Position, edges []models.ZoneEditorEdge) int

	// GridStep returns the snapping-grid cell size in canvas pixels.
	GridStep(zoneRadius float64) float64

	// SnapPosition nudges a dragged zone's center so its edges or center hold
	// onto nearby alignment guides and grid lines. It never pulls from afar:
	// only a position already within the hold distance sticks.
	SnapPosition(
		position models.Position,
		positions map[string]models.Position,
		zoneRadius float64,
		draggedZone string) models.ZoneEditorSnapResult
}
