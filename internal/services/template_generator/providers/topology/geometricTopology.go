package topology

import (
	"math"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
)

// GeometricTopologyService arranges zones into a centrally symmetric flower: a
// hub zone at the centre with one petal per player radiating outward. Each petal
// is a closed leaf-shaped loop of neutral zones that bulges out to the player
// zone at its tip and curves back to the hub, so the connections trace the
// rounded petals of the example templates (Shamrock, One for All, Nuclear,
// Kerberos, Infinity).
type GeometricTopologyService struct {
	RandomTopologyService
}

func NewGeometricTopologyService() *GeometricTopologyService {
	return &GeometricTopologyService{
		RandomTopologyService: *NewRandomTopologyService(),
	}
}

func (this *GeometricTopologyService) CreateTopologyVariant(
	configuration config.GeneratorConfig,
	playerLabels []string,
	neutralZones models.NeutralZonePlans,
	tuning models.GenerationTuning,
	holdCityNeutralLabel string) entities.Variant {
	isIsolated := configuration.NoDirectPlayerConnections && len(playerLabels) > 1
	allLabels, positions, pairs := this.createGeometricLayout(playerLabels, neutralZones)
	connectionNames := this.createConnectionNames(playerLabels, allLabels, pairs, isIsolated)

	zones := this.createZones(configuration, playerLabels, allLabels, tuning, neutralZones, holdCityNeutralLabel, connectionNames)
	for index := range zones {
		position := positions[index]
		zones[index].GeneratorPosition = &[2]float64{position.X, position.Y}
	}

	conns := this.createConnections(playerLabels, allLabels, tuning, isIsolated, neutralZones, connectionNames, pairs)
	if configuration.RandomPortals {
		conns = append(conns, this.CreateRandomPortalConnections(playerLabels, allLabels, tuning, configuration.MaxPortalConnections)...)
	}
	if isIsolated {
		conns = append(conns, this.CreateMissingPlayerConnections(playerLabels, zones, conns, tuning)...)
	}
	conns = append(conns, this.CreateMissingConnections(playerLabels, allLabels, positions, zones, conns, tuning, neutralZones)...)
	return this.CreateVariant(playerLabels, allLabels[0], len(allLabels), zones, conns)
}

// petal holds the zone indices that make up one flower petal in ring order:
// the neutral zones tracing the rounded lobe outline (from the hub-facing base
// on one side, around the outer rim, to the base on the other side) with the
// player zone sitting at the outer tip in the middle of that order.
type petal struct {
	ring   []int // ordered base → … → tip(player) → … → base
	player int
}

// createGeometricLayout builds the flower: a hub neutral at the centre and one
// rounded petal (lobe) per player. Each petal's neutrals are placed on a circle
// (the lobe) so the outline is genuinely round — like the leaves of a clover —
// and the player sits at the outer tip. Every petal is closed back to the hub
// on both sides and split by a central "vein" spoke, so each player has three
// guarded routes inward, matching the rich, rounded figures of the example
// templates (Shamrock, One for All, Nuclear, Kerberos, Infinity).
func (this *GeometricTopologyService) createGeometricLayout(
	playerLabels []string,
	neutralZones models.NeutralZonePlans) ([]string, models.Positions, [][2]int) {
	const (
		centreX    = 0.5
		centreY    = 0.5
		lobeDist   = 0.255 // distance from hub to the centre of each lobe
		startAngle = -math.Pi / 2.0
		openHalf   = 42.0 * math.Pi / 180.0 // half-angle of the hub-facing gap
	)
	playerCount := len(playerLabels)
	neutralCount := len(neutralZones)

	var allLabels []string
	var positions models.Positions

	// Hub: the first neutral zone anchors the centre of the flower.
	centreIndex := -1
	neutralCursor := 0
	if neutralCount >= 1 {
		centreIndex = len(allLabels)
		allLabels = append(allLabels, neutralZones[0].Label)
		positions.Add(models.NewPosition(centreX, centreY))
		neutralCursor = 1
	}

	// Distribute the remaining neutral zones across the petals round-robin so
	// every lobe grows at the same rate and zone tiers stay balanced.
	petalPlans := make([][]int, playerCount)
	for offset := 0; neutralCursor+offset < neutralCount && playerCount > 0; offset++ {
		petalPlans[offset%playerCount] = append(petalPlans[offset%playerCount], neutralCursor+offset)
	}

	// Size each lobe to fill its angular sector without letting neighbours
	// collide: the lobe radius is bounded by how wide the sector is.
	sectorHalf := math.Pi
	if playerCount >= 2 {
		sectorHalf = math.Pi / float64(playerCount)
	}
	lobeRadius := math.Min(0.205, lobeDist*math.Sin(sectorHalf*0.85))
	arcMax := math.Pi - openHalf // how far around the lobe the rim wraps

	petals := make([]petal, playerCount)
	for index := 0; index < playerCount; index++ {
		axis := startAngle + 2.0*math.Pi*float64(index)/float64(playerCount)
		lobeX := centreX + math.Cos(axis)*lobeDist
		lobeY := centreY + math.Sin(axis)*lobeDist
		plan := petalPlans[index]
		neutralPerSide := len(plan)
		leftCount := (neutralPerSide + 1) / 2
		rightCount := neutralPerSide - leftCount

		// lobePoint places a node on the lobe circle at local angle alpha,
		// measured from the outward (tip) direction so alpha=0 is the tip.
		lobePoint := func(alpha float64) models.Vector2 {
			return models.NewPosition(
				lobeX+math.Cos(axis+alpha)*lobeRadius,
				lobeY+math.Sin(axis+alpha)*lobeRadius)
		}

		var ring []int
		addNode := func(planIndex int, alpha float64) {
			ring = append(ring, len(allLabels))
			allLabels = append(allLabels, neutralZones[planIndex].Label)
			positions.Add(lobePoint(alpha))
		}

		// Left rim: outermost (near the hub gap) down to the node beside the tip.
		for j := leftCount; j >= 1; j-- {
			addNode(plan[j-1], -arcMax*float64(j)/float64(leftCount))
		}
		// Tip: the player zone at the outermost point of the lobe.
		playerIndex := len(allLabels)
		allLabels = append(allLabels, playerLabels[index])
		positions.Add(lobePoint(0))
		ring = append(ring, playerIndex)
		// Right rim: node beside the tip out to the outermost (near the hub gap).
		for j := 1; j <= rightCount; j++ {
			addNode(plan[leftCount+j-1], arcMax*float64(j)/float64(rightCount))
		}

		petals[index] = petal{ring: ring, player: playerIndex}
	}

	pairs := this.createGeometricPairs(centreIndex, petals, playerCount)
	return allLabels, positions, pairs
}

func (this *GeometricTopologyService) createGeometricPairs(
	centreIndex int,
	petals []petal,
	playerCount int) [][2]int {
	builder := newPairBuilder()

	for _, p := range petals {
		// Trace the rounded lobe outline.
		for step := 0; step+1 < len(p.ring); step++ {
			builder.add(p.ring[step], p.ring[step+1])
		}
		if centreIndex >= 0 && len(p.ring) > 0 {
			// Close the petal back to the hub on both sides so it reads as a
			// distinct leaf, and run a central vein straight to the player tip —
			// giving every player three guarded routes inward (left rim, vein,
			// right rim) for strategic depth.
			builder.add(centreIndex, p.ring[0])
			builder.add(centreIndex, p.ring[len(p.ring)-1])
			builder.add(centreIndex, p.player)
		}
	}

	// No neutral zones at all: there is no hub to build petals around, so fall
	// back to a closed player polygon to keep a geometric outline.
	if centreIndex < 0 && playerCount >= 2 {
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
