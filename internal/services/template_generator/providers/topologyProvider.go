package providers

import (
	"math/rand/v2"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/linq"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones"
)

type TopologyProvider struct {
	shufflePlayerZones bool
	zoneFactory        *zones.ZoneFactory
	roadFactory        *zones.RoadFactory
}

func NewTopologyProvider(
	zoneFactory *zones.ZoneFactory,
	roadFactory *zones.RoadFactory,
) *TopologyProvider {
	return &TopologyProvider{zoneFactory: zoneFactory, roadFactory: roadFactory}
}

func (this *TopologyProvider) CreateTopologyVariant(
	configuration config.GeneratorConfig,
	playerLabels []string,
	neutralZones neutral_zone.Plans,
	tuning models.GenerationTuning,
	holdCityNeutralLabel string) entities.Variant {
	playerLabelsCopy := this.copyLabels(playerLabels)

	if configuration.IsTournamentMode() && len(playerLabelsCopy) == 2 {
		return topology.NewTournamentTopologyService(this.zoneFactory, this.roadFactory).
			CreateTopologyVariant(configuration, playerLabelsCopy, neutralZones, tuning)
	}

	switch configuration.Topology {
	case config.TopologyHubAndSpoke:
		return topology.NewHubTopologyService(this.zoneFactory, this.roadFactory).
			CreateTopologyVariant(
				configuration, playerLabelsCopy, neutralZones, tuning, configuration.IsHubCityToHold())
	case config.TopologyGeometricHub:
		return topology.NewGeometricHubTopologyService(this.zoneFactory, this.roadFactory).
			CreateTopologyVariant(
				configuration, playerLabelsCopy, neutralZones, tuning, configuration.IsHubCityToHold())
	case config.TopologyChain:
		return topology.NewChainTopologyService(this.zoneFactory, this.roadFactory).
			CreateTopologyVariant(configuration, playerLabelsCopy, neutralZones, tuning, holdCityNeutralLabel)
	case config.TopologySharedWeb:
		return topology.NewSharedWebTopologyService(this.zoneFactory, this.roadFactory).
			CreateTopologyVariant(configuration, playerLabelsCopy, neutralZones, tuning, holdCityNeutralLabel)
	case config.TopologyRandom:
		return topology.NewRandomTopologyService(this.zoneFactory, this.roadFactory).
			CreateTopologyVariant(configuration, playerLabelsCopy, neutralZones, tuning, holdCityNeutralLabel)
	case config.TopologyCircles:
		return topology.NewCirclesTopologyService(this.zoneFactory, this.roadFactory).
			CreateTopologyVariant(configuration, playerLabelsCopy, neutralZones, tuning, holdCityNeutralLabel)
	case config.TopologySquare:
		return topology.NewSquareTopologyService(this.zoneFactory, this.roadFactory).
			CreateTopologyVariant(configuration, playerLabelsCopy, neutralZones, tuning, holdCityNeutralLabel)
	case config.TopologyGeometric:
		return topology.NewGeometricTopologyService(this.zoneFactory, this.roadFactory).
			CreateTopologyVariant(configuration, playerLabelsCopy, neutralZones, tuning, holdCityNeutralLabel)
	case config.TopologyCross:
		return topology.NewCrossTopologyService(this.zoneFactory, this.roadFactory).
			CreateTopologyVariant(configuration, playerLabelsCopy, neutralZones, tuning, holdCityNeutralLabel)
	case config.TopologyFractal:
		return topology.NewFractalTopologyService(this.zoneFactory, this.roadFactory).
			CreateTopologyVariant(configuration, playerLabelsCopy, neutralZones, tuning, holdCityNeutralLabel)
	default: // config.TopologyDefault
		return topology.NewRingTopologyService(this.zoneFactory, this.roadFactory).
			CreateTopologyVariant(configuration, playerLabelsCopy, neutralZones, tuning, holdCityNeutralLabel)
	}
}

func (this *TopologyProvider) ShufflePlayerZones(enabled bool) *TopologyProvider {
	this.shufflePlayerZones = enabled
	return this
}

func (this *TopologyProvider) copyLabels(playerLabels []string) []string {
	playerLabelsCopy := linq.FromSlice(playerLabels).ToSlice()
	if this.shufflePlayerZones {
		rand.Shuffle(len(playerLabelsCopy),
			func(i, j int) { playerLabelsCopy[i], playerLabelsCopy[j] = playerLabelsCopy[j], playerLabelsCopy[i] })
	}
	return playerLabelsCopy
}
