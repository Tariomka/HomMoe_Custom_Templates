package topology

import (
	"math"
	"slices"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
)

// geometricHubLayout is the computed slot structure of the Geometric Hub
// topology: which neutral plan occupies which hexagon slot, the normalized
// position of every zone, the direct ring/interior edges, and the labels that
// connect to the hub through portals.
type geometricHubLayout struct {
	// gapStables[j] holds the stable-zone labels between player j and player
	// j+1 (0, 1 shared, or 2 labels; index 0 is adjacent to player j).
	gapStables [][]string
	// gapCorners[j] holds the merged corner-zone label of the same gap (0 or 1
	// labels).
	gapCorners [][]string
	// hexInteriors[i] holds the interior-zone labels of player i's hexagon in
	// polygon-vertex order: index m is vertex x(m+1) of the regular k-gon
	// (x1 = hub-side-left, x2 = hub-side-right, x3/x4 the next left/right
	// pair, and so on toward the player).
	hexInteriors [][]string
	// positions maps every player and neutral label to its normalized layout
	// position.
	positions map[string]models.Position
	// directEdges lists the Direct connections as label pairs.
	directEdges [][2]string
	// hubPortalLabels lists every label that gets a Portal connection to the
	// hub (corners, hub-adjacent stables, isolated players, interiors).
	hubPortalLabels []string
}

// newGeometricHubLayout distributes the neutral plans over the hexagon slots
// following the confirmed growth order: stables (2 per player) -> merged
// corners (1 per adjacent-player pair) -> interiors (round-robin per hexagon,
// uncapped). Interiors take the highest-quality plans, corners the lowest,
// stables the middle.
func newGeometricHubLayout(playerLabels []string, plans neutral_zone.Plans) *geometricHubLayout {
	playerCount := len(playerLabels)
	layout := &geometricHubLayout{
		gapStables:   make([][]string, playerCount),
		gapCorners:   make([][]string, playerCount),
		hexInteriors: make([][]string, playerCount),
		positions:    map[string]models.Position{},
	}

	stableCounts, cornerCounts, interiorCounts := distributeGeometricHubSlots(playerCount, len(plans))
	layout.assignPlans(*neutral_zone.NewNeutralZonePlansSorted(plans), stableCounts, cornerCounts, interiorCounts)
	layout.computePositions(playerLabels)
	layout.buildEdges(playerLabels)
	return layout
}

// distributeGeometricHubSlots splits the neutral budget into per-gap stable
// counts (cap 2), per-gap merged-corner counts (cap 1), and per-hexagon
// interior counts (uncapped), each filled round-robin.
func distributeGeometricHubSlots(playerCount, neutralCount int) (stables, corners, interiors []int) {
	budget := neutralCount
	take := func(maximum int) int {
		taken := min(budget, maximum)
		budget -= taken
		return taken
	}
	stables = fillRoundRobin(playerCount, take(2*playerCount), 2)
	corners = fillRoundRobin(playerCount, take(playerCount), 1)
	interiors = fillRoundRobin(playerCount, budget, budget)
	return stables, corners, interiors
}

// fillRoundRobin deals the budget one slot at a time across bucketCount
// buckets, never exceeding capPerBucket in any bucket.
func fillRoundRobin(bucketCount, budget, capPerBucket int) []int {
	counts := make([]int, bucketCount)
	for budget > 0 {
		dealt := false
		for index := 0; index < bucketCount && budget > 0; index++ {
			if counts[index] < capPerBucket {
				counts[index]++
				budget--
				dealt = true
			}
		}
		if !dealt {
			break
		}
	}
	return counts
}

// assignPlans hands the quality-sorted plans (highest first) to the slots:
// interiors pop from the high end, corners from the low end, stables receive
// the middle, dealt round-robin so qualities spread evenly across the gaps.
// Interiors are dealt in polygon-vertex order (x1 first), so the hub-facing
// vertices of every hexagon carry its best plans.
func (this *geometricHubLayout) assignPlans(
	sortedPlans neutral_zone.Plans,
	stableCounts, cornerCounts, interiorCounts []int) {
	popHighest := func() string {
		label := sortedPlans[0].Label
		sortedPlans = sortedPlans[1:]
		return label
	}

	popLowest := func() string {
		label := sortedPlans[len(sortedPlans)-1].Label
		sortedPlans = sortedPlans[:len(sortedPlans)-1]
		return label
	}

	for round := range maxCount(interiorCounts) {
		for hexagon, count := range interiorCounts {
			if count > round {
				this.hexInteriors[hexagon] = append(this.hexInteriors[hexagon], popHighest())
			}
		}
	}
	for gap, count := range cornerCounts {
		if count > 0 {
			this.gapCorners[gap] = append(this.gapCorners[gap], popLowest())
		}
	}
	for round := range 2 {
		for gap, count := range stableCounts {
			if count > round {
				this.gapStables[gap] = append(this.gapStables[gap], popHighest())
			}
		}
	}
}

// computePositions stamps every player and neutral label with its normalized
// hexagon position (see [newGeoHubGeometry] for the per-player-count shape).
func (this *geometricHubLayout) computePositions(playerLabels []string) {
	playerCount := len(playerLabels)
	sector := 2 * math.Pi / float64(playerCount)
	playerAngle := func(index int) float64 { return -math.Pi/2 + float64(index)*sector }
	geometry := newGeoHubGeometry(playerCount)

	for index, label := range playerLabels {
		this.positions[label] = circlePoint(playerAngle(index), geometry.playerRadius)
	}
	for gap, stables := range this.gapStables {
		gapMidAngle := playerAngle(gap) + sector/2
		if len(stables) == 1 {
			this.positions[stables[0]] = circlePoint(gapMidAngle, geometry.stableRadius)
		} else if len(stables) == 2 {
			this.positions[stables[0]] = circlePoint(playerAngle(gap)+geometry.stableOffset, geometry.stableRadius)
			this.positions[stables[1]] = circlePoint(playerAngle(gap+1)-geometry.stableOffset, geometry.stableRadius)
		}
	}
	for gap, corners := range this.gapCorners {
		if len(corners) == 1 {
			this.positions[corners[0]] = circlePoint(playerAngle(gap)+sector/2, geometry.cornerRadius)
		}
	}
	for hexagon, interiors := range this.hexInteriors {
		this.computeInteriorPositions(playerAngle(hexagon), geometry, interiors)
	}
}

// computeInteriorPositions places a hexagon's k interiors as a regular k-gon
// centered on the hexagon's centroid (radius interiorCenter along the player
// axis), with the x1-x2 edge facing the hub: vertex x(m+1) sits at
// hub-relative angle ±(2⌈(m+1)/2⌉-1)·π/k, odd vertices on the sL side.
// A single interior sits exactly on the hexagon center (rule 7).
func (this *geometricHubLayout) computeInteriorPositions(
	playerAxisAngle float64, geometry geoHubGeometry, interiors []string) {
	vertexCount := len(interiors)
	if vertexCount == 0 {
		return
	}

	center := circlePoint(playerAxisAngle, geometry.interiorCenter)
	if vertexCount == 1 {
		this.positions[interiors[0]] = center
		return
	}

	for index, label := range interiors {
		step := float64(2*((index+2)/2) - 1)
		side := 1.0
		if index%2 == 1 {
			side = -1.0
		}
		vertexAngle := playerAxisAngle + math.Pi + side*step*math.Pi/float64(vertexCount)
		this.positions[label] = data.NewVec2(
			center.X+geometry.interiorCircumradius*math.Cos(vertexAngle),
			center.Y+geometry.interiorCircumradius*math.Sin(vertexAngle))
	}
}

// buildEdges derives the Direct ring/interior edges and the hub portal labels
// from the assigned slots.
func (this *geometricHubLayout) buildEdges(playerLabels []string) {
	playerCount := len(playerLabels)
	seen := map[[2]string]bool{}
	addDirect := func(from, to string) {
		if from == to {
			return
		}

		key := [2]string{min(from, to), max(from, to)}
		if !seen[key] {
			seen[key] = true
			this.directEdges = append(this.directEdges, [2]string{from, to})
		}
	}

	for gap := range playerCount {
		stables := this.gapStables[gap]
		var chain []string
		if len(stables) > 0 {
			chain = append(chain, stables[0])
		}
		chain = append(chain, this.gapCorners[gap]...)
		if len(stables) == 2 {
			chain = append(chain, stables[1])
		}
		if len(chain) == 0 {
			continue // an empty gap never links the two players directly
		}

		addDirect(playerLabels[gap], chain[0])
		for index := 0; index+1 < len(chain); index++ {
			addDirect(chain[index], chain[index+1])
		}
		addDirect(chain[len(chain)-1], playerLabels[(gap+1)%playerCount])

		// Rule 5/11: with corners present only the corners touch the hub;
		// without them the stables take over the hub link.
		if len(this.gapCorners[gap]) > 0 {
			this.hubPortalLabels = append(this.hubPortalLabels, this.gapCorners[gap]...)
		} else {
			this.hubPortalLabels = append(this.hubPortalLabels, this.gapStables[gap]...)
		}
	}

	for index, label := range playerLabels {
		previousGap := (index - 1 + playerCount) % playerCount
		if len(this.gapStables[index]) == 0 && len(this.gapStables[previousGap]) == 0 {
			this.hubPortalLabels = append(this.hubPortalLabels, label)
		}
	}

	this.buildInteriorEdges(addDirect)
}

// buildInteriorEdges wires each hexagon's interiors as a regular k-gon
// centered inside the hexagon. hexInteriors[i][m] is polygon vertex x(m+1):
// x1 = hub-side-left, x2 = hub-side-right, x3/x4 the next left/right pair,
// and so on toward the player (an odd k puts the last vertex on the player
// axis). Ring edges connect angular neighbors only (never diagonals); the
// left stable links x1 and x3, the right stable links x2 and x4. With k=3
// the single player-side vertex x3 serves both stables (the one allowed
// 4-connection exception), and with k=1 the lone x1 links both stables
// (rule 7, splitting the hexagon into rhombuses). Only x1 and x2 portal to
// the hub.
func (this *geometricHubLayout) buildInteriorEdges(addDirect func(from, to string)) {
	for hexagon, interiors := range this.hexInteriors {
		if len(interiors) == 0 {
			continue
		}

		this.hubPortalLabels = append(this.hubPortalLabels, interiors[:min(len(interiors), 2)]...)
		ordered := interiorAngularOrder(interiors)
		for index := range ordered {
			addDirect(ordered[index], ordered[(index+1)%len(ordered)])
		}
		this.connectInteriorStables(hexagon, interiors, addDirect)
	}
}

// interiorAngularOrder returns the interior labels in polygon ring order:
// the odd vertices (x1, x3, x5, ...) descend the hub-left side and the even
// vertices (x2, x4, ...) ascend the hub-right side, so consecutive entries
// (cyclically) are the k-gon's ring neighbors.
func interiorAngularOrder(interiors []string) []string {
	var leftSide, rightSide []string
	for index, label := range interiors {
		if index%2 == 0 {
			leftSide = append(leftSide, label)
		} else {
			rightSide = append(rightSide, label)
		}
	}
	slices.Reverse(leftSide)
	return append(leftSide, rightSide...)
}

// connectInteriorStables links the hexagon's stables to the polygon: the
// left stable to x1 and x3, the right stable to x2 and x4; smaller polygons
// reuse the nearest existing vertex (k=1: both stables to x1; k=3: both
// stables to x3).
func (this *geometricHubLayout) connectInteriorStables(
	hexagon int, interiors []string, addDirect func(from, to string)) {
	stables := this.hexagonStables(hexagon)
	if len(stables) == 0 {
		return
	}

	stableLeft := stables[0]
	stableRight := stables[len(stables)-1]
	addDirect(stableLeft, interiors[0])
	addDirect(stableRight, interiors[min(1, len(interiors)-1)])
	if len(interiors) < 3 {
		return
	}

	addDirect(stableLeft, interiors[2])
	addDirect(stableRight, interiors[min(3, len(interiors)-1)])
}

// hexagonStables returns the stable zones flanking player i: the last stable
// of the gap on the player's left and the first stable of the gap on its
// right (a shared stable serves both roles).
func (this *geometricHubLayout) hexagonStables(hexagon int) []string {
	playerCount := len(this.gapStables)
	var stables []string
	if left := this.gapStables[(hexagon-1+playerCount)%playerCount]; len(left) > 0 {
		stables = append(stables, left[len(left)-1])
	}
	if right := this.gapStables[hexagon]; len(right) > 0 && (len(stables) == 0 || right[0] != stables[0]) {
		stables = append(stables, right[0])
	}
	return stables
}

func maxCount(counts []int) int {
	highest := 0
	for _, count := range counts {
		if count > highest {
			highest = count
		}
	}
	return highest
}
