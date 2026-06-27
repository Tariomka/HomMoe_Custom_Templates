package dtos

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
)

type TemplateUpdateDto struct {
	Template    *entities.RmgTemplate
	Zones       []entities.Zone
	Connections []entities.Connection

	// Config, when supplied, lets UpdateTemplate rebuild the template's
	// mandatory content from the final (manually edited) zones so a zone whose
	// quality was changed in the editor gets the content of its new tier.
	Config *config.GeneratorConfig
}
