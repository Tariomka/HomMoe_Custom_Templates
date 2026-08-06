package topology

import (
	"slices"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/geometry_helpers"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology/base"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology/position_layout"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones/zone_interfaces"
)

type CirclesTopologyService struct {
	PositionedTopologyBuilder

	positionLayoutService *position_layout.PositionLayoutService
}

func NewCirclesTopologyService(
	zoneFactory zone_interfaces.IZoneFactory,
	roadFactory zone_interfaces.IRoadFactory,
	zoneLabelProvider zone_interfaces.IZoneLabelProvider,
	connectionService base.ITopologyConnectionService,
) *CirclesTopologyService {
	return &CirclesTopologyService{
		PositionedTopologyBuilder: *NewPositionedTopologyBuilder(
			zoneFactory, roadFactory, zoneLabelProvider, connectionService),
		positionLayoutService: position_layout.NewPositionLayoutService(),
	}
}

func (this *CirclesTopologyService) CreateTopologyVariant(
	configuration config.GeneratorConfig,
	playerLabels []string,
	neutralZones neutral_zone.Plans,
	tuning models.GenerationTuning,
	holdCityNeutralLabel string) entities.Variant {
	return this.BuildVariant(
		configuration, playerLabels, neutralZones, tuning, holdCityNeutralLabel,
		this.createCirclesLayout, stampGeneratorRings)
}

func (this *CirclesTopologyService) createCirclesLayout(
	playerLabels []string,
	neutralZones neutral_zone.Plans,
) ([]string, models.Positions, []models.ConnectionIndexes) {
	allLabels := this.ZoneLabelProvider.CreateBalancedRingZoneLabels(playerLabels, neutralZones)
	positions := this.positionLayoutService.CreatePositionsFromPlans(allLabels, playerLabels, neutralZones)
	pairs := this.createCirclesPairs(
		geometry_helpers.CreateDelaunayTriangulation(positions),
		allLabels,
		playerLabels,
		neutralZones,
	)
	return allLabels, positions, pairs
}

func stampGeneratorRings(
	zones []entities.Zone,
	allLabels, playerLabels []string,
	neutralZones neutral_zone.Plans,
) {
	for index := range zones {
		tier := 0
		if !slices.Contains(playerLabels, allLabels[index]) {
			tier = neutralZones.GetTier(allLabels[index])
		}
		zones[index].GeneratorRing = &tier
	}
}

func (this *CirclesTopologyService) createCirclesPairs(
	pairs []models.ConnectionIndexes,
	allLabels, playerLabels []string,
	neutralZones neutral_zone.Plans) []models.ConnectionIndexes {
	presentTiers := map[int]bool{}
	for _, label := range allLabels {
		tier := 0
		if !slices.Contains(playerLabels, label) {
			tier = neutralZones.GetTier(label)
		}
		presentTiers[tier] = true
	}

	var filtered []models.ConnectionIndexes
	for _, pair := range pairs {
		tierA := 0
		if !slices.Contains(playerLabels, allLabels[pair.X]) {
			tierA = neutralZones.GetTier(allLabels[pair.X])
		}
		tierB := 0
		if !slices.Contains(playerLabels, allLabels[pair.Y]) {
			tierB = neutralZones.GetTier(allLabels[pair.Y])
		}
		low, high := tierA, tierB
		if low > high {
			low, high = high, low
		}
		if high-low <= 1 {
			filtered = append(filtered, pair)
			continue
		}

		skip := false
		for tier := low + 1; tier < high; tier++ {
			if presentTiers[tier] {
				skip = true
				break
			}
		}

		if !skip {
			filtered = append(filtered, pair)
		}
	}
	return filtered
}
