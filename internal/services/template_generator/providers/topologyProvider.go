package providers

import (
	"fmt"

	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_errors"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
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
	holdCityNeutralLabel string) (template_model.Variant, error) {
	creator, supported := this.services.Resolve(configuration.Topology)
	if !supported {
		return template_model.Variant{},
			fmt.Errorf("%w: %q", common_errors.ErrUnsupportedTopology, configuration.Topology)
	}

	if configuration.IsTournamentMode() && len(playerLabels) == 2 {
		creator = this.services.Tournament()
	}
	return creator(configuration, playerLabels, neutralZones, tuning, holdCityNeutralLabel), nil
}
