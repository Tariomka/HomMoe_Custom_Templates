package template_generator

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
)

type ITemplateGenerator interface {
	SetConfiguration(configuration *config.GeneratorConfig)
	Generate() (*template_model.Template, []string)
}
