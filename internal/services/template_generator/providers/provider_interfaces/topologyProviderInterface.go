package provider_interfaces

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
)

type ITopologyProvider interface {
	CreateTopologyVariant(
		configuration config.GeneratorConfig,
		playerLabels []string,
		neutralZones neutral_zone.Plans,
		tuning models.GenerationTuning,
		holdCityNeutralLabel string) template_model.Variant
}
