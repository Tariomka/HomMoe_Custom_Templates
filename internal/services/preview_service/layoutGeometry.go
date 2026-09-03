package preview_service

import (
	"math"
	"sort"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/registry"
)

const (
	csCanvasSide      = 700.0
	csMargin          = 18.0
	csZoneRadiusMax   = 38.0
	csZoneRadiusFloor = 8.0
	csHubRadiusMin    = 28.0
	csMinGap          = 6.0
	csConnectionGap   = 26.0 // ring layout - visible chord clearance between zones
	// csGeoHubEdgeInset pulls the Geometric Hub figure away from the canvas
	// border: the fill-fit would otherwise push the player zones right up to
	// the padded edge, hiding the hub-centric shape of the layout.
	csGeoHubEdgeInset = 48.0
	// csGeoHubEdgeInsetCrowded replaces csGeoHubEdgeInset once
	// csGeoHubCrowdedMinPlayers or more players share the figure: the smaller
	// inset lets the fill-fit scale the figure further out, spacing the
	// player zones away from the central hub.
	csGeoHubEdgeInsetCrowded  = 12.0
	csGeoHubCrowdedMinPlayers = 6
	scatterIdealMultiplier    = 3.2
	scatterMinDist            = 3.8
	scatterEdgeClear          = 1.2
)

// canvasMetrics bundles the canvas-side-scaled layout distances every
// renderer needs, replacing the identical prologue each of them used to
// repeat. All values are in canvas pixels for the given side length.
type canvasMetrics struct {
	side          float64
	scale         float64
	margin        float64
	minGap        float64
	zoneRadiusMax float64
	hubRadiusMin  float64
	connectionGap float64
	// cx and cy are useless, they need to be data.Vec2[float64] 'center', removing the need for center() method
	cx float64
	cy float64
}

func newCanvasMetrics(side float64) canvasMetrics {
	scale := side / csCanvasSide
	return canvasMetrics{
		side:          side,
		scale:         scale,
		margin:        csMargin * scale,
		minGap:        csMinGap * scale,
		zoneRadiusMax: csZoneRadiusMax * scale,
		hubRadiusMin:  csHubRadiusMin * scale,
		connectionGap: csConnectionGap * scale,
		cx:            side / 2.0,
		cy:            side / 2.0,
	}
}

func (this canvasMetrics) center() data.Vec2[float64] {
	return data.NewVec2(this.cx, this.cy)
}

// placeTrivial handles the degenerate 0- and 1-zone cases shared by every
// renderer: the zone radius maxes out and a lone zone sits at the canvas
// center. Reports whether the layout was fully handled.
func (this *PreviewLayoutService) placeTrivial(zones []template_model.Zone, metrics canvasMetrics) bool {
	if len(zones) > 1 {
		return false
	}
	this.layout.ZoneRadius = metrics.zoneRadiusMax
	if len(zones) == 1 {
		this.layout.Positions[zones[0].Name] = metrics.center()
	}
	return true
}

// commitPositions stores the working coordinates and the final zone radius in the layout.
func (this *PreviewLayoutService) commitPositions(
	zones []template_model.Zone,
	positions models.Positions,
	radius float64) {
	this.layout.ZoneRadius = radius
	for i, zone := range zones {
		position := positions[i]
		this.layout.Positions[zone.Name] = position
	}
}

// getGeneratorCoordinates copies the zones' GeneratorPosition stamps into mutable
// coordinate slices. Callers must have verified allHavePosition.
func getGeneratorCoordinates(zones []template_model.Zone) models.Positions {
	var positions models.Positions
	for _, zone := range zones {
		p := *zone.GeneratorPosition // Is this required to be copied? can't it be used directly safely?
		positions.Add(data.NewVec2(p[0], p[1]))
	}
	return positions
}

// radiusFromClosestPair shrinks maxRadius just enough that the closest pair
// of zones keeps minGap of clearance, never dropping below the readability
// floor.
func radiusFromClosestPair(positions models.Positions, maxRadius, minGap float64) float64 {
	minDist := closestPairDistance(positions)
	radius := maxRadius
	if minDist < math.MaxFloat64 {
		radius = math.Min(radius, (minDist-minGap)/2.0)
	}
	return math.Max(radius, csZoneRadiusFloor)
}

func closestPairDistance(positions models.Positions) float64 {
	minDist := math.MaxFloat64
	for i := range positions {
		for j := i + 1; j < len(positions); j++ {
			minDist = math.Min(minDist, positions[i].Subtract(positions[j]).Distance())
		}
	}
	return minDist
}

// fitToCanvas centers the bounding box of the points on the canvas and
// uniformly scales them about the center so the box fits inside the padded
// draw area. With fill=true the figure is also scaled up to fill the area;
// otherwise it is only ever shrunk. Returns the applied scale factor.
func fitToCanvas(positions models.Positions, metrics canvasMetrics, pad float64, fill bool) float64 {
	minPos, maxPos := getMinMaxPositions(positions)
	boxCenter := minPos.Add(maxPos).DivideScalar(2)
	allowedDimension := metrics.side - 2.0*pad

	fitScale := math.MaxFloat64
	if maxPos.X-minPos.X > 1e-6 {
		fitScale = math.Min(fitScale, allowedDimension/(maxPos.X-minPos.X))
	}
	if maxPos.Y-minPos.Y > 1e-6 {
		fitScale = math.Min(fitScale, allowedDimension/(maxPos.Y-minPos.Y))
	}
	if fitScale == math.MaxFloat64 || fitScale <= 0 {
		fitScale = 1.0
	}
	if !fill {
		fitScale = math.Min(fitScale, 1.0)
	}
	for i := range positions {
		positions[i] = positions[i].Subtract(boxCenter).MultiplyScalar(fitScale).Add(metrics.center())
	}
	return fitScale
}

func getMinMaxPositions(positions models.Positions) (minPos, maxPos models.Position) {
	if len(positions) == 0 {
		return models.Position{}, models.Position{}
	}

	minPos, maxPos = positions[0], positions[0]
	for _, position := range positions[1:] {
		if position.X < minPos.X {
			minPos.X = position.X
		}
		if position.Y < minPos.Y {
			minPos.Y = position.Y
		}
		if position.X > maxPos.X {
			maxPos.X = position.X
		}
		if position.Y > maxPos.Y {
			maxPos.Y = position.Y
		}
	}
	return minPos, maxPos
}

// isStructuralIgnored reports whether a connection type is skipped when the
// connection graph drives geometry (adjacency, hub spokes).
func isStructuralIgnored(connectionType string) bool {
	connectionTypes := registry.GetConnectionTypeValues()
	return connectionType == connectionTypes.Proximity || connectionType == connectionTypes.Portal
}

func orderZonesByZeroAngle(zones []template_model.Zone, zeroAngleZone string) []template_model.Zone {
	if zeroAngleZone == "" {
		return zones
	}
	pivot := -1
	for i, z := range zones {
		if z.Name == zeroAngleZone {
			pivot = i
			break
		}
	}
	if pivot <= 0 {
		return zones
	}
	out := make([]template_model.Zone, 0, len(zones))
	out = append(out, zones[pivot:]...)
	out = append(out, zones[:pivot]...)
	return out
}

func allHavePosition(zones []template_model.Zone) bool {
	for _, z := range zones {
		if z.GeneratorPosition == nil {
			return false
		}
	}

	return true
}

func allHaveManualPosition(zones []template_model.Zone) bool {
	for _, z := range zones {
		if z.ManualPosition == nil {
			return false
		}
	}

	return len(zones) > 0
}

func allHaveRing(zones []template_model.Zone) bool {
	for _, z := range zones {
		if z.GeneratorRing == nil {
			return false
		}
	}

	return true
}

func connectedComponents(n int, adj [][]int) [][]int {
	id := make([]int, n)
	for i := range id {
		id[i] = -1
	}
	var comps [][]int
	for start := range n {
		if id[start] >= 0 {
			continue
		}

		var comp []int
		queue := []int{start}
		id[start] = len(comps)
		for len(queue) > 0 {
			u := queue[0]
			queue = queue[1:]
			comp = append(comp, u)
			for _, v := range adj[u] {
				if id[v] < 0 {
					id[v] = len(comps)
					queue = append(queue, v)
				}
			}
		}
		comps = append(comps, comp)
	}
	return comps
}

func positionCentroid(zones []template_model.Zone) data.Vec2[float64] {
	if len(zones) == 0 {
		return data.NewVec2(0.5, 0.5)
	}

	sum := data.NewVec2(0.0, 0.0)
	for _, z := range zones {
		p := *z.GeneratorPosition // Is this required to be copied? can't it be used directly safely?
		sum = sum.Add(data.NewVec2(p[0], p[1]))
	}
	return sum.DivideScalar(float64(len(zones)))
}

func positionAngle(z template_model.Zone, rawCenter data.Vec2[float64]) float64 {
	p := *z.GeneratorPosition
	return math.Atan2(p[1]-rawCenter.Y, p[0]-rawCenter.X)
}

// sortIndicesByAngle returns the given zone indices reordered by their raw
// generator position's angle around the raw centroid, preserving neighbor
// ordering when zones are re-projected onto a canvas ring.
func sortIndicesByAngle(zones []template_model.Zone, indices []int, rawCenter data.Vec2[float64]) []int {
	sorted := append([]int(nil), indices...)
	sort.SliceStable(sorted, func(i, j int) bool {
		return positionAngle(zones[sorted[i]], rawCenter) < positionAngle(zones[sorted[j]], rawCenter)
	})

	return sorted
}
