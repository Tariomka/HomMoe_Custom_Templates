package providers

import (
	"math/rand/v2"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/linq"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutralZone"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology"
)

type TopologyProvider struct {
	shufflePlayerZones bool
}

func NewTopologyProvider() *TopologyProvider {
	return &TopologyProvider{}
}

func (this *TopologyProvider) CreateTopologyVariant(
	configuration config.GeneratorConfig,
	playerLabels []string,
	neutralZones neutralZone.Plans,
	tuning models.GenerationTuning,
	holdCityNeutralLabel string) entities.Variant {
	playerLabelsCopy := this.copyLabels(playerLabels)

	if configuration.IsTournamentMode() && len(playerLabelsCopy) == 2 {
		return topology.NewTournamentTopologyService().
			CreateTopologyVariant(configuration, playerLabelsCopy, neutralZones, tuning)
	}

	switch configuration.Topology {
	case config.TopologyHubAndSpoke:
		return topology.NewHubTopologyService().
			CreateTopologyVariant(configuration, playerLabelsCopy, neutralZones, tuning, configuration.IsHubCityToHold())
	case config.TopologyChain:
		return topology.NewChainTopologyService().
			CreateTopologyVariant(configuration, playerLabelsCopy, neutralZones, tuning, holdCityNeutralLabel)
	case config.TopologySharedWeb:
		return topology.NewSharedWebTopologyService().
			CreateTopologyVariant(configuration, playerLabelsCopy, neutralZones, tuning, holdCityNeutralLabel)
	case config.TopologyRandom:
		return topology.NewRandomTopologyService().
			CreateTopologyVariant(configuration, playerLabelsCopy, neutralZones, tuning, holdCityNeutralLabel)
	case config.TopologyCircles:
		return topology.NewCirclesTopologyService().
			CreateTopologyVariant(configuration, playerLabelsCopy, neutralZones, tuning, holdCityNeutralLabel)
	case config.TopologySquare:
		return topology.NewSquareTopologyService().
			CreateTopologyVariant(configuration, playerLabelsCopy, neutralZones, tuning, holdCityNeutralLabel)
	case config.TopologyGeometric:
		return topology.NewGeometricTopologyService().
			CreateTopologyVariant(configuration, playerLabelsCopy, neutralZones, tuning, holdCityNeutralLabel)
	case config.TopologyCross:
		return topology.NewCrossTopologyService().
			CreateTopologyVariant(configuration, playerLabelsCopy, neutralZones, tuning, holdCityNeutralLabel)
	case config.TopologyFractal:
		return topology.NewFractalTopologyService().
			CreateTopologyVariant(configuration, playerLabelsCopy, neutralZones, tuning, holdCityNeutralLabel)
	default: // config.TopologyDefault
		return topology.NewRingTopologyService().
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
