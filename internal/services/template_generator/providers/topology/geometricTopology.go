package topology

import (
	"math"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	zone_services "github.com/Tariomka/hommoe_custom_templates/internal/services/zones"
)

// GeometricTopologyService arranges zones into a centrally symmetric flower: a
// hub zone at the center with one petal per player radiating outward. Each petal
// is a closed leaf-shaped loop of neutral zones that bulges out to the player
// zone at its tip and curves back to the hub, so the connections trace the
// rounded petals of the example templates (Shamrock, One for All, Nuclear,
// Kerberos, Infinity).
type GeometricTopologyService struct {
	PositionedTopologyBuilder
}

func NewGeometricTopologyService() *GeometricTopologyService {
	return NewGeometricTopologyServiceWithCreationServices(zone_services.NewCreationServices(nil, nil))
}

func NewGeometricTopologyServiceWithCreationServices(
	creationServices *zone_services.CreationServices,
) *GeometricTopologyService {
	return &GeometricTopologyService{
		PositionedTopologyBuilder: *NewPositionedTopologyBuilderWithCreationServices(creationServices),
	}
}

func (this *GeometricTopologyService) CreateTopologyVariant(
	configuration config.GeneratorConfig,
	playerLabels []string,
	neutralZones neutral_zone.Plans,
	tuning models.GenerationTuning,
	holdCityNeutralLabel string) entities.Variant {
	return this.BuildVariant(
		configuration, playerLabels, neutralZones, tuning, holdCityNeutralLabel, this.createGeometricLayout, nil)
}

// petal holds the zone indices that make up one flower petal in ring order:
// the neutral zones tracing the rounded lobe outline (from the center-facing
// base on one side, around the outer rim, to the base on the other side) with
// the player zone sitting at the outer tip in the middle of that order.
type petal struct {
	ring   []int // ordered base → ... → tip(player) → ... → base
	player int
}

// createGeometricLayout builds the flower: a central neutral at the middle and
// one petal per player. Each petal is a fat teardrop (leaf) that fans out from
// the center to fill the player's whole angular sector, reaching almost to the
// canvas edge where the player zone sits at the tip. The neutrals trace the two
// bowed edges of the leaf so the lobes are large and space-filling - like the
// example templates (Shamrock, One for All, Nuclear, Kerberos, Infinity) which
// are built purely from player and neutral zones with no dedicated hub.
func (this *GeometricTopologyService) createGeometricLayout(
	playerLabels []string,
	neutralZones neutral_zone.Plans) ([]string, models.Positions, []models.ConnectionIndexes) {
	const (
		tipRadius  = 0.46 // player tip, almost at the canvas edge
		startAngle = -math.Pi / 2.0
	)
	playerCount := len(playerLabels)
	neutralCount := len(neutralZones)

	var allLabels []string
	var positions models.Positions

	// Center: the first neutral zone anchors the middle of the flower. It is a
	// regular neutral zone - never a dedicated hub - so it only ever borders the
	// petals' base neutrals, not the players.
	centerIndex := -1
	neutralCursor := 0
	if neutralCount >= 1 {
		centerIndex = len(allLabels)
		allLabels = append(allLabels, neutralZones[0].Label)
		positions.Add(data.NewVec2(layoutCenter, layoutCenter))
		neutralCursor = 1
	}

	// Distribute the remaining neutral zones across the petals round-robin so
	// every lobe grows at the same rate and zone tiers stay balanced.
	petalPlans := make([][]int, playerCount)
	for offset := 0; neutralCursor+offset < neutralCount && playerCount > 0; offset++ {
		petalPlans[offset%playerCount] = append(petalPlans[offset%playerCount], neutralCursor+offset)
	}

	// Each petal fills its angular sector. The leaf edges bow out almost to the
	// sector boundary so neighbouring petals nearly touch, leaving no large gaps.
	sectorHalf := math.Pi
	if playerCount >= 2 {
		sectorHalf = math.Pi / float64(playerCount)
	}
	bowAngle := sectorHalf * 0.92
	// Bézier control distance: pushing the control point out past the tip radius
	// fattens the leaf so it bulges toward the edge instead of staying a thin
	// sliver. Wider sectors (few players) get rounder leaves.
	ctrlDist := tipRadius * (0.82 + 0.30*math.Min(1.0, sectorHalf/(math.Pi/3.0)))

	petals := make([]petal, playerCount)
	for index := range playerCount {
		axis := startAngle + 2.0*math.Pi*float64(index)/float64(playerCount)
		petals[index] = buildPetal(
			axis, bowAngle, ctrlDist, tipRadius,
			playerLabels[index], petalPlans[index], neutralZones,
			&allLabels, &positions)
	}

	pairs := this.createGeometricPairs(centerIndex, petals, playerCount)
	return allLabels, positions, pairs
}

// buildPetal places one player's fat teardrop leaf: the neutral zones trace the
// two bowed Bézier edges from the center-facing bases out to the tip, where the
// player zone sits almost at the canvas edge. Labels and positions are appended
// in place; the returned petal records the ring order base → tip → base.
func buildPetal(
	axis, bowAngle, ctrlDist, tipRadius float64,
	playerLabel string,
	plan []int,
	neutralZones neutral_zone.Plans,
	allLabels *[]string,
	positions *models.Positions) petal {
	var ring []int
	addNode := func(planIndex int, point models.Position) {
		ring = append(ring, len(*allLabels))
		*allLabels = append(*allLabels, neutralZones[planIndex].Label)
		positions.Add(point)
	}
	createPosition := func(angle, radius float64) models.Position {
		return data.NewVec2(
			layoutCenter+math.Cos(angle)*radius,
			layoutCenter+math.Sin(angle)*radius)
	}

	startPosition := data.NewVec2(layoutCenter, layoutCenter)
	endPosition := createPosition(axis, tipRadius)
	leftCtrl := axis + bowAngle
	rightCtrl := axis - bowAngle
	leftCount := (len(plan) + 1) / 2
	rightCount := len(plan) - leftCount

	// Left edge: from the base near the center out toward the tip.
	for index := 1; index <= leftCount; index++ {
		t := float64(index) / float64(leftCount+1)
		addNode(plan[index-1], helpers.GetVectorOnQuadraticBezierCurve(
			startPosition, createPosition(leftCtrl, ctrlDist), endPosition, t))
	}
	// Tip: the player zone, almost at the canvas edge.
	playerIndex := len(*allLabels)
	*allLabels = append(*allLabels, playerLabel)
	positions.Add(endPosition)
	ring = append(ring, playerIndex)
	// Right edge: from the tip back down to the base near the center.
	for index := rightCount; index >= 1; index-- {
		t := float64(index) / float64(rightCount+1)
		addNode(plan[leftCount+rightCount-index],
			helpers.GetVectorOnQuadraticBezierCurve(
				startPosition, createPosition(rightCtrl, ctrlDist), endPosition, t))
	}

	return petal{ring: ring, player: playerIndex}
}

func (this *GeometricTopologyService) createGeometricPairs(
	centerIndex int,
	petals []petal,
	playerCount int) []models.ConnectionIndexes {
	builder := newPairBuilder()

	for _, p := range petals {
		// Trace the rounded lobe outline (left rim → player tip → right rim).
		for step := 0; step+1 < len(p.ring); step++ {
			builder.add(p.ring[step], p.ring[step+1])
		}
		if centerIndex < 0 || len(p.ring) == 0 {
			continue
		}

		// Anchor the petal to the center through its base neutrals - the rim ends
		// nearest the center gap. The center must never border a player tip, or
		// it would connect to every player and read as a dedicated hub, so a base
		// that is the player itself (a petal with no neutrals on that side) is
		// skipped in favour of the opposite base.
		first, last := p.ring[0], p.ring[len(p.ring)-1]
		anchored := false
		if first != p.player {
			builder.add(centerIndex, first)
			anchored = true
		}
		if last != p.player && last != first {
			builder.add(centerIndex, last)
			anchored = true
		}
		// Degenerate petal with no neutrals at all: the player has to hang off
		// the center directly. Only reachable when there are fewer neutrals than
		// players, an extreme low-zone configuration.
		if !anchored {
			builder.add(centerIndex, p.player)
		}
	}

	// No neutral zones at all: there is no center to build petals around, so fall
	// back to a closed player polygon to keep a geometric outline.
	if centerIndex < 0 && playerCount >= 2 {
		ring := make([]int, playerCount)
		for index, p := range petals {
			ring[index] = p.player
		}
		connectRing(builder, ring)
	}
	return builder.pairs
}

// connectRing joins the given node indices into a closed loop (or a single edge
// for a pair). Rings of fewer than two nodes are left untouched.
func connectRing(builder *pairBuilder, ring []int) {
	switch len(ring) {
	case 0, 1:
		return
	case 2:
		builder.add(ring[0], ring[1])
	default:
		for i := range ring {
			builder.add(ring[i], ring[(i+1)%len(ring)])
		}
	}
}
