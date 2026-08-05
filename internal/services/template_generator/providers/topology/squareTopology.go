package topology

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology/base"
	zone_services "github.com/Tariomka/hommoe_custom_templates/internal/services/zones"
)

// SquareTopologyService lays the player zones out along the edges of a square.
// The perimeter zones are joined in a closed loop so the connections trace the
// square outline, while a share of the neutral zones is pulled into the middle
// of the square and tied back to the nearest edge.
type SquareTopologyService struct {
	PositionedTopologyBuilder
}

func NewSquareTopologyService(
	zoneFactory *zone_services.ZoneFactory,
	roadFactory *zone_services.RoadFactory,
	zoneLabelProvider zone_services.IZoneLabelProvider,
	connectionService *base.TopologyConnectionService) *SquareTopologyService {
	return &SquareTopologyService{
		PositionedTopologyBuilder: *NewPositionedTopologyBuilder(
			zoneFactory, roadFactory, zoneLabelProvider, connectionService),
	}
}

func (this *SquareTopologyService) CreateTopologyVariant(
	configuration config.GeneratorConfig,
	playerLabels []string,
	neutralZones neutral_zone.Plans,
	tuning models.GenerationTuning,
	holdCityNeutralLabel string) entities.Variant {
	return this.BuildVariant(
		configuration, playerLabels, neutralZones, tuning, holdCityNeutralLabel, this.createSquareLayout, nil)
}

// createSquareLayout builds the parallel label, position and connection-pair
// slices for the square. Index 0..perimeterCount-1 are the edge zones (players
// and the neutral zones that share the perimeter); the remaining indices are
// the interior neutral zones.
func (this *SquareTopologyService) createSquareLayout(
	playerLabels []string,
	neutralZones neutral_zone.Plans) ([]string, models.Positions, []models.ConnectionIndexes) {
	const half = 0.42

	// Roughly a third of the neutral zones (the higher tiers, which sort last)
	// are pulled inside the square; the rest line the edges between players.
	interiorCount := len(neutralZones) / 3
	perimeterPlans := neutralZones[:len(neutralZones)-interiorCount]
	interiorPlans := neutralZones[len(neutralZones)-interiorCount:]

	// Even player spacing around the perimeter, separated by the edge neutrals.
	perimeterLabels := this.ZoneLabelProvider.CreateBalancedRingZoneLabels(playerLabels, perimeterPlans)

	var allLabels []string
	var positions models.Positions

	perimeterCount := len(perimeterLabels)
	for i, label := range perimeterLabels {
		fraction := 0.0
		if perimeterCount > 0 {
			fraction = float64(i) / float64(perimeterCount)
		}
		allLabels = append(allLabels, label)
		positions.Add(squarePerimeterPoint(fraction, half))
	}

	// Interior neutral zones sit on a smaller inner square (or the exact center
	// when there is only one) so they read as being inside the perimeter.
	interiorHalf := half * 0.45
	for i, plan := range interiorPlans {
		var point models.Position
		if len(interiorPlans) == 1 {
			point = data.NewVec2(layoutCenter, layoutCenter)
		} else {
			point = squarePerimeterPoint(float64(i)/float64(len(interiorPlans)), interiorHalf)
		}
		allLabels = append(allLabels, plan.Label)
		positions.Add(point)
	}

	pairs := this.createSquarePairs(perimeterCount, positions)
	return allLabels, positions, pairs
}

func (this *SquareTopologyService) createSquarePairs(
	perimeterCount int,
	positions models.Positions) []models.ConnectionIndexes {
	builder := newPairBuilder()

	// Perimeter loop draws the square outline.
	if perimeterCount >= 2 {
		for i := range perimeterCount {
			builder.add(i, (i+1)%perimeterCount)
		}
	}

	// Interior loop keeps the inner neutral zones connected to one another.
	interiorCount := len(positions) - perimeterCount
	if interiorCount >= 2 {
		for i := range interiorCount {
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
