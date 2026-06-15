package topology

import (
	"math"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
)

// GeometricTopologyService arranges zones into symmetric concentric polygons: a
// player ring on the outside, a neutral ring on the inside and a central zone,
// joined by ring loops and radial spokes so the connections trace clean
// geometric shapes (petals, wheels and stars).
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

// createGeometricLayout builds a centrally symmetric figure: a hub neutral zone
// at the centre, the remaining neutral zones on an inner ring forming the
// central polygon (joined in a loop and spoked to the hub like a star), and the
// player zones on an outer ring, each hanging off the figure via a single
// radial petal stem to its nearest inner-ring neutral. Players never border one
// another, so the figure stays intact under isolated-player-starts and reads
// like the example templates (Shamrock, One for All, Nuclear, Kerberos).
func (this *GeometricTopologyService) createGeometricLayout(
	playerLabels []string,
	neutralZones models.NeutralZonePlans) ([]string, models.Positions, [][2]int) {
	const (
		centreX     = 0.5
		centreY     = 0.5
		outerRadius = 0.42
		innerRadius = 0.20
		startAngle  = -math.Pi / 2.0
	)
	playerCount := len(playerLabels)
	neutralCount := len(neutralZones)

	var allLabels []string
	var positions models.Positions

	// Hub: the first neutral zone anchors the centre of the figure.
	centreIndex := -1
	neutralCursor := 0
	if neutralCount >= 1 {
		centreIndex = len(allLabels)
		allLabels = append(allLabels, neutralZones[0].Label)
		positions.Add(models.NewPosition(centreX, centreY))
		neutralCursor = 1
	}

	// Inner ring: every remaining neutral zone, evenly spaced, forming the
	// central polygon.
	innerStart := len(allLabels)
	innerCount := neutralCount - neutralCursor
	for j := 0; j < innerCount; j++ {
		angle := startAngle + 2.0*math.Pi*float64(j)/float64(innerCount)
		allLabels = append(allLabels, neutralZones[neutralCursor+j].Label)
		positions.Add(circlePoint(angle, centreX, centreY, innerRadius))
	}

	// Outer ring: player zones, aligned to the inner-ring spokes so the petal
	// stems run radially.
	playerStart := len(allLabels)
	for i, label := range playerLabels {
		angle := startAngle
		if playerCount > 0 {
			angle += 2.0 * math.Pi * float64(i) / float64(playerCount)
		}
		allLabels = append(allLabels, label)
		positions.Add(circlePoint(angle, centreX, centreY, outerRadius))
	}

	pairs := this.createGeometricPairs(centreIndex, innerStart, innerCount, playerStart, playerCount, positions)
	return allLabels, positions, pairs
}

func (this *GeometricTopologyService) createGeometricPairs(
	centreIndex, innerStart, innerCount, playerStart, playerCount int,
	positions models.Positions) [][2]int {
	builder := newPairBuilder()

	// Central polygon joining the inner neutral ring.
	if innerCount >= 2 {
		for j := 0; j < innerCount; j++ {
			builder.add(innerStart+j, innerStart+(j+1)%innerCount)
		}
	}

	// Hub star: spokes from the centre out to every inner-ring node (or straight
	// to the players when there is no inner ring).
	if centreIndex >= 0 {
		if innerCount > 0 {
			for j := 0; j < innerCount; j++ {
				builder.add(centreIndex, innerStart+j)
			}
		} else {
			for i := 0; i < playerCount; i++ {
				builder.add(centreIndex, playerStart+i)
			}
		}
	}

	// Player petal stems: each player hangs off its nearest inner-ring neutral,
	// or the hub when the figure has no inner ring.
	for i := 0; i < playerCount; i++ {
		playerIndex := playerStart + i
		switch {
		case innerCount > 0:
			nearest := nearestIndexInRange(positions, playerIndex, innerStart, innerStart+innerCount)
			if nearest >= 0 {
				builder.add(playerIndex, nearest)
			}
		case centreIndex >= 0:
			builder.add(playerIndex, centreIndex)
		}
	}

	// With no neutral zones at all there is no central figure, so fall back to a
	// closed player polygon to keep a geometric outline.
	if centreIndex < 0 && playerCount >= 2 {
		for i := 0; i < playerCount; i++ {
			builder.add(playerStart+i, playerStart+(i+1)%playerCount)
		}
	}
	return builder.pairs
}
