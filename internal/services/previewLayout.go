package services

import (
	"image"
	"math"
	"sort"
	"strings"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/generator"
)

// PreviewZone is one zone laid out on the preview canvas.
type PreviewZone struct {
	Name      string
	Letter    string
	Center    image.Point
	IsPlayer  bool
	IsHub     bool
	Tier      int // 0 unknown, 1 bronze, 2 silver, 3 gold
	Owner     int
	HasCastle bool
	Castles   int
}

// PreviewConnection is a drawn link between two zones on the preview canvas.
type PreviewConnection struct {
	A, B   image.Point
	Portal bool
}

// PreviewLayout is the full geometry of a preview rendered into a square
// canvas of the requested side length.
type PreviewLayout struct {
	Positions   map[string]image.Point
	Zones       []PreviewZone
	Connections []PreviewConnection
	ZoneRadius  int
}

const (
	csCanvasSide     = 700.0
	csMargin         = 18.0
	csZoneRadiusMax  = 38.0
	csHubRadiusMin   = 28.0
	csMinGap         = 6.0
	csConnectionGap  = 26.0 // ring layout — visible chord clearance between zones
	csClusterGap     = 20.0
	scatterIdealMult = 3.2
	scatterMinDist   = 3.8
	scatterEdgeClear = 1.2
)

func canvasScale(side float64) float64 { return side / csCanvasSide }

// BuildPreviewLayout computes zone positions, radius and connections for a
// preview canvas of the given side length. The layout strategy is picked to
// match the in-game generator: Balanced uses concentric rings keyed off the
// GeneratorRing stamps; Random scatters zones using the GeneratorPosition
// stamps with hard-floor and edge-clearance correction passes; all other
// topologies fall back to the classic ring / hub-and-spoke renderer.
func BuildPreviewLayout(template *models.RmgTemplate, topology models.MapTopology, side float64) PreviewLayout {
	layout := PreviewLayout{Positions: map[string]image.Point{}}
	if template == nil || len(template.Variants) == 0 {
		return layout
	}
	variant := template.Variants[0]
	if len(variant.Zones) == 0 {
		return layout
	}

	// Tournament single-cluster strip: if Direct-only adjacency has exactly two
	// components, render only the first cluster at full canvas size so the
	// preview reads like a non-tournament layout
	zones, connections := stripFirstClusterIfTwo(variant.Zones, variant.Connections)

	// Apply the optional ZeroAngleZone rotation so the first ring slot lines
	// up with the template author's chosen anchor.
	zones = orderZonesByZeroAngle(zones, variant.Orientation.ZeroAngleZone)

	// Dispatch to the topology-specific layout. Each path writes positions
	// into `layout.Positions` and sets `layout.ZoneRadius`.
	switch {
	case (topology == generator.TopologyBalanced) && allHaveRing(zones):
		layoutBalancedRings(&layout, zones, side)
	case (topology == generator.TopologyRandom || topology == generator.TopologyBalanced) && allHavePosition(zones):
		layoutScatter(&layout, zones, connections, side)
	default:
		layoutRingOrHub(&layout, zones, connections, side)
	}

	implicitHub := findImplicitHubName(zones, connections)
	for _, zone := range variant.Zones {
		pos, ok := layout.Positions[zone.Name]
		if !ok {
			continue
		}
		isHub := strings.EqualFold(zone.Name, "Hub") || strings.HasPrefix(zone.Name, "Hub-")
		if implicitHub != "" && zone.Name == implicitHub {
			isHub = true
		}
		preview := PreviewZone{
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

	// Connections — only render those whose endpoints survived the strip.
	for _, conn := range variant.Connections {
		a, okA := layout.Positions[conn.From]
		b, okB := layout.Positions[conn.To]
		if !okA || !okB {
			continue
		}
		isPortal := len(conn.PortalPlacementRulesFrom) > 0 ||
			len(conn.PortalPlacementRulesTo) > 0 ||
			conn.ConnectionType == "Portal"
		layout.Connections = append(layout.Connections, PreviewConnection{A: a, B: b, Portal: isPortal})
	}
	return layout
}

func stripFirstClusterIfTwo(zones []models.RmgZone, conns []models.RmgConnection) ([]models.RmgZone, []models.RmgConnection) {
	n := len(zones)
	idx := make(map[string]int, n)
	for i, z := range zones {
		idx[z.Name] = i
	}
	adj := make([][]int, n)
	for _, c := range conns {
		if isStructuralIgnored(c.ConnectionType) {
			continue
		}
		ai, ok1 := idx[c.From]
		bi, ok2 := idx[c.To]
		if !ok1 || !ok2 {
			continue
		}
		adj[ai] = append(adj[ai], bi)
		adj[bi] = append(adj[bi], ai)
	}
	compID := make([]int, n)
	for i := range compID {
		compID[i] = -1
	}
	var comps [][]int
	for start := 0; start < n; start++ {
		if compID[start] >= 0 {
			continue
		}
		var comp []int
		queue := []int{start}
		compID[start] = len(comps)
		for len(queue) > 0 {
			u := queue[0]
			queue = queue[1:]
			comp = append(comp, u)
			for _, v := range adj[u] {
				if compID[v] < 0 {
					compID[v] = len(comps)
					queue = append(queue, v)
				}
			}
		}
		comps = append(comps, comp)
	}
	if len(comps) != 2 {
		return zones, conns
	}
	keep := make(map[int]bool, len(comps[0]))
	for _, i := range comps[0] {
		keep[i] = true
	}
	keptZones := make([]models.RmgZone, 0, len(comps[0]))
	keptNames := make(map[string]bool, len(comps[0]))
	for i, z := range zones {
		if keep[i] {
			keptZones = append(keptZones, z)
			keptNames[z.Name] = true
		}
	}
	keptConns := make([]models.RmgConnection, 0, len(conns))
	for _, c := range conns {
		if keptNames[c.From] && keptNames[c.To] {
			keptConns = append(keptConns, c)
		}
	}
	return keptZones, keptConns
}

func isStructuralIgnored(connectionType string) bool {
	return connectionType == "Proximity" || connectionType == "Portal"
}

func orderZonesByZeroAngle(zones []models.RmgZone, zeroAngleZone string) []models.RmgZone {
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
	out := make([]models.RmgZone, 0, len(zones))
	out = append(out, zones[pivot:]...)
	out = append(out, zones[:pivot]...)
	return out
}

func allHavePosition(zones []models.RmgZone) bool {
	for _, z := range zones {
		if z.GeneratorPosition == nil {
			return false
		}
	}
	return true
}

func allHaveRing(zones []models.RmgZone) bool {
	for _, z := range zones {
		if z.GeneratorRing == nil {
			return false
		}
	}
	return true
}

func layoutBalancedRings(layout *PreviewLayout, zones []models.RmgZone, side float64) {
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
		// All zones in a single ring — degenerate; fall back to the ring path.
		layoutRingOrHub(layout, zones, nil, side)
		return
	}

	ringIndices := make([][]int, ringCount)
	ringLabel := make([]int, n)
	for i, z := range zones {
		r := tierToRing[*z.GeneratorRing]
		ringLabel[i] = r
		ringIndices[r] = append(ringIndices[r], i)
	}

	drawRadius := side/2.0 - margin - zoneRadiusMax

	assignRingRadii := func(zr float64) []float64 {
		mc := 2.0*zr + minGap
		radii := make([]float64, ringCount)
		for r := 0; r < ringCount; r++ {
			cnt := len(ringIndices[r])
			natural := drawRadius * float64(r+1) / float64(ringCount)
			withinRing := 0.0
			if cnt >= 2 {
				withinRing = mc / (2.0 * math.Sin(math.Pi/float64(cnt)))
			} else if cnt == 1 && r > 0 {
				withinRing = mc
			}
			afterPrev := 0.0
			if r > 0 {
				afterPrev = radii[r-1] + mc
			}
			radii[r] = math.Max(natural, math.Max(withinRing, afterPrev))
		}
		return radii
	}

	// Binary-search the largest zone radius that keeps the outer ring inside
	// the available draw radius.
	lo, hi := 8.0, zoneRadiusMax
	for iter := 0; iter < 32; iter++ {
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

	for r := 0; r < ringCount; r++ {
		group := ringIndices[r]
		cnt := len(group)
		if cnt == 0 {
			continue
		}
		if cnt == 1 && r == 0 {
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
		canvasRadius := ringRadii[r]
		for j, idx := range sorted {
			angle := firstAngle + 2.0*math.Pi*float64(j)/float64(cnt)
			x := cx + math.Cos(angle)*canvasRadius
			y := cy + math.Sin(angle)*canvasRadius
			layout.Positions[zones[idx].Name] = image.Pt(int(math.Round(x)), int(math.Round(y)))
		}
	}
}

func layoutScatter(layout *PreviewLayout, zones []models.RmgZone, conns []models.RmgConnection, side float64) {
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
		for _, v := range adj[a] {
			if v == b {
				return
			}
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
	for i := 0; i < n; i++ {
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
	for i := 0; i < n; i++ {
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
	for i := 0; i < n; i++ {
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
		for i := 0; i < n; i++ {
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

	for pass := 0; pass < 500; pass++ {
		moved := false

		// A: hard floor.
		for i := 0; i < n; i++ {
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
		for a := 0; a < n; a++ {
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
				for c := 0; c < n; c++ {
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
func layoutRingOrHub(layout *PreviewLayout, zones []models.RmgZone, conns []models.RmgConnection, side float64) {
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

	// Hub detection: explicit "Hub" zone, or a single non-player zone
	// connected to every player zone (implicit hub-and-spoke).
	hubIdx := -1
	for i, z := range zones {
		if z.Name == "Hub" {
			hubIdx = i
			break
		}
	}
	if hubIdx < 0 {
		hubName := findImplicitHubName(zones, conns)
		if hubName != "" {
			for i, z := range zones {
				if z.Name == hubName {
					hubIdx = i
					break
				}
			}
		}
	}

	var outer []int
	for i := range zones {
		if i != hubIdx {
			outer = append(outer, i)
		}
	}
	outerN := len(outer)
	if outerN < 1 {
		outerN = 1
	}
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

func layoutMultiHub(layout *PreviewLayout, zones []models.RmgZone, conns []models.RmgConnection, hubIndices []int, side float64) {
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

// findImplicitHubName returns the single non-player zone connected (Direct
// edges) to every player zone, or "" if no such zone exists. Used to render
// shared hubs that were not literally named "Hub".
func findImplicitHubName(zones []models.RmgZone, conns []models.RmgConnection) string {
	playerNames := map[string]bool{}
	for _, z := range zones {
		if strings.HasPrefix(z.Name, "Spawn-") {
			playerNames[z.Name] = true
		}
	}
	if len(playerNames) == 0 {
		return ""
	}
	neighbours := make(map[string]map[string]bool, len(zones))
	for _, c := range conns {
		if isStructuralIgnored(c.ConnectionType) {
			continue
		}
		if neighbours[c.From] == nil {
			neighbours[c.From] = map[string]bool{}
		}
		if neighbours[c.To] == nil {
			neighbours[c.To] = map[string]bool{}
		}
		neighbours[c.From][c.To] = true
		neighbours[c.To][c.From] = true
	}
	bestName := ""
	bestDeg := -1
	for _, z := range zones {
		if strings.HasPrefix(z.Name, "Spawn-") {
			continue
		}
		nb := neighbours[z.Name]
		if len(nb) < 2 {
			continue
		}
		connectsAll := true
		for p := range playerNames {
			if !nb[p] {
				connectsAll = false
				break
			}
		}
		if !connectsAll {
			continue
		}
		if len(nb) > bestDeg {
			bestDeg = len(nb)
			bestName = z.Name
		}
	}
	return bestName
}

// ── Tier / letter helpers (kept for compatibility) ────────────────────────

// ExtractZoneLetter returns the trailing letter portion of a zone name like
// "Spawn-A" → "A" or "Neutral-C" → "C". Plain names (e.g. "Hub") pass through.
func ExtractZoneLetter(zoneName string) string {
	if strings.HasPrefix(zoneName, "Spawn-") {
		return strings.TrimPrefix(zoneName, "Spawn-")
	}
	if strings.HasPrefix(zoneName, "Neutral-") {
		return strings.TrimPrefix(zoneName, "Neutral-")
	}
	return zoneName
}

// ClassifyZoneTier guesses a neutral zone's tier from its layout name and
// falls back to keyword matching on the zone name
func ClassifyZoneTier(zone models.RmgZone) int {
	if strings.HasPrefix(zone.Name, "Spawn-") {
		return 0
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

func connectedComponents(n int, adj [][]int) [][]int {
	id := make([]int, n)
	for i := range id {
		id[i] = -1
	}
	var comps [][]int
	for start := 0; start < n; start++ {
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

func positionCentroid(zones []models.RmgZone) (float64, float64) {
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

func positionAngle(z models.RmgZone, rawCx, rawCy float64) float64 {
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
