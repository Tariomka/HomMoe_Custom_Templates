package preview_service

import (
	"math"
	"slices"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
)

// layoutScatter renders the organic position-driven topologies (Random and
// Circles without ring stamps): the raw GeneratorPosition stamps are scaled so
// the mean direct-edge length is readable, then relaxation passes push
// overlapping zones apart and nudge zones off connection lines before a final
// shrink-to-fit.
func (this *PreviewLayoutService) layoutScatter(
	zones []entities.Zone,
	conns []entities.Connection,
	side float64) {
	metrics := newCanvasMetrics(side)
	if this.placeTrivial(zones, metrics) {
		return
	}

	adj := buildScatterAdjacency(zones, conns)
	zoneRadius := scatterZoneRadius(adj, metrics)
	px, py := projectScatterPositions(zones, adj, zoneRadius, metrics)

	relaxPasses(px, py, adj, zoneRadius)

	// Final fit: recentre the bounding box then shrink only if it overflows
	// the padded canvas.
	shrink := fitToCanvas(px, py, metrics, zoneRadius+metrics.margin, false)
	if shrink < 1.0 {
		zoneRadius = math.Max(zoneRadius*shrink, csZoneRadiusFloor)
	}
	this.commitPositions(zones, px, py, zoneRadius)
}

// buildScatterAdjacency builds the direct-only, deduplicated adjacency lists
// that drive both the radius heuristic and the edge-clearance pass.
func buildScatterAdjacency(zones []entities.Zone, conns []entities.Connection) [][]int {
	idx := make(map[string]int, len(zones))
	for i, zone := range zones {
		idx[zone.Name] = i
	}
	adj := make([][]int, len(zones))
	for _, conn := range conns {
		if isStructuralIgnored(conn.ConnectionType) {
			continue
		}
		a, ok1 := idx[conn.From]
		b, ok2 := idx[conn.To]
		if !ok1 || !ok2 || a == b || slices.Contains(adj[a], b) {
			continue
		}
		adj[a] = append(adj[a], b)
		adj[b] = append(adj[b], a)
	}
	return adj
}

// scatterZoneRadius sizes zones off the largest connected component so
// readability stays consistent regardless of how many neutrals exist.
func scatterZoneRadius(adj [][]int, metrics canvasMetrics) float64 {
	maxCompSize := 1
	for _, comp := range connectedComponents(len(adj), adj) {
		maxCompSize = max(maxCompSize, len(comp))
	}
	if maxCompSize < 2 {
		return metrics.zoneRadiusMax
	}
	ringRadius := metrics.side/2.0 - metrics.margin
	chord := 2.0 * ringRadius * math.Sin(math.Pi/float64(maxCompSize))
	radius := math.Min(metrics.zoneRadiusMax, (chord-metrics.minGap)/2.0)
	return math.Max(radius, csZoneRadiusFloor)
}

// projectScatterPositions maps the raw [0,1] generator positions onto the
// canvas: scaled so the mean direct-edge length matches the ideal for the
// zone radius (falling back to spanning the draw area for empty graphs),
// centred, and shrunk to fit the padded canvas.
func projectScatterPositions(
	zones []entities.Zone,
	adj [][]int,
	zoneRadius float64,
	metrics canvasMetrics) (px, py []float64) {
	gScale := metrics.side - 2.0*metrics.margin
	if mean := meanRawEdgeLength(zones, adj); mean > 1e-6 {
		gScale = zoneRadius * scatterIdealMult / mean
	}

	rawCx, rawCy := positionCentroid(zones)
	px = make([]float64, len(zones))
	py = make([]float64, len(zones))
	for i, zone := range zones {
		p := *zone.GeneratorPosition
		px[i] = (p[0] - rawCx) * gScale
		py[i] = (p[1] - rawCy) * gScale
	}
	fitToCanvas(px, py, metrics, zoneRadius+metrics.margin, false)
	return px, py
}

// meanRawEdgeLength averages the raw generator-space length of the direct
// edges; returns 0 when the graph has none.
func meanRawEdgeLength(zones []entities.Zone, adj [][]int) float64 {
	sum, count := 0.0, 0
	for i := range adj {
		for _, j := range adj[i] {
			if j <= i {
				continue
			}
			pi := *zones[i].GeneratorPosition
			pj := *zones[j].GeneratorPosition
			sum += math.Hypot(pi[0]-pj[0], pi[1]-pj[1])
			count++
		}
	}
	if count == 0 {
		return 0
	}
	return sum / float64(count)
}

// relaxPasses alternates two correction passes until the layout settles: a
// hard-floor pass that pushes overlapping zones apart and an edge-clearance
// pass that nudges zones off connection lines.
func relaxPasses(px, py []float64, adj [][]int, zoneRadius float64) {
	minDist := zoneRadius * scatterMinDist
	edgeClear := zoneRadius * scatterEdgeClear
	for range 500 {
		pushed := pushApartPass(px, py, minDist)
		nudged := nudgeOffEdgesPass(px, py, adj, edgeClear)
		if !pushed && !nudged {
			break
		}
	}
}

// pushApartPass symmetrically pushes every zone pair closer than minDist
// apart. Reports whether anything moved.
func pushApartPass(px, py []float64, minDist float64) bool {
	moved := false
	for i := range px {
		for j := i + 1; j < len(px); j++ {
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
	return moved
}

// nudgeOffEdgesPass moves any zone that sits too close to another pair's
// connection line perpendicularly off it. Reports whether anything moved.
func nudgeOffEdgesPass(px, py []float64, adj [][]int, edgeClear float64) bool {
	moved := false
	for a := range px {
		for _, b := range adj[a] {
			if b <= a {
				continue
			}
			if nudgeZonesOffEdge(px, py, adj, a, b, edgeClear) {
				moved = true
			}
		}
	}
	return moved
}

// nudgeZonesOffEdge checks every third zone against the a-b segment and moves
// the ones within edgeClear to whichever perpendicular side sits farther from
// their own neighbours. Reports whether anything moved.
func nudgeZonesOffEdge(px, py []float64, adj [][]int, a, b int, edgeClear float64) bool {
	ex := px[b] - px[a]
	ey := py[b] - py[a]
	elen2 := ex*ex + ey*ey
	if elen2 < 1e-3 {
		return false
	}
	elenInv := 1.0 / math.Sqrt(elen2)

	moved := false
	for c := range px {
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
		px[c], py[c] = preferredNudge(px, py, adj[c],
			projX+perpX*edgeClear, projY+perpY*edgeClear,
			projX-perpX*edgeClear, projY-perpY*edgeClear)
		moved = true
	}
	return moved
}

// preferredNudge picks whichever candidate point sits farther (in summed
// squared distance) from the zone's own neighbours, keeping its edges long.
func preferredNudge(px, py []float64, neighbours []int, ax, ay, bx, by float64) (float64, float64) {
	scoreA, scoreB := 0.0, 0.0
	for _, nb := range neighbours {
		dax := ax - px[nb]
		day := ay - py[nb]
		dbx := bx - px[nb]
		dby := by - py[nb]
		scoreA += dax*dax + day*day
		scoreB += dbx*dbx + dby*dby
	}
	if scoreB < scoreA {
		return bx, by
	}
	return ax, ay
}
