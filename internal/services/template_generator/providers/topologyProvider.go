package providers

import (
	"math/rand/v2"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/linq"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/provider_interfaces"
)

type TopologyProvider struct {
	shufflePlayerZones bool
	services           provider_interfaces.ITopologyServiceLookup
}

func NewTopologyProvider(
	services provider_interfaces.ITopologyServiceLookup) provider_interfaces.ITopologyProvider {
	return &TopologyProvider{services: services}
}

func (this *TopologyProvider) CreateTopologyVariant(
	configuration config.GeneratorConfig,
	playerLabels []string,
	neutralZones neutral_zone.Plans,
	tuning models.GenerationTuning,
	holdCityNeutralLabel string) entities.Variant {
	playerLabelsCopy := this.copyLabels(playerLabels)

	if configuration.IsTournamentMode() && len(playerLabelsCopy) == 2 {
		return this.services.Tournament()(
			configuration, playerLabelsCopy, neutralZones, tuning, holdCityNeutralLabel)
	}

	return this.services.Resolve(configuration.Topology)(
		configuration, playerLabelsCopy, neutralZones, tuning, holdCityNeutralLabel)
}

func (this *TopologyProvider) ShufflePlayerZones(enabled bool) provider_interfaces.ITopologyProvider {
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
