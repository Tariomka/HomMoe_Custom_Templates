package connection_editor

import (
	"image"
	"math"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
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
	zones []entities.Zone,
	connections []entities.Connection,
	topology config.MapTopology,
	canvasSide int) models.ZoneEditorGeometry {
	template := &entities.RmgTemplate{
		Variants: []entities.Variant{variant_content.NewVariantBuilder().
			WithZones(zones...).
			WithConnections(connections...).
			Build()},
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
	position image.Point,
	positions map[string]image.Point,
	zoneRadius int) string {
	best := ""
	bestDistance := math.MaxFloat64
	reach := float64(zoneRadius)
	for name, center := range positions {
		distance := math.Hypot(float64(position.X-center.X), float64(position.Y-center.Y))
		if distance <= reach && distance < bestDistance {
			bestDistance = distance
			best = name
		}
	}
	return best
}

func (this *ZoneEditorGeometryService) HitTestEdge(
	position image.Point,
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
			distance := math.Hypot(float64(position.X)-bezierPoint.X, float64(position.Y)-bezierPoint.Y)
			if distance < bestDistance {
				bestDistance = distance
				best = index
			}
		}
	}
	return best
}

func (this *ZoneEditorGeometryService) GridStep(zoneRadius int) float64 {
	return float64(zoneRadius) * 2.0 / gridCellsPerZoneDiameter
}

func (this *ZoneEditorGeometryService) SnapPosition(
	position image.Point,
	positions map[string]image.Point,
	zoneRadius int,
	draggedZone string) models.ZoneEditorSnapResult {
	if zoneRadius <= 0 {
		return models.ZoneEditorSnapResult{Position: position}
	}
	radius := float64(zoneRadius)
	// The dragged zone's own snap points on each axis: leading edge, center,
	// trailing edge.
	offsets := [3]float64{-radius, 0, radius}
	guidesX, guidesY := otherZoneGuides(positions, draggedZone, radius)
	step := this.GridStep(zoneRadius)
	x, guideX, hitX := snapAxis(float64(position.X), offsets, guidesX, step)
	y, guideY, hitY := snapAxis(float64(position.Y), offsets, guidesY, step)
	result := models.ZoneEditorSnapResult{
		Position: image.Pt(int(math.Round(x)), int(math.Round(y))),
	}
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
	connections []entities.Connection,
	positions map[string]image.Point,
	zoneRadius int) []models.ZoneEditorEdge {
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
			deltaX := float64(canonicalB.X - canonicalA.X)
			deltaY := float64(canonicalB.Y - canonicalA.Y)
			distance := math.Hypot(deltaX, deltaY)
			if distance < 1 {
				distance = 1
			}
			normalX := deltaY / distance
			normalY := -deltaX / distance
			spread := (float64(slot) - float64(count-1)/2.0) * parallelEdgeGapPx
			bulge := spread + obstacleBulge(positions, zoneRadius, canonicalA, canonicalB, normalX, normalY)
			midX := float64(startPoint.X+endPoint.X) / 2.0
			midY := float64(startPoint.Y+endPoint.Y) / 2.0
			controlX := midX + 2.0*bulge*normalX
			controlY := midY + 2.0*bulge*normalY
			labelX := 0.25*float64(startPoint.X) + 0.5*controlX + 0.25*float64(endPoint.X)
			labelY := 0.25*float64(startPoint.Y) + 0.5*controlY + 0.25*float64(endPoint.Y)
			edges = append(edges, models.ZoneEditorEdge{
				ConnectionIndex: connectionIndex,
				StartPoint:      data.NewVec2(float64(startPoint.X), float64(startPoint.Y)),
				EndPoint:        data.NewVec2(float64(endPoint.X), float64(endPoint.Y)),
				ControlPoint:    data.NewVec2(controlX, controlY),
				MidPoint:        image.Pt(int(labelX), int(labelY)),
			})
		}
	}
	return edges
}

// groupConnectionsByPair buckets connection indices by unordered endpoint pair,
// preserving first-seen order so parallel edges spread deterministically from
// frame to frame.
func groupConnectionsByPair(
	connections []entities.Connection) ([]connectionPairKey, map[connectionPairKey][]int) {
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
	positions map[string]image.Point,
	zoneRadius int,
	chordStart, chordEnd image.Point,
	normalX, normalY float64) float64 {
	clearance := float64(zoneRadius) + obstacleClearancePx
	startX, startY := float64(chordStart.X), float64(chordStart.Y)
	segmentX := float64(chordEnd.X - chordStart.X)
	segmentY := float64(chordEnd.Y - chordStart.Y)
	segmentLengthSquared := segmentX*segmentX + segmentY*segmentY
	if segmentLengthSquared < 1 {
		return 0
	}
	best := 0.0
	bestMagnitude := 0.0
	for _, center := range positions {
		centerX, centerY := float64(center.X), float64(center.Y)
		ratio := ((centerX-startX)*segmentX + (centerY-startY)*segmentY) / segmentLengthSquared
		if ratio <= obstacleChordMargin || ratio >= 1-obstacleChordMargin {
			continue
		}
		closestX := startX + ratio*segmentX
		closestY := startY + ratio*segmentY
		perpendicular := math.Hypot(centerX-closestX, centerY-closestY)
		if perpendicular >= clearance {
			continue
		}
		side := (centerX-closestX)*normalX + (centerY-closestY)*normalY
		need := (clearance - perpendicular) + obstacleBulgePaddingPx
		signed := need
		if side >= 0 {
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
	positions map[string]image.Point,
	draggedZone string,
	radius float64) (guidesX, guidesY []float64) {
	for name, center := range positions {
		if name == draggedZone {
			continue
		}
		centerX, centerY := float64(center.X), float64(center.Y)
		guidesX = append(guidesX, centerX-radius, centerX, centerX+radius)
		guidesY = append(guidesY, centerY-radius, centerY, centerY+radius)
	}
	return guidesX, guidesY
}

// snapAxis snaps a single axis value. Zone-alignment guides win over the grid;
// within each class the smallest correction wins. When a zone guide is hit its
// coordinate is returned so the caller can draw an alignment indicator.
func snapAxis(
	value float64,
	offsets [3]float64,
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
