package topology

import (
	"slices"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
)

type CirclesTopologyService struct {
	RandomTopologyService
}

func NewCirclesTopologyService() *CirclesTopologyService {
	return &CirclesTopologyService{
		RandomTopologyService: *NewRandomTopologyService(),
	}
}

func (this *CirclesTopologyService) CreateTopologyVariant(
	configuration config.GeneratorConfig,
	playerLabels []string,
	neutralZones models.NeutralZonePlans,
	tuning models.GenerationTuning,
	holdCityNeutralLabel string) entities.Variant {
	neutralLabels := make([]string, len(neutralZones))
	for i, zonePlan := range neutralZones {
		neutralLabels[i] = zonePlan.Label
	}
	isIsolated := configuration.NoDirectPlayerConnections && len(playerLabels) > 1
	allLabels := this.ZoneLabelProvider.CreateBalancedRingZoneLabels(playerLabels, neutralZones, 0)
	positions := models.CreatePositionsFromPlans(allLabels, playerLabels, neutralZones)
	pairs := this.createCirclesPairs(positions.CreateDelaunayTriangulation(), allLabels, playerLabels, neutralZones)
	connectionNames := this.createConnectionNames(playerLabels, allLabels, pairs, isIsolated)

	zones := this.createZones(configuration, playerLabels, allLabels, tuning, neutralZones, holdCityNeutralLabel, connectionNames)
	for index := range zones {
		position := positions[index]
		zones[index].GeneratorPosition = &[2]float64{position.X, position.Y}
		tier := 0
		if !slices.Contains(playerLabels, allLabels[index]) {
			tier = neutralZones.GetTier(allLabels[index])
		}
		zones[index].GeneratorRing = &tier
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

func (this *CirclesTopologyService) createCirclesPairs(
	pairs []models.ConnectionIndexes,
	allLabels, playerLabels []string,
	neutralZones models.NeutralZonePlans) []models.ConnectionIndexes {
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
