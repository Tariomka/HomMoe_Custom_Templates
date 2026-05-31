package topology

import (
	"math/rand/v2"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template"
)

type TopologyFactory struct {
	shufflePlayerZones bool
}

func NewTopologyFactory() *TopologyFactory {
	return &TopologyFactory{}
}

func (this *TopologyFactory) CreateTopologyVariant(
	configuration *config.GeneratorConfig,
	playerLetters []string,
	neutralZones []models.NeutralZonePlan,
	tuning models.GenerationTuning,
	holdCityNeutralLetter string) template.Variant {
	pl := make([]string, len(playerLetters))
	copy(pl, playerLetters)
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
		return NewRingTopologyService().
			GetTopologyVariant(configuration, pl, neutralZones, tuning, holdCityNeutralLetter)
	}
}

func (this *TopologyFactory) ShufflePlayerZones(enabled bool) *TopologyFactory {
	this.shufflePlayerZones = enabled
	return this
}
