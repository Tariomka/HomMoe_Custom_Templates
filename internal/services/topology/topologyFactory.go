package topology

import (
	"math/rand/v2"

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
	neutralZones []neutralZonePlan,
	tuning generationTuning,
	holdCityNeutralLetter string,
	hubIsHoldCity bool) template.Variant {
	pl := make([]string, len(playerLetters))
	copy(pl, playerLetters)
	if this.shufflePlayerZones {
		rand.Shuffle(len(pl), func(i, j int) { pl[i], pl[j] = pl[j], pl[i] })
	}

	isTournament := (configuration.TournamentRules != nil && configuration.TournamentRules.Enabled) ||
		(configuration.GameEndConditions != nil && configuration.GameEndConditions.VictoryCondition == "win_condition_6")
	if isTournament && len(pl) == 2 {
		return buildVariantTournament(configuration, pl, neutralZones, tuning)
	}

	switch configuration.Topology {
	case config.TopologyHubAndSpoke:
		return buildVariantHubAndSpoke(configuration, pl, neutralZones, tuning, hubIsHoldCity)
	case config.TopologyChain:
		return buildVariantChain(configuration, pl, neutralZones, tuning, holdCityNeutralLetter)
	case config.TopologySharedWeb:
		return buildVariantSharedWeb(configuration, pl, neutralZones, tuning, holdCityNeutralLetter)
	case config.TopologyRandom:
		return buildVariantRandom(configuration, pl, neutralZones, tuning, holdCityNeutralLetter)
	case config.TopologyBalanced:
		return buildVariantBalanced(configuration, pl, neutralZones, tuning, holdCityNeutralLetter)
	default:
		return buildVariantDefault(configuration, pl, neutralZones, tuning, holdCityNeutralLetter)
	}
}

func (this *TopologyFactory) ShufflePlayerZones(flag bool) *TopologyFactory {
	this.shufflePlayerZones = flag
	return this
}
