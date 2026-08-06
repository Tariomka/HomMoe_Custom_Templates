package template_generator

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
)

type ITemplateGenerator interface {
	SetConfiguration(configuration *config.GeneratorConfig)
	Generate() (*entities.RmgTemplate, []string)
}
