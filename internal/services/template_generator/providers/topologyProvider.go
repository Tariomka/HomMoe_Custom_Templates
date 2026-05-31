package providers

import (
	"math/rand/v2"

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
	holdCityNeutralLetter string) template.Variant {
	pl := make([]string, len(playerLabels))
	copy(pl, playerLabels)
	if this.shufflePlayerZones {
		rand.Shuffle(len(pl), func(i, j int) { pl[i], pl[j] = pl[j], pl[i] })
	}

	if configuration.IsTournamentMode() && len(pl) == 2 {
		return buildVariantTournament(configuration, pl, neutralZones, tuning)
	}

	switch configuration.Topology {
	case config.TopologyHubAndSpoke:
		return buildVariantHubAndSpoke(configuration, pl, neutralZones, tuning, configuration.IsHubCityToHold())
	case config.TopologyChain:
		return buildVariantChain(configuration, pl, neutralZones, tuning, holdCityNeutralLetter)
	case config.TopologySharedWeb:
		return buildVariantSharedWeb(configuration, pl, neutralZones, tuning, holdCityNeutralLetter)
	case config.TopologyRandom, config.TopologyBalanced:
		return buildVariantRandom(configuration, pl, neutralZones, tuning, holdCityNeutralLetter)
	default:
		return topology.NewRingTopologyService().
			GetTopologyVariant(configuration, pl, neutralZones, tuning, holdCityNeutralLetter)
	}
}

func (this *TopologyProvider) ShufflePlayerZones(enabled bool) *TopologyProvider {
	this.shufflePlayerZones = enabled
	return this
}
