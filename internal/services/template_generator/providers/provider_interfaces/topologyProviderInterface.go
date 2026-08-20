package provider_interfaces

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
)

type ITopologyProvider interface {
	CreateTopologyVariant(
		configuration config.GeneratorConfig,
		playerLabels []string,
		neutralZones neutral_zone.Plans,
		tuning models.GenerationTuning,
		holdCityNeutralLabel string) entities.Variant
}
