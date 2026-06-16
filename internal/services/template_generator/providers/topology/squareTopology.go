package topology

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
)

// SquareTopologyService lays the player zones out along the edges of a square.
// The perimeter zones are joined in a closed loop so the connections trace the
// square outline, while a share of the neutral zones is pulled into the middle
// of the square and tied back to the nearest edge.
type SquareTopologyService struct {
	RandomTopologyService
}

func NewSquareTopologyService() *SquareTopologyService {
	return &SquareTopologyService{
		RandomTopologyService: *NewRandomTopologyService(),
	}
}

func (this *SquareTopologyService) CreateTopologyVariant(
	configuration config.GeneratorConfig,
	playerLabels []string,
	neutralZones models.NeutralZonePlans,
	tuning models.GenerationTuning,
	holdCityNeutralLabel string) entities.Variant {
	isIsolated := configuration.NoDirectPlayerConnections && len(playerLabels) > 1
	allLabels, positions, pairs := this.createSquareLayout(playerLabels, neutralZones)
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

// createSquareLayout builds the parallel label, position and connection-pair
// slices for the square. Index 0..perimeterCount-1 are the edge zones (players
// and the neutral zones that share the perimeter); the remaining indices are
// the interior neutral zones.
func (this *SquareTopologyService) createSquareLayout(
	playerLabels []string,
	neutralZones models.NeutralZonePlans) ([]string, models.Positions, [][2]int) {
	const (
		centreX = 0.5
		centreY = 0.5
		half    = 0.42
	)

	// Roughly a third of the neutral zones (the higher tiers, which sort last)
	// are pulled inside the square; the rest line the edges between players.
	interiorCount := len(neutralZones) / 3
	perimeterPlans := neutralZones[:len(neutralZones)-interiorCount]
	interiorPlans := neutralZones[len(neutralZones)-interiorCount:]

	// Even player spacing around the perimeter, separated by the edge neutrals.
	perimeterLabels := this.ZoneLabelProvider.CreateBalancedRingZoneLabels(playerLabels, perimeterPlans, 0)

	var allLabels []string
	var positions models.Positions

	perimeterCount := len(perimeterLabels)
	for i, label := range perimeterLabels {
		fraction := 0.0
		if perimeterCount > 0 {
			fraction = float64(i) / float64(perimeterCount)
		}
		allLabels = append(allLabels, label)
		positions.Add(squarePerimeterPoint(fraction, centreX, centreY, half))
	}

	// Interior neutral zones sit on a smaller inner square (or the exact centre
	// when there is only one) so they read as being inside the perimeter.
	interiorHalf := half * 0.45
	for i, plan := range interiorPlans {
		var point models.Vector2
		if len(interiorPlans) == 1 {
			point = models.NewPosition(centreX, centreY)
		} else {
			point = squarePerimeterPoint(float64(i)/float64(len(interiorPlans)), centreX, centreY, interiorHalf)
		}
		allLabels = append(allLabels, plan.Label)
		positions.Add(point)
	}

	pairs := this.createSquarePairs(perimeterCount, positions)
	return allLabels, positions, pairs
}

func (this *SquareTopologyService) createSquarePairs(perimeterCount int, positions models.Positions) [][2]int {
	builder := newPairBuilder()

	// Perimeter loop draws the square outline.
	if perimeterCount >= 2 {
		for i := 0; i < perimeterCount; i++ {
			builder.add(i, (i+1)%perimeterCount)
		}
	}

	// Interior loop keeps the inner neutral zones connected to one another.
	interiorCount := len(positions) - perimeterCount
	if interiorCount >= 2 {
		for i := 0; i < interiorCount; i++ {
			builder.add(perimeterCount+i, perimeterCount+(i+1)%interiorCount)
		}
	}

	// Tie every interior zone back to its nearest edge zone (inward spokes).
	if perimeterCount > 0 {
		for i := perimeterCount; i < len(positions); i++ {
			nearest := nearestIndexInRange(positions, i, 0, perimeterCount)
			if nearest >= 0 {
				builder.add(i, nearest)
			}
		}
	}
	return builder.pairs
}
