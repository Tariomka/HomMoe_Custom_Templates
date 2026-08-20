package preview_service

import (
	"math"
	"slices"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
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
	positions := projectScatterPositions(zones, adj, zoneRadius, metrics)

	relaxPasses(positions, adj, zoneRadius)

	// Final fit: recenter the bounding box then shrink only if it overflows
	// the padded canvas.
	shrink := fitToCanvas(positions, metrics, zoneRadius+metrics.margin, false)
	if shrink < 1.0 {
		zoneRadius = math.Max(zoneRadius*shrink, csZoneRadiusFloor)
	}
	this.commitPositions(zones, positions, zoneRadius)
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
// centered, and shrunk to fit the padded canvas.
func projectScatterPositions(
	zones []entities.Zone,
	adj [][]int,
	zoneRadius float64,
	metrics canvasMetrics) models.Positions {
	gScale := metrics.side - 2.0*metrics.margin
	if mean := meanRawEdgeLength(zones, adj); mean > 1e-6 {
		gScale = zoneRadius * scatterIdealMultiplier / mean
	}

	rawCenter := positionCentroid(zones)
	var positions models.Positions
	for _, zone := range zones {
		p := *zone.GeneratorPosition // Is this required to be copied? can't it be used directly safely?
		positions.Add(data.NewVec2(p[0], p[1]).Subtract(rawCenter).MultiplyScalar(gScale))
	}
	fitToCanvas(positions, metrics, zoneRadius+metrics.margin, false)
	return positions
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
func relaxPasses(positions models.Positions, adj [][]int, zoneRadius float64) {
	minDist := zoneRadius * scatterMinDist
	edgeClear := zoneRadius * scatterEdgeClear
	for range 500 {
		pushed := pushApartPass(positions, minDist)
		nudged := nudgeOffEdgesPass(positions, adj, edgeClear)
		if !pushed && !nudged {
			break
		}
	}
}

// pushApartPass symmetrically pushes every zone pair closer than minDist
// apart. Reports whether anything moved.
func pushApartPass(positions models.Positions, minDist float64) bool {
	moved := false
	for i := range positions {
		for j := i + 1; j < len(positions); j++ {
			deltaPosition := positions[i].Subtract(positions[j])
			distance := deltaPosition.Distance()
			if distance >= minDist {
				continue
			}

			if distance < 1e-3 {
				deltaPosition, distance = data.NewVec2(1.0, 0.0), 1e-3
			}
			push := (minDist - distance) / 2.0
			positions[i] = positions[i].Add(deltaPosition.MultiplyScalar(push / distance))
			positions[j] = positions[j].Subtract(deltaPosition.MultiplyScalar(push / distance))
			moved = true
		}
	}
	return moved
}

// nudgeOffEdgesPass moves any zone that sits too close to another pair's
// connection line perpendicularly off it. Reports whether anything moved.
func nudgeOffEdgesPass(positions models.Positions, adj [][]int, edgeClear float64) bool {
	moved := false
	for index := range positions {
		for _, adjacencyValue := range adj[index] {
			if adjacencyValue <= index {
				continue
			}

			if nudgeZonesOffEdge(positions, adj, index, adjacencyValue, edgeClear) {
				moved = true
			}
		}
	}
	return moved
}

// nudgeZonesOffEdge checks every third zone against the a-b segment and moves
// the ones within edgeClear to whichever perpendicular side sits farther from
// their own neighbours. Reports whether anything moved.
func nudgeZonesOffEdge(positions models.Positions, adj [][]int, a, b int, edgeClear float64) bool {
	delta := positions[b].Subtract(positions[a])
	elen2 := delta.SquaredLength()
	if elen2 < 1e-3 {
		return false
	}
	elenInv := 1.0 / math.Sqrt(elen2)

	moved := false
	for c := range positions {
		if c == a || c == b {
			continue
		}

		tProj := positions[c].Subtract(positions[a]).DotProduct(delta) / elen2
		if tProj < 0 || tProj > 1 {
			continue
		}

		projected := positions[a].Add(delta.MultiplyScalar(tProj))
		projectedDelta := positions[c].Subtract(projected)
		distance := projectedDelta.Distance()
		if distance >= edgeClear {
			continue
		}

		var perp data.Vec2[float64]
		if distance < 1e-3 {
			perp = data.NewVec2(delta.Y, -delta.X).MultiplyScalar(elenInv)
		} else {
			perp = projectedDelta.DivideScalar(distance)
		}
		positions[c] = preferredNudge(positions, adj[c],
			projected.Add(perp.MultiplyScalar(edgeClear)),
			projected.Subtract(perp.MultiplyScalar(edgeClear)))
		moved = true
	}
	return moved
}

// preferredNudge picks whichever candidate point sits farther (in summed
// squared distance) from the zone's own neighbors, keeping its edges long.
func preferredNudge(
	positions models.Positions,
	neighbors []int,
	positionA, positionB data.Vec2[float64]) data.Vec2[float64] {
	scoreA, scoreB := 0.0, 0.0
	for _, nb := range neighbors {
		deltaA := positionA.Subtract(positions[nb])
		deltaB := positionB.Subtract(positions[nb])
		scoreA += deltaA.SquaredLength()
		scoreB += deltaB.SquaredLength()
	}
	if scoreB < scoreA {
		return positionB
	}

	return positionA
}
