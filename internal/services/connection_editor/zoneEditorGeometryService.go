package connection_editor

import (
	"math"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/preview_service"
)

// Curve tunables - adjust the shape of the connection curves here.
const (
	// parallelEdgeGapPx is the perpendicular distance between two curves that
	// connect the same pair of zones.
	parallelEdgeGapPx = 18.0
	// obstacleClearancePx is how far beyond a zone's radius a curve must pass
	// before that zone stops pushing it aside.
	obstacleClearancePx = 8.0
	// obstacleBulgePaddingPx is the extra push applied on top of the clearance
	// so a deflected curve does not merely graze the obstacle.
	obstacleBulgePaddingPx = 6.0
	// obstacleChordMargin ignores obstacles sitting near either end of the
	// chord, where a bulge would only distort the curve's attachment.
	obstacleChordMargin = 0.08
	// edgeHitReachPx is how close a point must be to a curve to hit it.
	edgeHitReachPx = 9.0
	// edgeHitSampleCount is how many points along a curve are distance-tested.
	edgeHitSampleCount = 21
)

// Snapping tunables - adjust the feel of the snap toggle here.
const (
	// gridCellsPerZoneDiameter sets the snap-grid density: the distance
	// between adjacent grid lines is (zone diameter) / gridCellsPerZoneDiameter.
	gridCellsPerZoneDiameter = 7.0
	// gridSnapThresholdPx is the "light" hold distance (canvas px) within
	// which a dragged zone's edges/center stick to a grid line.
	gridSnapThresholdPx = 4.0
	// zoneSnapThresholdPx is the "heavier" hold distance (canvas px) within
	// which a dragged zone's edges/center stick to the horizontal or vertical
	// extension of another zone's edge or center.
	zoneSnapThresholdPx = 9.0
)

// ZoneEditorGeometryService computes the manual zone editor's canvas geometry.
// Node placement is delegated to the preview layout service, so the editor and
// the preview tab always agree on where a zone sits.
type ZoneEditorGeometryService struct {
	previewLayout preview_service.IPreviewLayoutService
}

func NewZoneEditorGeometryService(
	previewLayout preview_service.IPreviewLayoutService) IZoneEditorGeometryService {
	return &ZoneEditorGeometryService{previewLayout: previewLayout}
}

func (this *ZoneEditorGeometryService) BuildGeometry(
	zones []template_model.Zone,
	connections []template_model.Connection,
	topology config.MapTopology,
	canvasSide int) models.ZoneEditorGeometry {
	template := &template_model.Template{
		Variants: []template_model.Variant{{
			Zones:       zones,
			Connections: connections,
		}},
	}
	layout := this.previewLayout.BuildPreviewLayout(template, topology, float64(canvasSide))
	return models.ZoneEditorGeometry{
		Positions:  layout.Positions,
		Zones:      layout.Zones,
		ZoneRadius: layout.ZoneRadius,
		Edges:      buildEdges(connections, layout.Positions, layout.ZoneRadius),
	}
}

func (this *ZoneEditorGeometryService) HitTestNode(
	position models.Position,
	positions map[string]models.Position,
	zoneRadius float64) string {
	best := ""
	bestDistance := math.MaxFloat64
	for name, center := range positions {
		distance := position.Subtract(center).Distance()
		if distance <= zoneRadius && distance < bestDistance {
			bestDistance = distance
			best = name
		}
	}
	return best
}

func (this *ZoneEditorGeometryService) HitTestEdge(
	position models.Position,
	edges []models.ZoneEditorEdge) int {
	best := -1
	bestDistance := edgeHitReachPx
	for index := range edges {
		edge := edges[index]
		for step := range edgeHitSampleCount {
			ratio := float64(step) / float64(edgeHitSampleCount-1)
			bezierPoint := helpers.GetVectorOnQuadraticBezierCurve(
				edge.StartPoint,
				edge.ControlPoint,
				edge.EndPoint,
				ratio)
			distance := position.Subtract(bezierPoint).Distance()
			if distance < bestDistance {
				bestDistance = distance
				best = index
			}
		}
	}
	return best
}

func (this *ZoneEditorGeometryService) GridStep(zoneRadius float64) float64 {
	return zoneRadius * 2.0 / gridCellsPerZoneDiameter
}

func (this *ZoneEditorGeometryService) SnapPosition(
	position models.Position,
	positions map[string]models.Position,
	zoneRadius float64,
	draggedZone string) models.ZoneEditorSnapResult {
	if zoneRadius <= 0 {
		return models.ZoneEditorSnapResult{Position: position}
	}

	// The dragged zone's own snap points on each axis: leading edge, center,
	// trailing edge.
	offsets := [3]float64{-zoneRadius, 0, zoneRadius} // Probably can just be a Vec3
	guidesX, guidesY := otherZoneGuides(positions, draggedZone, zoneRadius)
	step := this.GridStep(zoneRadius)
	x, guideX, hitX := snapAxis(position.X, offsets, guidesX, step)
	y, guideY, hitY := snapAxis(position.Y, offsets, guidesY, step)
	result := models.ZoneEditorSnapResult{Position: data.NewVec2(x, y)}
	if hitX {
		result.GuideX, result.HasGuideX = guideX, true
	}
	if hitY {
		result.GuideY, result.HasGuideY = guideY, true
	}
	return result
}

// buildEdges turns every connection whose two endpoints have a position into a
// curve, spreading connections that share a zone pair symmetrically around the
// straight chord and bending clear of intermediate nodes.
func buildEdges(
	connections []template_model.Connection,
	positions map[string]models.Position,
	zoneRadius float64) []models.ZoneEditorEdge {
	order, groups := groupConnectionsByPair(connections)
	edges := make([]models.ZoneEditorEdge, 0, len(connections))
	for _, key := range order {
		group := groups[key]
		count := len(group)
		for slot, connectionIndex := range group {
			connection := connections[connectionIndex]
			startPoint, hasStart := positions[connection.From]
			endPoint, hasEnd := positions[connection.To]
			if !hasStart || !hasEnd {
				continue
			}
			// Bend around the canonical (lexicographic) endpoint order so that
			// A->B and B->A spread to opposite sides of the same chord.
			canonicalA, canonicalB := startPoint, endPoint
			if connection.From > connection.To {
				canonicalA, canonicalB = canonicalB, canonicalA
			}
			delta := canonicalB.Subtract(canonicalA)
			distance := delta.Distance()
			if distance < 1 {
				distance = 1
			}
			normal := delta.RotateClockwise().DivideScalar(distance)
			spread := (float64(slot) - float64(count-1)/2.0) * parallelEdgeGapPx
			bulge := spread + obstacleBulge(positions, zoneRadius, canonicalA, canonicalB, normal)
			midPoint := startPoint.Add(endPoint).MultiplyScalar(0.5)
			controlPoint := midPoint.Add(normal.MultiplyScalar(2.0 * bulge))
			labelPoint := startPoint.MultiplyScalar(0.25).
				Add(controlPoint.MultiplyScalar(0.5)).
				Add(endPoint.MultiplyScalar(0.25))
			edges = append(edges, models.ZoneEditorEdge{
				ConnectionIndex: connectionIndex,
				StartPoint:      startPoint,
				EndPoint:        endPoint,
				ControlPoint:    controlPoint,
				MidPoint:        labelPoint,
			})
		}
	}
	return edges
}

// groupConnectionsByPair buckets connection indices by unordered endpoint pair,
// preserving first-seen order so parallel edges spread deterministically from
// frame to frame.
func groupConnectionsByPair(
	connections []template_model.Connection) ([]connectionPairKey, map[connectionPairKey][]int) {
	groups := make(map[connectionPairKey][]int)
	order := make([]connectionPairKey, 0)
	for index, connection := range connections {
		from, to := connection.From, connection.To
		if from > to {
			from, to = to, from
		}
		key := connectionPairKey{from: from, to: to}
		if _, seen := groups[key]; !seen {
			order = append(order, key)
		}
		groups[key] = append(groups[key], index)
	}
	return order, groups
}

// obstacleBulge returns a perpendicular push so a curve bends clear of any zone
// node that lies close to the straight chord between its two endpoints.
func obstacleBulge(
	positions map[string]models.Position,
	zoneRadius float64,
	chordStart, chordEnd models.Position,
	normal models.Position) float64 {
	clearance := zoneRadius + obstacleClearancePx
	segment := chordEnd.Subtract(chordStart)
	segmentLengthSquared := segment.SquaredLength()
	if segmentLengthSquared < 1 {
		return 0
	}

	best := 0.0
	bestMagnitude := 0.0
	for _, center := range positions {
		ratio := center.Subtract(chordStart).DotProduct(segment) / segmentLengthSquared
		if ratio <= obstacleChordMargin || ratio >= 1-obstacleChordMargin {
			continue
		}

		offset := center.Subtract(chordStart.Add(segment.MultiplyScalar(ratio)))
		perpendicular := offset.Distance()
		if perpendicular >= clearance {
			continue
		}

		need := (clearance - perpendicular) + obstacleBulgePaddingPx
		signed := need
		if offset.DotProduct(normal) >= 0 {
			signed = -need
		}
		if math.Abs(signed) > bestMagnitude {
			bestMagnitude = math.Abs(signed)
			best = signed
		}
	}
	return best
}

// otherZoneGuides collects the horizontal and vertical guide coordinates
// (edge / center / edge) of every zone except the dragged one.
func otherZoneGuides(
	positions map[string]models.Position,
	draggedZone string,
	radius float64,
) (guidesX, guidesY []float64) { // This needs to return data.Positions, but for not I'm leaving as is to not allocate more arrays in snapAxis
	for name, center := range positions {
		if name == draggedZone {
			continue
		}
		guidesX = append(guidesX, center.X-radius, center.X, center.X+radius)
		guidesY = append(guidesY, center.Y-radius, center.Y, center.Y+radius)
	}
	return guidesX, guidesY
}

// snapAxis snaps a single axis value. Zone-alignment guides win over the grid;
// within each class the smallest correction wins. When a zone guide is hit its
// coordinate is returned so the caller can draw an alignment indicator.
func snapAxis(
	value float64,
	offsets [3]float64, // Probably can just be a Vec3
	guides []float64,
	gridStep float64,
) (snapped float64, guide float64, zoneGuideHit bool) {
	best := math.MaxFloat64
	bestGuide := 0.0
	for _, offset := range offsets {
		point := value + offset
		for _, candidate := range guides {
			if delta := candidate - point; math.Abs(delta) < math.Abs(best) {
				best = delta
				bestGuide = candidate
			}
		}
	}
	if math.Abs(best) <= zoneSnapThresholdPx {
		return value + best, bestGuide, true
	}

	if gridStep > 0 {
		best = math.MaxFloat64
		for _, offset := range offsets {
			point := value + offset
			if delta := math.Round(point/gridStep)*gridStep - point; math.Abs(delta) < math.Abs(best) {
				best = delta
			}
		}
		if math.Abs(best) <= gridSnapThresholdPx {
			return value + best, 0, false
		}
	}

	return value, 0, false
}
