package provider_interfaces

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
)

type IMandatoryContentProvider interface {
	CreateContents(
		configuration config.GeneratorConfig,
		playerLabels []string,
		neutralZones neutral_zone.Plans) []entities.MandatoryContent

	CreateContentsForZones(
		configuration config.GeneratorConfig,
		zones []template_model.Zone) []entities.MandatoryContent
}
