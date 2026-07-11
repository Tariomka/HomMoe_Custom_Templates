package services

import (
	"image"
	"math"
	"slices"
	"sort"
	"strings"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/preview"
)

// PreviewLayout is the full geometry of a preview rendered into a square
// canvas of the requested side length.
type PreviewLayout struct {
	Positions   map[string]image.Point
	Zones       []preview.PreviewZone
	Connections []preview.PreviewConnection
	ZoneRadius  int
}

const (
	csCanvasSide     = 700.0
	csMargin         = 18.0
	csZoneRadiusMax  = 38.0
	csHubRadiusMin   = 28.0
	csMinGap         = 6.0
	csConnectionGap  = 26.0 // ring layout - visible chord clearance between zones
	csClusterGap     = 20.0
	scatterIdealMult = 3.2
	scatterMinDist   = 3.8
	scatterEdgeClear = 1.2
)

func canvasScale(side float64) float64 { return side / csCanvasSide }

// BuildPreviewLayout computes zone positions, radius and connections for a
// preview canvas of the given side length. The layout strategy is picked to
// match the in-game generator: Circles uses concentric rings keyed off the
// GeneratorRing stamps; Square, Geometric and Cross are placed verbatim from
// their GeneratorPosition stamps (centred and scaled to fit) so the exact
// geometric figure is preserved; Random scatters zones using the
// GeneratorPosition stamps with hard-floor and edge-clearance correction
// passes; all other topologies fall back to the classic ring / hub-and-spoke
// renderer.
func BuildPreviewLayout(template *entities.RmgTemplate, topology config.MapTopology, side float64) PreviewLayout {
	layout := PreviewLayout{Positions: map[string]image.Point{}}
	if template == nil || len(template.Variants) == 0 {
		return layout
	}
	variant := template.Variants[0]
	if len(variant.Zones) == 0 {
		return layout
	}

	// Apply the optional ZeroAngleZone rotation so the first ring slot lines
	// up with the template author's chosen anchor, then lay out every zone with
	// the topology-specific renderer. Tournament templates are not special-
	// cased here: both player clusters are laid out together at full canvas
	// size (the generator seeds the two halves with mirrored positions and, for
	// hub topologies, layoutMultiHub fans the clusters out), so the preview and
	// the zone editor share one consistent, fully reversible coordinate system.
	zones := orderZonesByZeroAngle(variant.Zones, variant.Orientation.ZeroAngleZone)
	connections := variant.Connections

	dispatchClusterLayout(&layout, zones, connections, topology, side)

	for _, zone := range variant.Zones {
		pos, ok := layout.Positions[zone.Name]
		if !ok {
			continue
		}
		// A zone is only drawn as a hub when the template actually contains a
		// hub zone (named "Hub" or "Hub-*"). Connectivity-based guesses are not
		// used here: in topologies like Random or Circles an ordinary neutral
		// can happen to touch every spawn without being a hub, which previously
		// made the hub marker appear (and flicker) on non-hub zones.
		isHub := strings.EqualFold(zone.Name, "Hub") || strings.HasPrefix(zone.Name, "Hub-")
		preview := preview.PreviewZone{
			Name:     zone.Name,
			Letter:   ExtractZoneLetter(zone.Name),
			Center:   pos,
			Tier:     ClassifyZoneTier(zone),
			IsHub:    isHub,
			IsPlayer: strings.HasPrefix(zone.Name, "Spawn-"),
		}
		for _, mainObject := range zone.MainObjects {
			switch {
			case strings.EqualFold(mainObject.Type, "Spawn"):
				preview.HasCastle = true
				preview.Castles++
				if strings.HasPrefix(mainObject.Spawn, "Player") {
					for _, ch := range mainObject.Spawn[len("Player"):] {
						if ch >= '0' && ch <= '9' {
							preview.Owner = preview.Owner*10 + int(ch-'0')
						}
					}
				}
			case strings.EqualFold(mainObject.Type, "City"):
				preview.HasCastle = true
				preview.Castles++
			}
		}
		layout.Zones = append(layout.Zones, preview)
	}

	// Connections - only render those whose endpoints survived the strip.
	// Parallel edges between the same unordered pair are fanned out into
	// distinct curves (matching the manual zone editor); a lone edge keeps its
	// control point on the midpoint and therefore renders straight.
	layout.Connections = buildPreviewConnections(variant.Connections, layout.Positions)
	return layout
}

// previewParallelGap is the perpendicular spacing between parallel preview
// edges, matched to the manual zone editor's bulge spacing so both views fan
// multiple connections out the same way.
const previewParallelGap = 22.0

// buildPreviewConnections turns the variant's connections into drawable preview
// edges. Connections sharing the same unordered endpoint pair are grouped and
// each is given a perpendicular bulge so they do not collapse onto a single
// overlapping line.
func buildPreviewConnections(
	connections []entities.Connection,
	positions map[string]image.Point) []preview.PreviewConnection {
	type pairKey struct{ a, b string }
	groups := make(map[pairKey][]entities.Connection)
	order := make([]pairKey, 0)
	for _, conn := range connections {
		if _, ok := positions[conn.From]; !ok {
			continue
		}
		if _, ok := positions[conn.To]; !ok {
			continue
		}
		a, b := conn.From, conn.To
		if a > b {
			a, b = b, a
		}
		key := pairKey{a, b}
		if _, seen := groups[key]; !seen {
			order = append(order, key)
		}
		groups[key] = append(groups[key], conn)
	}

	result := make([]preview.PreviewConnection, 0, len(connections))
	for _, key := range order {
		group := groups[key]
		count := len(group)
		for index, conn := range group {
			a := positions[conn.From]
			b := positions[conn.To]
			// Bulge off a canonical baseline (sorted endpoints) so every
			// parallel edge fans out from the same side regardless of the
			// direction in which it happens to be stored.
			canonicalA, canonicalB := a, b
			if conn.From > conn.To {
				canonicalA, canonicalB = canonicalB, canonicalA
			}
			dx := float64(canonicalB.X - canonicalA.X)
			dy := float64(canonicalB.Y - canonicalA.Y)
			distance := math.Hypot(dx, dy)
			if distance < 1 {
				distance = 1
			}
			normalX := dy / distance
			normalY := -dx / distance
			spread := (float64(index) - float64(count-1)/2.0) * previewParallelGap
			midX := float64(a.X+b.X) / 2.0
			midY := float64(a.Y+b.Y) / 2.0
			ctrl := image.Pt(
				int(math.Round(midX+2.0*spread*normalX)),
				int(math.Round(midY+2.0*spread*normalY)),
			)
			isPortal := len(conn.PortalPlacementRulesFrom) > 0 ||
				len(conn.PortalPlacementRulesTo) > 0 ||
				conn.ConnectionType == "Portal"
			result = append(result, preview.PreviewConnection{Start: a, End: b, Ctrl: ctrl, Portal: isPortal})
		}
	}
	return result
}

// dispatchClusterLayout writes positions for the given zones into the layout,
// picking the topology-specific renderer. Each path sets layout.Positions and
// layout.ZoneRadius.
func dispatchClusterLayout(
	layout *PreviewLayout,
	zones []entities.Zone,
	connections []entities.Connection,
	topology config.MapTopology,
	side float64) {
	switch {
	case allHaveManualPosition(zones):
		layoutManualPositions(layout, zones, side)
	case (topology == config.TopologyCircles) && allHaveRing(zones):
		layoutBalancedRings(layout, zones, side)
	case isFixedGeometryTopology(topology) && allHavePosition(zones):
		layoutFixedPositions(layout, zones, side)
	case isScatterTopology(topology) && allHavePosition(zones):
		layoutScatter(layout, zones, connections, side)
	default:
		layoutRingOrHub(layout, zones, connections, side)
	}
}

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
// placed verbatim (only centred and scaled to fit) so the preview reproduces
// the intended shape instead of relaxing it into a scatter.
func isFixedGeometryTopology(topology config.MapTopology) bool {
	switch topology {
	case config.TopologySquare, config.TopologyGeometric, config.TopologyCross, config.TopologyFractal:
		return true
	default:
		return false
	}
}

func allHaveManualPosition(zones []entities.Zone) bool {
	for _, z := range zones {
		if z.ManualPosition == nil {
			return false
		}
	}
	return len(zones) > 0
}

// layoutManualPositions places zones exactly where the manual zone editor put
// them: canvas = normalized position × side. The mapping must stay trivially
// invertible (p = pos / side) so dragging in the editor is exact. The zone
// radius shrinks just enough to keep the closest pair of zones from
// overlapping.
func layoutManualPositions(layout *PreviewLayout, zones []entities.Zone, side float64) {
	scale := canvasScale(side)
	radius := csZoneRadiusMax * scale
	minGap := csMinGap * scale

	minDist := math.MaxFloat64
	for i := range zones {
		for j := i + 1; j < len(zones); j++ {
			pi := *zones[i].ManualPosition
			pj := *zones[j].ManualPosition
			dist := math.Hypot((pi[0]-pj[0])*side, (pi[1]-pj[1])*side)
			minDist = math.Min(minDist, dist)
		}
	}
	if minDist < math.MaxFloat64 {
		radius = math.Min(radius, (minDist-minGap)/2.0)
	}
	radius = math.Max(radius, 8.0)

	layout.ZoneRadius = int(math.Round(radius))
	for _, zone := range zones {
		p := *zone.ManualPosition
		layout.Positions[zone.Name] = image.Pt(
			int(math.Round(p[0]*side)),
			int(math.Round(p[1]*side)))
	}
}

// layoutFixedPositions places zones at their exact GeneratorPosition stamps,
// preserving the deterministic geometric figure built by the Square, Geometric
// and Cross topologies. The normalized positions are centred and uniformly
// scaled to fill the padded canvas (never relaxed), then the zone radius is
// shrunk just enough to keep the closest pair from overlapping.
func layoutFixedPositions(layout *PreviewLayout, zones []entities.Zone, side float64) {
	n := len(zones)
	if n == 0 {
		layout.ZoneRadius = scaledInt(csZoneRadiusMax, side)
		return
	}
	scale := canvasScale(side)
	margin := csMargin * scale
	minGap := csMinGap * scale
	zoneRadiusMax := csZoneRadiusMax * scale
	cx := side / 2.0
	cy := side / 2.0

	if n == 1 {
		layout.ZoneRadius = int(math.Round(zoneRadiusMax))
		layout.Positions[zones[0].Name] = image.Pt(int(cx), int(cy))
		return
	}

	px := make([]float64, n)
	py := make([]float64, n)
	for i, z := range zones {
		p := *z.GeneratorPosition
		px[i] = p[0]
		py[i] = p[1]
	}

	// Centre the bounding box and scale uniformly so the figure fills the
	// canvas inside a margin that reserves room for the largest possible zone
	// radius (so the eventual radius, which is never larger, always fits).
	minX, maxX := minMax(px)
	minY, maxY := minMax(py)
	spanX := maxX - minX
	spanY := maxY - minY
	pad := margin + zoneRadiusMax
	drawW := side - 2.0*pad
	drawH := side - 2.0*pad
	fitScale := math.MaxFloat64
	if spanX > 1e-6 {
		fitScale = math.Min(fitScale, drawW/spanX)
	}
	if spanY > 1e-6 {
		fitScale = math.Min(fitScale, drawH/spanY)
	}
	if fitScale == math.MaxFloat64 || fitScale <= 0 {
		fitScale = 1.0
	}
	boxCx := (minX + maxX) / 2.0
	boxCy := (minY + maxY) / 2.0
	for i := range px {
		px[i] = cx + (px[i]-boxCx)*fitScale
		py[i] = cy + (py[i]-boxCy)*fitScale
	}

	// Radius from the closest pair so neighbouring zones never overlap.
	minDist := math.MaxFloat64
	for i := range n {
		for j := i + 1; j < n; j++ {
			minDist = math.Min(minDist, math.Hypot(px[i]-px[j], py[i]-py[j]))
		}
	}
	radius := zoneRadiusMax
	if minDist < math.MaxFloat64 {
		radius = math.Min(radius, (minDist-minGap)/2.0)
	}
	radius = math.Max(radius, 8.0)

	layout.ZoneRadius = int(math.Round(radius))
	for i, z := range zones {
		layout.Positions[z.Name] = image.Pt(int(math.Round(px[i])), int(math.Round(py[i])))
	}
}

func allHaveRing(zones []entities.Zone) bool {
	for _, z := range zones {
		if z.GeneratorRing == nil {
			return false
		}
	}
	return true
}

func layoutBalancedRings(layout *PreviewLayout, zones []entities.Zone, side float64) {
	zoneCount := len(zones)
	if zoneCount == 0 {
		layout.ZoneRadius = scaledInt(csZoneRadiusMax, side)
		return
	}
	scale := canvasScale(side)
	margin := csMargin * scale
	minGap := csMinGap * scale
	zoneRadiusMax := csZoneRadiusMax * scale
	cx := side / 2.0
	cy := side / 2.0

	if zoneCount == 1 {
		layout.ZoneRadius = int(math.Round(zoneRadiusMax))
		layout.Positions[zones[0].Name] = image.Pt(int(cx), int(cy))
		return
	}

	// Collect & order the present GeneratorRing values, then build the
	// outer→inner ring-index map (largest tier → outermost ring index 0).
	presentSet := map[int]bool{}
	for _, z := range zones {
		presentSet[*z.GeneratorRing] = true
	}
	present := make([]int, 0, len(presentSet))
	for t := range presentSet {
		present = append(present, t)
	}
	sort.Ints(present)
	ringCount := len(present)
	tierToRing := make(map[int]int, ringCount)
	for ri, tier := range present {
		tierToRing[tier] = ringCount - 1 - ri
	}

	if ringCount < 2 {
		// All zones in a single ring - degenerate; fall back to the ring path.
		layoutRingOrHub(layout, zones, nil, side)
		return
	}

	ringIndices := make([][]int, ringCount)
	ringLabel := make([]int, zoneCount)
	for i, z := range zones {
		r := tierToRing[*z.GeneratorRing]
		ringLabel[i] = r
		ringIndices[r] = append(ringIndices[r], i)
	}

	drawRadius := side/2.0 - margin - zoneRadiusMax

	assignRingRadii := func(zr float64) []float64 {
		mc := 2.0*zr + minGap
		radii := make([]float64, ringCount)
		for ringIndex := range ringCount {
			cnt := len(ringIndices[ringIndex])
			natural := drawRadius * float64(ringIndex+1) / float64(ringCount)
			withinRing := 0.0
			if cnt >= 2 {
				withinRing = mc / (2.0 * math.Sin(math.Pi/float64(cnt)))
			} else if cnt == 1 && ringIndex > 0 {
				withinRing = mc
			}
			afterPrev := 0.0
			if ringIndex > 0 {
				afterPrev = radii[ringIndex-1] + mc
			}
			radii[ringIndex] = math.Max(natural, math.Max(withinRing, afterPrev))
		}
		return radii
	}

	// Binary-search the largest zone radius that keeps the outer ring inside
	// the available draw radius.
	lo, hi := 8.0, zoneRadiusMax
	for range 32 {
		mid := (lo + hi) / 2.0
		r2 := assignRingRadii(mid)
		if r2[ringCount-1] <= drawRadius {
			lo = mid
		} else {
			hi = mid
		}
	}
	zoneRadius := math.Max(lo, 8.0)
	ringRadii := assignRingRadii(zoneRadius)
	layout.ZoneRadius = int(math.Round(zoneRadius))

	rawCx, rawCy := positionCentroid(zones)

	for ringIndex := range ringCount {
		group := ringIndices[ringIndex]
		cnt := len(group)
		if cnt == 0 {
			continue
		}
		if cnt == 1 && ringIndex == 0 {
			layout.Positions[zones[group[0]].Name] = image.Pt(int(cx), int(cy))
			continue
		}
		// Sort zones in this ring by their raw position's angle around the
		// cluster centroid so the rendered ring preserves neighbour ordering.
		sorted := append([]int(nil), group...)
		sort.SliceStable(sorted, func(i, j int) bool {
			return positionAngle(zones[sorted[i]], rawCx, rawCy) <
				positionAngle(zones[sorted[j]], rawCx, rawCy)
		})
		firstAngle := positionAngle(zones[sorted[0]], rawCx, rawCy)
		canvasRadius := ringRadii[ringIndex]
		for j, idx := range sorted {
			angle := firstAngle + 2.0*math.Pi*float64(j)/float64(cnt)
			x := cx + math.Cos(angle)*canvasRadius
			y := cy + math.Sin(angle)*canvasRadius
			layout.Positions[zones[idx].Name] = image.Pt(int(math.Round(x)), int(math.Round(y)))
		}
	}
}

func layoutScatter(layout *PreviewLayout, zones []entities.Zone, conns []entities.Connection, side float64) {
	n := len(zones)
	if n == 0 {
		layout.ZoneRadius = scaledInt(csZoneRadiusMax, side)
		return
	}
	scale := canvasScale(side)
	margin := csMargin * scale
	minGap := csMinGap * scale
	zoneRadiusMax := csZoneRadiusMax * scale
	cx := side / 2.0
	cy := side / 2.0

	if n == 1 {
		layout.ZoneRadius = int(math.Round(zoneRadiusMax))
		layout.Positions[zones[0].Name] = image.Pt(int(cx), int(cy))
		return
	}

	idx := make(map[string]int, n)
	for i, z := range zones {
		idx[z.Name] = i
	}

	// Direct-only adjacency drives both the radius heuristic and Pass B.
	adj := make([][]int, n)
	addAdj := func(a, b int) {
		if slices.Contains(adj[a], b) {
			return
		}
		adj[a] = append(adj[a], b)
		adj[b] = append(adj[b], a)
	}
	for _, c := range conns {
		if isStructuralIgnored(c.ConnectionType) {
			continue
		}
		ai, ok1 := idx[c.From]
		bi, ok2 := idx[c.To]
		if !ok1 || !ok2 || ai == bi {
			continue
		}
		addAdj(ai, bi)
	}

	// Component sizing: zoneRadius is derived from the largest component so
	// readability stays consistent regardless of how many neutrals exist.
	comps := connectedComponents(n, adj)
	maxCompSize := 1
	for _, c := range comps {
		if len(c) > maxCompSize {
			maxCompSize = len(c)
		}
	}
	ringRadius0 := side/2.0 - margin
	zoneRadius := zoneRadiusMax
	if maxCompSize >= 2 {
		chord0 := 2.0 * ringRadius0 * math.Sin(math.Pi/float64(maxCompSize))
		zoneRadius = math.Min(zoneRadiusMax, (chord0-minGap)/2.0)
		zoneRadius = math.Max(zoneRadius, 8.0)
	}
	idealEdge := zoneRadius * scatterIdealMult

	// Scale raw [0,1] generator positions so the mean direct-edge length
	// matches the ideal. Empty graphs fall back to spanning the draw area.
	rawEdgeSum, rawEdgeCount := 0.0, 0
	for i := range n {
		for _, j := range adj[i] {
			if j <= i {
				continue
			}
			pi := *zones[i].GeneratorPosition
			pj := *zones[j].GeneratorPosition
			rawEdgeSum += math.Hypot(pi[0]-pj[0], pi[1]-pj[1])
			rawEdgeCount++
		}
	}
	gScale := math.Min(side-2.0*margin, side-2.0*margin)
	if rawEdgeCount > 0 {
		mean := rawEdgeSum / float64(rawEdgeCount)
		if mean > 1e-6 {
			gScale = idealEdge / mean
		}
	}

	// Translate to centroid frame, scale, then translate to canvas centre.
	rawCx, rawCy := positionCentroid(zones)
	px := make([]float64, n)
	py := make([]float64, n)
	for i, z := range zones {
		p := *z.GeneratorPosition
		px[i] = (p[0] - rawCx) * gScale
		py[i] = (p[1] - rawCy) * gScale
	}
	pad := zoneRadius + margin
	drawW := side - 2.0*pad
	drawH := side - 2.0*pad
	rawMinX, rawMaxX := minMax(px)
	rawMinY, rawMaxY := minMax(py)
	fitScale := 1.0
	if rawMaxX-rawMinX > drawW && rawMaxX-rawMinX > 1e-3 {
		fitScale = math.Min(fitScale, drawW/(rawMaxX-rawMinX))
	}
	if rawMaxY-rawMinY > drawH && rawMaxY-rawMinY > 1e-3 {
		fitScale = math.Min(fitScale, drawH/(rawMaxY-rawMinY))
	}
	for i := range n {
		px[i] = cx + px[i]*fitScale
		py[i] = cy + py[i]*fitScale
	}

	// ── Pass A + Pass B: tight intra-zone clearance, off-line nudging ────
	relaxPasses(px, py, adj, zoneRadius)

	// Final fit: recentre the bounding box then shrink only if it overflows
	// the padded canvas
	finalMinX, finalMaxX := minMax(px)
	finalMinY, finalMaxY := minMax(py)
	finalCx := (finalMinX + finalMaxX) / 2.0
	finalCy := (finalMinY + finalMaxY) / 2.0
	for i := range n {
		px[i] += cx - finalCx
		py[i] += cy - finalCy
	}
	finalMinX, finalMaxX = minMax(px)
	finalMinY, finalMaxY = minMax(py)
	spanX := finalMaxX - finalMinX
	spanY := finalMaxY - finalMinY
	allowW := side - 2.0*pad
	allowH := side - 2.0*pad
	shrink := 1.0
	if spanX > allowW && spanX > 1e-3 {
		shrink = math.Min(shrink, allowW/spanX)
	}
	if spanY > allowH && spanY > 1e-3 {
		shrink = math.Min(shrink, allowH/spanY)
	}
	if shrink < 1.0 {
		for i := range n {
			px[i] = cx + (px[i]-cx)*shrink
			py[i] = cy + (py[i]-cy)*shrink
		}
		zoneRadius = math.Max(zoneRadius*shrink, 8.0)
	}

	layout.ZoneRadius = int(math.Round(zoneRadius))
	for i, z := range zones {
		layout.Positions[z.Name] = image.Pt(int(math.Round(px[i])), int(math.Round(py[i])))
	}
}

func relaxPasses(px, py []float64, adj [][]int, zoneRadius float64) {
	n := len(px)
	minDist := zoneRadius * scatterMinDist
	edgeClear := zoneRadius * scatterEdgeClear

	for range 500 {
		moved := false

		// A: hard floor.
		for i := range n {
			for j := i + 1; j < n; j++ {
				dx := px[i] - px[j]
				dy := py[i] - py[j]
				d := math.Hypot(dx, dy)
				if d >= minDist {
					continue
				}
				if d < 1e-3 {
					dx, dy, d = 1, 0, 1e-3
				}
				push := (minDist - d) / 2.0
				px[i] += dx / d * push
				py[i] += dy / d * push
				px[j] -= dx / d * push
				py[j] -= dy / d * push
				moved = true
			}
		}

		// B: edge clearance.
		for a := range n {
			for _, b := range adj[a] {
				if b <= a {
					continue
				}
				ex := px[b] - px[a]
				ey := py[b] - py[a]
				elen2 := ex*ex + ey*ey
				if elen2 < 1e-3 {
					continue
				}
				elenInv := 1.0 / math.Sqrt(elen2)
				for c := range n {
					if c == a || c == b {
						continue
					}
					tProj := ((px[c]-px[a])*ex + (py[c]-py[a])*ey) / elen2
					if tProj < 0 || tProj > 1 {
						continue
					}
					projX := px[a] + tProj*ex
					projY := py[a] + tProj*ey
					nx := px[c] - projX
					ny := py[c] - projY
					dist := math.Hypot(nx, ny)
					if dist >= edgeClear {
						continue
					}
					var perpX, perpY float64
					if dist < 1e-3 {
						perpX = ey * elenInv
						perpY = -ex * elenInv
					} else {
						perpX = nx / dist
						perpY = ny / dist
					}
					ax := projX + perpX*edgeClear
					ay := projY + perpY*edgeClear
					bx := projX - perpX*edgeClear
					by := projY - perpY*edgeClear
					scoreA, scoreB := 0.0, 0.0
					for _, nb := range adj[c] {
						dax := ax - px[nb]
						day := ay - py[nb]
						dbx := bx - px[nb]
						dby := by - py[nb]
						scoreA += dax*dax + day*day
						scoreB += dbx*dbx + dby*dby
					}
					if scoreB < scoreA {
						px[c], py[c] = bx, by
					} else {
						px[c], py[c] = ax, ay
					}
					moved = true
				}
			}
		}

		if !moved {
			break
		}
	}
}

// Used for structured topologies (Default, HubAndSpoke, Chain, SharedWeb).
// Multi-hub "Hub-*" templates fan their spokes out from each cluster centre;
// otherwise zones land on a single outer ring with an optional centre hub.
func layoutRingOrHub(layout *PreviewLayout, zones []entities.Zone, conns []entities.Connection, side float64) {
	n := len(zones)
	scale := canvasScale(side)
	margin := csMargin * scale
	minGap := csMinGap * scale
	zoneRadiusMax := csZoneRadiusMax * scale
	hubRadiusMin := csHubRadiusMin * scale
	connectionGap := csConnectionGap * scale
	cx := side / 2.0
	cy := side / 2.0

	if n == 0 {
		layout.ZoneRadius = int(math.Round(zoneRadiusMax))
		return
	}

	// Multi-hub tournament layout: clusters fan out around the canvas.
	var hubIndices []int
	for i, z := range zones {
		if strings.HasPrefix(z.Name, "Hub-") {
			hubIndices = append(hubIndices, i)
		}
	}
	if len(hubIndices) >= 2 {
		layoutMultiHub(layout, zones, conns, hubIndices, side)
		return
	}

	// Hub detection: only an explicitly named "Hub" zone is treated as a hub.
	// The preview is a faithful representation of the template data, so
	// connectivity is never used to guess an implicit hub.
	hubIdx := -1
	for i, z := range zones {
		if z.Name == "Hub" {
			hubIdx = i
			break
		}
	}

	var outer []int
	for i := range zones {
		if i != hubIdx {
			outer = append(outer, i)
		}
	}
	outerN := max(len(outer), 1)
	ringRadius0 := side/2.0 - margin
	sinA := 1.0
	if outerN > 1 {
		sinA = math.Sin(math.Pi / float64(outerN))
	}
	var zoneRadius float64
	if hubIdx < 0 {
		zoneRadius = (2.0*ringRadius0*sinA - connectionGap) / (2.0 * (1.0 + sinA))
	} else {
		chord0 := 2.0 * ringRadius0 * math.Sin(math.Pi/math.Max(1, float64(outerN)))
		zoneRadius = (chord0 - connectionGap) / 2.0
	}
	zoneRadius = math.Min(zoneRadiusMax, math.Max(zoneRadius, 4.0))
	layout.ZoneRadius = int(math.Round(zoneRadius))

	ringRadius := math.Max(hubRadiusMin+zoneRadius+minGap,
		math.Min(ringRadius0, side/2.0-zoneRadius-margin))

	if hubIdx >= 0 {
		layout.Positions[zones[hubIdx].Name] = image.Pt(int(cx), int(cy))
	}
	if n == 1 {
		layout.Positions[zones[0].Name] = image.Pt(int(cx), int(cy))
		return
	}
	for i, idx := range outer {
		angle := -math.Pi/2.0 + float64(i)*2.0*math.Pi/float64(outerN)
		x := cx + math.Cos(angle)*ringRadius
		y := cy + math.Sin(angle)*ringRadius
		layout.Positions[zones[idx].Name] = image.Pt(int(math.Round(x)), int(math.Round(y)))
	}
}

func layoutMultiHub(
	layout *PreviewLayout,
	zones []entities.Zone,
	conns []entities.Connection,
	hubIndices []int,
	side float64,
) {
	scale := canvasScale(side)
	margin := csMargin * scale
	minGap := csMinGap * scale
	zoneRadiusMax := csZoneRadiusMax * scale
	hubRadiusMin := csHubRadiusMin * scale
	cx := side / 2.0
	cy := side / 2.0

	zoneIdx := make(map[string]int, len(zones))
	for i, z := range zones {
		zoneIdx[z.Name] = i
	}
	// Build per-hub spoke list (Direct connections only, dedup'd).
	hubSpokes := make(map[string][]int, len(hubIndices))
	for _, h := range hubIndices {
		hub := zones[h].Name
		seen := map[int]bool{}
		for _, c := range conns {
			if isStructuralIgnored(c.ConnectionType) {
				continue
			}
			other := ""
			switch {
			case c.From == hub:
				other = c.To
			case c.To == hub:
				other = c.From
			}
			if other == "" {
				continue
			}
			oi, ok := zoneIdx[other]
			if !ok || seen[oi] {
				continue
			}
			seen[oi] = true
			hubSpokes[hub] = append(hubSpokes[hub], oi)
		}
	}
	numHubs := len(hubIndices)
	maxSpokes := 1
	for _, s := range hubSpokes {
		if len(s) > maxSpokes {
			maxSpokes = len(s)
		}
	}

	canvasHalf := side/2.0 - margin
	sinB := 0.0
	if numHubs > 1 {
		sinB = math.Sin(math.Pi / float64(numHubs))
	}
	sinA := 1.0
	if maxSpokes > 1 {
		sinA = math.Sin(math.Pi / float64(maxSpokes))
	}
	hubRing := 0.0
	if numHubs > 1 {
		hubRing = (canvasHalf + minGap/2.0) / (1.0 + sinB)
	}
	radialLeft := canvasHalf - hubRing
	minSpokeR := hubRadiusMin + minGap
	zoneRadius := math.Min(zoneRadiusMax, (radialLeft*sinA-minGap/2.0)/(1.0+sinA))
	zoneRadius = math.Max(1.0, zoneRadius)
	spokeRing := math.Max(radialLeft-zoneRadius, minSpokeR+zoneRadius)
	layout.ZoneRadius = int(math.Round(zoneRadius))

	for h, hubIndex := range hubIndices {
		hubAngle := -math.Pi/2.0 + float64(h)*2.0*math.Pi/float64(numHubs)
		hx, hy := cx, cy
		if numHubs > 1 {
			hx = cx + math.Cos(hubAngle)*hubRing
			hy = cy + math.Sin(hubAngle)*hubRing
		}
		layout.Positions[zones[hubIndex].Name] = image.Pt(int(math.Round(hx)), int(math.Round(hy)))
		spokes := hubSpokes[zones[hubIndex].Name]
		if len(spokes) == 0 {
			continue
		}
		spokeBase := hubAngle
		if numHubs == 1 {
			spokeBase = -math.Pi / 2.0
		}
		for i, si := range spokes {
			angle := spokeBase + float64(i)*2.0*math.Pi/float64(len(spokes))
			x := hx + math.Cos(angle)*spokeRing
			y := hy + math.Sin(angle)*spokeRing
			layout.Positions[zones[si].Name] = image.Pt(int(math.Round(x)), int(math.Round(y)))
		}
	}
	// Stragglers (e.g. cross-cluster zones) collapse to canvas centre.
	for _, z := range zones {
		if _, ok := layout.Positions[z.Name]; !ok {
			layout.Positions[z.Name] = image.Pt(int(cx), int(cy))
		}
	}
}

// ── Tier / letter helpers (kept for compatibility) ────────────────────────

// ExtractZoneLetter returns the trailing letter portion of a zone name like
// "Spawn-A" → "A" or "Neutral-C" → "C". Plain names (e.g. "Hub") pass through.
func ExtractZoneLetter(zoneName string) string {
	if after, ok := strings.CutPrefix(zoneName, "Spawn-"); ok {
		return after
	}
	if after, ok := strings.CutPrefix(zoneName, "Neutral-"); ok {
		return after
	}
	return zoneName
}

// ClassifyZoneTier guesses a neutral zone's tier from its content pools,
// layout name, and zone name (in order of reliability).
func ClassifyZoneTier(zone entities.Zone) int {
	if strings.HasPrefix(zone.Name, "Spawn-") {
		return 0
	}
	// Most reliable for templates generated by this tool: the guarded content
	// pool names embed the encounter tier (t2/t3/t4/t5). High-tier zones share
	// the "treasure" layout with medium-tier ones, so layout alone cannot
	// distinguish them.
	if tier := tierFromContentPools(zone.GuardedContentPool); tier > 0 {
		return tier
	}
	if tier := tierFromContentPools(zone.UnguardedContentPool); tier > 0 {
		return tier
	}
	layout := strings.ToLower(zone.Layout)
	switch {
	case strings.Contains(layout, "sides"):
		return 1
	case strings.Contains(layout, "treasure"):
		return 2
	case strings.Contains(layout, "center"):
		return 3
	}
	name := strings.ToLower(zone.Name)
	switch {
	case strings.Contains(name, "low") || strings.Contains(name, "side"):
		return 1
	case strings.Contains(name, "med") || strings.Contains(name, "treasure"):
		return 2
	case strings.Contains(name, "high") || strings.Contains(name, "center") || strings.Contains(name, "core"):
		return 3
	}
	return 1
}

// tierFromContentPools scans pool SIDs for tier markers ("_t2_".."_t5_") and
// returns the highest tier found mapped to preview tiers: t4/t5 → 3 (gold),
// t3 → 2 (silver), t2 → 1 (bronze). Returns 0 if no marker is present.
func tierFromContentPools(pools []string) int {
	best := 0
	for _, p := range pools {
		lp := strings.ToLower(p)
		var t int
		switch {
		case strings.Contains(lp, "_t5_") || strings.Contains(lp, "_t4_"):
			t = 3
		case strings.Contains(lp, "_t3_"):
			t = 2
		case strings.Contains(lp, "_t2_"):
			t = 1
		}
		if t > best {
			best = t
		}
	}
	return best
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

func scaledInt(value, side float64) int {
	return int(math.Round(value * canvasScale(side)))
}
