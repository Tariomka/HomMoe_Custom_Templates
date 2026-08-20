package providers

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/provider_interfaces"
)

type TopologyProvider struct {
	services provider_interfaces.ITopologyServiceLookup
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
	if configuration.IsTournamentMode() && len(playerLabels) == 2 {
		return this.services.Tournament()(
			configuration, playerLabels, neutralZones, tuning, holdCityNeutralLabel)
	}

	return this.services.Resolve(configuration.Topology)(
		configuration, playerLabels, neutralZones, tuning, holdCityNeutralLabel)
}
