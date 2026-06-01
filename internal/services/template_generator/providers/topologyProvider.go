package providers

import (
	"math/rand/v2"

	"github.com/Tariomka/hommoe_custom_templates/internal/linq"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template"
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
	neutralZones []models.NeutralZonePlan,
	tuning models.GenerationTuning,
	holdCityNeutralLabel string) template.Variant {
	playerLabelsCopy := this.copyLabels(playerLabels)

	if configuration.IsTournamentMode() && len(playerLabelsCopy) == 2 {
		return buildVariantTournament(configuration, playerLabelsCopy, neutralZones, tuning)
	}

	switch configuration.Topology {
	case config.TopologyHubAndSpoke:
		return buildVariantHubAndSpoke(configuration, playerLabelsCopy, neutralZones, tuning, configuration.IsHubCityToHold())
	case config.TopologyChain:
		return buildVariantChain(configuration, playerLabelsCopy, neutralZones, tuning, holdCityNeutralLabel)
	case config.TopologySharedWeb:
		return buildVariantSharedWeb(configuration, playerLabelsCopy, neutralZones, tuning, holdCityNeutralLabel)
	case config.TopologyRandom, config.TopologyBalanced:
		return buildVariantRandom(configuration, playerLabelsCopy, neutralZones, tuning, holdCityNeutralLabel)
	default:
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
		rand.Shuffle(len(playerLabelsCopy), func(i, j int) { playerLabelsCopy[i], playerLabelsCopy[j] = playerLabelsCopy[j], playerLabelsCopy[i] })
	}
	return playerLabelsCopy
}
