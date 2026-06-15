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

// petal holds the zone indices that make up one flower petal: the neutral zones
// forming the left and right edges of the leaf (ordered hub→tip) and the player
// zone at the tip.
type petal struct {
	left   []int
	right  []int
	player int
}

// createGeometricLayout builds the flower: a hub neutral at the centre and one
// petal per player. The remaining neutral zones are split evenly between the
// petals and arranged into a leaf outline — each side of the leaf curves out
// from near the hub to the player at the tip — so the connection graph reads as
// a ring of rounded petals (Shamrock, One for All, Kerberos, Infinity).
func (this *GeometricTopologyService) createGeometricLayout(
	playerLabels []string,
	neutralZones models.NeutralZonePlans) ([]string, models.Positions, [][2]int) {
	const (
		centreX    = 0.5
		centreY    = 0.5
		baseRadius = 0.10 // where the petal edges start, just off the hub
		tipRadius  = 0.44 // the player zone at the petal tip
		startAngle = -math.Pi / 2.0
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
	// every petal grows at the same rate and zone tiers stay balanced.
	petalPlans := make([][]int, playerCount)
	for offset := 0; neutralCursor+offset < neutralCount && playerCount > 0; offset++ {
		petalPlans[offset%playerCount] = append(petalPlans[offset%playerCount], neutralCursor+offset)
	}

	// Keep each petal narrow enough that neighbours never overlap: the half-width
	// stays under half the angular gap between adjacent players.
	halfGap := math.Pi
	if playerCount >= 2 {
		halfGap = math.Pi / float64(playerCount)
	}
	maxOffset := halfGap * 0.62

	petals := make([]petal, playerCount)
	for index := 0; index < playerCount; index++ {
		centreAngle := startAngle + 2.0*math.Pi*float64(index)/float64(playerCount)
		plan := petalPlans[index]
		leftCount := (len(plan) + 1) / 2

		// placeSide lays neutrals along one edge of the leaf. The angular offset
		// follows sin(pi*t) so the edge bows out at the middle and closes in at
		// both the hub (t→0) and the tip (t→1), tracing a leaf silhouette.
		placeSide := func(planSlice []int, sign float64) []int {
			indices := make([]int, 0, len(planSlice))
			total := len(planSlice)
			for step, planIndex := range planSlice {
				t := float64(step+1) / float64(total+1)
				radius := baseRadius + (tipRadius-baseRadius)*t
				angle := centreAngle + sign*maxOffset*math.Sin(math.Pi*t)
				indices = append(indices, len(allLabels))
				allLabels = append(allLabels, neutralZones[planIndex].Label)
				positions.Add(circlePoint(angle, centreX, centreY, radius))
			}
			return indices
		}

		left := placeSide(plan[:leftCount], -1.0)
		right := placeSide(plan[leftCount:], +1.0)

		playerIndex := len(allLabels)
		allLabels = append(allLabels, playerLabels[index])
		positions.Add(circlePoint(centreAngle, centreX, centreY, tipRadius))

		petals[index] = petal{left: left, right: right, player: playerIndex}
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
		// Climb the left edge from the hub up to the player tip.
		previous := centreIndex
		for _, node := range p.left {
			if previous >= 0 {
				builder.add(previous, node)
			}
			previous = node
		}
		if previous >= 0 {
			builder.add(previous, p.player)
		}

		// Descend the right edge from the tip back down to the hub, closing the
		// petal loop.
		previous = p.player
		for step := len(p.right) - 1; step >= 0; step-- {
			builder.add(previous, p.right[step])
			previous = p.right[step]
		}
		if centreIndex >= 0 {
			builder.add(previous, centreIndex)
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
