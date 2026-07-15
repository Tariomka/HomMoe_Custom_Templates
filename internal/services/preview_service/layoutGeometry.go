package preview_service

import (
	"image"
	"math"
	"sort"
	"strings"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
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
	scatterIdealMult  = 3.2
	scatterMinDist    = 3.8
	scatterEdgeClear  = 1.2
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
	cx            float64
	cy            float64
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

func (this canvasMetrics) center() image.Point {
	return image.Pt(int(this.cx), int(this.cy))
}

// placeTrivial handles the degenerate 0- and 1-zone cases shared by every
// renderer: the zone radius maxes out and a lone zone sits at the canvas
// center. Reports whether the layout was fully handled.
func (this *PreviewLayoutService) placeTrivial(zones []entities.Zone, metrics canvasMetrics) bool {
	if len(zones) > 1 {
		return false
	}
	this.layout.ZoneRadius = int(math.Round(metrics.zoneRadiusMax))
	if len(zones) == 1 {
		this.layout.Positions[zones[0].Name] = metrics.center()
	}
	return true
}

// commitPositions rounds the working coordinates into the layout and stores
// the final zone radius. Every renderer funnels through here.
func (this *PreviewLayoutService) commitPositions(
	zones []entities.Zone,
	px, py []float64,
	radius float64) {
	this.layout.ZoneRadius = int(math.Round(radius))
	for i, zone := range zones {
		this.layout.Positions[zone.Name] = image.Pt(int(math.Round(px[i])), int(math.Round(py[i])))
	}
}

// generatorCoords copies the zones' GeneratorPosition stamps into mutable
// coordinate slices. Callers must have verified allHavePosition.
func generatorCoords(zones []entities.Zone) (px, py []float64) {
	px = make([]float64, len(zones))
	py = make([]float64, len(zones))
	for i, zone := range zones {
		p := *zone.GeneratorPosition
		px[i] = p[0]
		py[i] = p[1]
	}
	return px, py
}

// radiusFromClosestPair shrinks maxRadius just enough that the closest pair
// of zones keeps minGap of clearance, never dropping below the readability
// floor.
func radiusFromClosestPair(px, py []float64, maxRadius, minGap float64) float64 {
	minDist := closestPairDistance(px, py)
	radius := maxRadius
	if minDist < math.MaxFloat64 {
		radius = math.Min(radius, (minDist-minGap)/2.0)
	}
	return math.Max(radius, csZoneRadiusFloor)
}

func closestPairDistance(px, py []float64) float64 {
	minDist := math.MaxFloat64
	for i := range px {
		for j := i + 1; j < len(px); j++ {
			minDist = math.Min(minDist, math.Hypot(px[i]-px[j], py[i]-py[j]))
		}
	}
	return minDist
}

// fitToCanvas centers the bounding box of the points on the canvas and
// uniformly scales them about the center so the box fits inside the padded
// draw area. With fill=true the figure is also scaled up to fill the area;
// otherwise it is only ever shrunk. Returns the applied scale factor.
func fitToCanvas(px, py []float64, metrics canvasMetrics, pad float64, fill bool) float64 {
	minX, maxX := minMax(px)
	minY, maxY := minMax(py)
	boxCx := (minX + maxX) / 2.0
	boxCy := (minY + maxY) / 2.0
	allowW := metrics.side - 2.0*pad
	allowH := metrics.side - 2.0*pad

	fitScale := math.MaxFloat64
	if maxX-minX > 1e-6 {
		fitScale = math.Min(fitScale, allowW/(maxX-minX))
	}
	if maxY-minY > 1e-6 {
		fitScale = math.Min(fitScale, allowH/(maxY-minY))
	}
	if fitScale == math.MaxFloat64 || fitScale <= 0 {
		fitScale = 1.0
	}
	if !fill {
		fitScale = math.Min(fitScale, 1.0)
	}
	for i := range px {
		px[i] = metrics.cx + (px[i]-boxCx)*fitScale
		py[i] = metrics.cy + (py[i]-boxCy)*fitScale
	}
	return fitScale
}

func minMax(values []float64) (float64, float64) {
	if len(values) == 0 {
		return 0, 0
	}
	lo, hi := values[0], values[0]
	for _, v := range values[1:] {
		if v < lo {
			lo = v
		}
		if v > hi {
			hi = v
		}
	}
	return lo, hi
}

// isStructuralIgnored reports whether a connection type is skipped when the
// connection graph drives geometry (adjacency, hub spokes).
func isStructuralIgnored(connectionType string) bool {
	return connectionType == "Proximity" || connectionType == "Portal"
}

func orderZonesByZeroAngle(zones []entities.Zone, zeroAngleZone string) []entities.Zone {
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
	out := make([]entities.Zone, 0, len(zones))
	out = append(out, zones[pivot:]...)
	out = append(out, zones[:pivot]...)
	return out
}

func allHavePosition(zones []entities.Zone) bool {
	for _, z := range zones {
		if z.GeneratorPosition == nil {
			return false
		}
	}
	return true
}

func allHaveManualPosition(zones []entities.Zone) bool {
	for _, z := range zones {
		if z.ManualPosition == nil {
			return false
		}
	}
	return len(zones) > 0
}

func allHaveRing(zones []entities.Zone) bool {
	for _, z := range zones {
		if z.GeneratorRing == nil {
			return false
		}
	}
	return true
}

// isScatterTopology reports whether the topology lays out zones from their
// GeneratorPosition stamps using the organic scatter renderer (mean-edge
// scaling plus the relaxation passes that nudge zones apart).
func isScatterTopology(topology config.MapTopology) bool {
	switch topology {
	case config.TopologyRandom, config.TopologyCircles:
		return true
	default:
		return false
	}
}

// isFixedGeometryTopology reports whether the topology defines an exact,
// deterministic geometric figure from its GeneratorPosition stamps. These are
// placed verbatim (only centered and scaled to fit) so the preview reproduces
// the intended shape instead of relaxing it into a scatter.
func isFixedGeometryTopology(topology config.MapTopology) bool {
	switch topology {
	case config.TopologySquare, config.TopologyGeometric, config.TopologyCross,
		config.TopologyFractal, config.TopologyGeometricHub:
		return true
	default:
		return false
	}
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

func positionCentroid(zones []entities.Zone) (float64, float64) {
	if len(zones) == 0 {
		return 0.5, 0.5
	}
	var sx, sy float64
	for _, z := range zones {
		p := *z.GeneratorPosition
		sx += p[0]
		sy += p[1]
	}
	return sx / float64(len(zones)), sy / float64(len(zones))
}

func positionAngle(z entities.Zone, rawCx, rawCy float64) float64 {
	p := *z.GeneratorPosition
	return math.Atan2(p[1]-rawCy, p[0]-rawCx)
}

// sortIndicesByAngle returns the given zone indices reordered by their raw
// generator position's angle around the raw centroid, preserving neighbour
// ordering when zones are re-projected onto a canvas ring.
func sortIndicesByAngle(zones []entities.Zone, indices []int, rawCx, rawCy float64) []int {
	sorted := append([]int(nil), indices...)
	sort.SliceStable(sorted, func(i, j int) bool {
		return positionAngle(zones[sorted[i]], rawCx, rawCy) <
			positionAngle(zones[sorted[j]], rawCx, rawCy)
	})
	return sorted
}

// hasHubName reports whether the zone name marks an explicit hub ("Hub" or
// "Hub-*"). Connectivity is never used to guess an implicit hub.
func hasHubName(name string) bool {
	return strings.EqualFold(name, "Hub") || strings.HasPrefix(name, "Hub-")
}
